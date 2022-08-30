package main

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/bot"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/douyu-task/api"
)

func main() {
	app.Add(startup)
	app.Add(cron)
	err := app.Go()
	if err != nil {
		log.Fatal().Msgf("app init fail: %v", err)
	}
}

func startup(context.Context) error {
	if config.GetBool("startup") {
		NewRenewal().Run()
	}
	return nil
}

func cron(context.Context) error {
	_, err := cronx.AddJob("0 0 12 * * 0", NewRenewal())
	if err != nil {
		return err
	}
	go cronx.Run()
	return nil
}

// --------------------------------------------------------------------------------

// NewRenewal 送免费的荧光棒续牌子
func NewRenewal() *Renewal {
	r := &Renewal{
		did:  config.GetString("douyu.did"),
		uid:  config.GetString("douyu.uid"),
		auth: config.GetString("douyu.auth"),
		ltp0: config.GetString("douyu.ltp0"),
	}
	if r.did == "" || r.uid == "" || (r.auth == "" && r.ltp0 == "") {
		log.Fatal().Msgf("douyu missing config")
	}
	r.stickRemaining = config.GetInt("douyu.stick.remaining")
	return r
}

type Renewal struct {
	did  string // cookie: dy_did
	uid  string // cookie: acf_uid
	auth string // cookie: acf_auth
	ltp0 string // cookie: ltp0

	stickRemaining int // 房间号，剩余的荧光棒送给谁
}

func (r *Renewal) Name() string {
	return "renewal"
}

func (r *Renewal) Run() {
	err := r.Send()
	if err != nil {
		log.Error().Msgf("renewal: %v", err)
		_ = bot.Send("续牌子失败：" + err.Error())
	} else {
		_ = bot.Send("续牌子成功")
	}
}

func (r *Renewal) Send() error {
	c := api.New(r.did, r.uid, r.auth, r.ltp0)

	if r.ltp0 != "" {
		err := c.Refresh()
		if err != nil {
			return fmt.Errorf("refresh fail: %v", err)
		}
	}

	gifts, err := c.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts fail: %v", err)
	}

	stick := gifts.Find(api.GiftGlowSticks).GetCount()
	if stick == 0 {
		return fmt.Errorf("no free gift")
	}

	badges, _, err := r.Badges(c, true)
	if err != nil {
		return fmt.Errorf("list badges fail: %v", err)
	}

	for _, badge := range badges {
		if badge.Room == r.stickRemaining {
			continue
		}
		gifts, err = c.SendGift(badge.Room, api.GiftGlowSticks, 1)
		if err != nil {
			log.Error().Msgf("send gift fail: %v", err)
			continue
		}
		stick = gifts.Find(api.GiftGlowSticks).GetCount()
		if stick == 0 {
			log.Error().Msg("no remaining free gift")
			break
		}
		log.Info().Msgf("send gift success, %s", badge.Anchor)
		time.Sleep(time.Second)
	}

	if stick == 0 {
		return nil
	}

	gifts, err = c.SendGift(r.stickRemaining, api.GiftGlowSticks, stick)
	if err != nil {
		return fmt.Errorf("send gift fail: %v", err)
	}

	_, _, err = r.Badges(c, true)
	if err != nil {
		return fmt.Errorf("list badges fail: %v", err)
	}

	return nil
}

func (r *Renewal) Badges(c *api.Client, output bool) (map[int]*api.Badge, []int, error) {
	badges, err := c.ListBadges()
	if err != nil {
		return nil, nil, err
	}

	bm, rs := map[int]*api.Badge{}, make([]int, len(badges))

	for i, badge := range badges {
		rs[i] = badge.Room
		bm[badge.Room] = badges[i]
	}

	sort.Slice(rs, func(i, j int) bool { return bm[rs[i]].Intimacy > bm[rs[j]].Intimacy })

	if output {
		bb := &bytes.Buffer{}
		tw := tablewriter.NewWriter(bb)
		tw.SetAlignment(tablewriter.ALIGN_CENTER)
		tw.SetHeader([]string{"room", "anchor", "name", "level", "intimacy", "rank"})
		for i := 0; i < len(rs); i++ {
			b := bm[rs[i]]
			tw.Append([]string{strconv.Itoa(b.Room), b.Anchor, b.Name, strconv.Itoa(b.Level), strconv.FormatFloat(b.Intimacy, 'f', -1, 64), strconv.Itoa(b.Rank)})
		}
		tw.Render()
		log.Info().Msgf("badges:\n%s", bb.String())
	}

	return bm, rs, nil
}
