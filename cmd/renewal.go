package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/starudream/go-lib/cronx"

	"github.com/starudream/go-lib/bot"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/douyu-task/api"
)

// NewRenewal 送免费的荧光棒续牌子
func NewRenewal() *Renewal {
	r := &Renewal{
		did:  config.GetString("douyu.did"),
		uid:  config.GetString("douyu.uid"),
		auth: config.GetString("douyu.auth"),
		ltp0: config.GetString("douyu.ltp0"),
	}
	if r.did == "" || r.uid == "" || (r.auth == "" && r.ltp0 == "") {
		log.Fatal().Msgf("douyu renewal missing config")
	}
	r.stickRemaining = config.GetInt("douyu.stick.remaining")
	r.client = api.New(r.did, r.uid, r.auth, r.ltp0)
	return r
}

type Renewal struct {
	did  string // cookie: dy_did
	uid  string // cookie: acf_uid
	auth string // cookie: acf_auth
	ltp0 string // cookie: ltp0

	stickRemaining int // 房间号，剩余的荧光棒送给谁

	client *api.Client
}

var _ cronx.Job = (*Renewal)(nil)

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
	gifts, err := r.Gifts()
	if err != nil {
		return err
	}

	id := gifts.NotEmpty(api.GiftFansGlowSticks, api.GiftGlowSticks)
	if id == -1 {
		return fmt.Errorf("no free gift")
	}

	badges, _, err := r.Badges()
	if err != nil {
		return fmt.Errorf("list badges fail: %v", err)
	}

	var stick int

	for _, badge := range badges {
		if badge.Room == r.stickRemaining {
			continue
		}
		gifts, err = r.client.SendGift(badge.Room, id, 1)
		if err != nil {
			log.Error().Msgf("send gift fail: %v", err)
			continue
		}
		stick = gifts.Find(id).GetCount()
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

	gifts, err = r.client.SendGift(r.stickRemaining, id, stick)
	if err != nil {
		return fmt.Errorf("send gift fail: %v", err)
	}

	_, _, err = r.Badges()
	if err != nil {
		return fmt.Errorf("list badges fail: %v", err)
	}

	return nil
}

func (r *Renewal) Gifts() (*api.Gifts, error) {
	if r.ltp0 != "" {
		err := r.client.Refresh()
		if err != nil {
			return nil, fmt.Errorf("refresh fail: %v", err)
		}
	}

	gifts, err := r.client.ListGifts()
	if err != nil {
		return nil, fmt.Errorf("list gifts fail: %v", err)
	}

	bb := &bytes.Buffer{}
	tw := tablewriter.NewWriter(bb)
	tw.SetAlignment(tablewriter.ALIGN_CENTER)
	tw.SetHeader([]string{"name", "price", "count", "expire"})
	for i := 0; i < len(gifts.List); i++ {
		g := gifts.List[i]
		tw.Append([]string{g.Name, strconv.Itoa(g.Price), strconv.Itoa(g.Count), time.Unix(int64(g.Met), 0).Format(time.DateTime)})
	}
	tw.Render()
	log.Info().Msgf("gifts:\n%s", bb.String())

	return gifts, nil
}

func (r *Renewal) Badges() (map[int]*api.Badge, []int, error) {
	badges, err := r.client.ListBadges()
	if err != nil {
		return nil, nil, err
	}

	bm, rs := map[int]*api.Badge{}, make([]int, len(badges))

	for i, badge := range badges {
		rs[i] = badge.Room
		bm[badge.Room] = badges[i]
	}

	sort.Slice(rs, func(i, j int) bool { return bm[rs[i]].Intimacy > bm[rs[j]].Intimacy })

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

	return bm, rs, nil
}
