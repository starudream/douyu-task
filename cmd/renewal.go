package main

import (
	"fmt"

	"github.com/starudream/go-lib/bot"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/consts"
)

// NewRenewal 送免费的荧光棒续牌子
func NewRenewal() *Renewal {
	r := &Renewal{
		room: config.GetInt("douyu.room"),
	}
	assigns, err := config.UnmarshalKeyTo[[]Assign]("douyu.assigns")
	if err != nil {
		log.Fatal().Msgf("douyu assigns error: %v", err)
	}
	r.assigns = assigns
	return r
}

type Assign struct {
	Room  int  `json:"room,omitempty" yaml:"room,omitempty"`
	Count int  `json:"count,omitempty" yaml:"count,omitempty"`
	All   bool `json:"all,omitempty" yaml:"all,omitempty"`
}

type Renewal struct {
	room int

	assigns []Assign
}

var _ cronx.Job = (*Renewal)(nil)

func (r *Renewal) Name() string {
	return "renewal"
}

func (r *Renewal) Run() {
	err := r.do()
	if err != nil {
		log.Error().Msgf("renewal error: %v", err)
		_ = bot.Send("续牌子失败：" + err.Error())
	} else {
		_ = bot.Send("续牌子成功")
	}
}

func (r *Renewal) do() error {
	client, err := api.NewFromEnv()
	if err != nil {
		return fmt.Errorf("init api error: %w", err)
	}

	badges1, err := client.ListBadges()
	if err != nil {
		return fmt.Errorf("list badges error: %w", err)
	}
	log.Info().Msgf("badges before:\n%s", badges1.TableString())

	gifts1, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	log.Info().Msgf("gifts before:\n%s", gifts1.TableString())

	id := gifts1.NotEmpty(consts.GiftFansGlowSticks, consts.GiftGlowSticks)
	if id == -1 {
		return fmt.Errorf("no free gift")
	}

	count := gifts1.Find(id).GetCount()

	for i := 0; i < len(r.assigns); i++ {
		a := r.assigns[i]

		if count <= 0 {
			break
		}

		if a.All {
			count, err = r.sendGift(client, a.Room, id, count)
			if err != nil {
				return err
			}
			continue
		}

		if a.Room <= 0 {
			for j := 0; j < len(badges1); j++ {
				count, err = r.sendGift(client, badges1[j].Room, id, a.Count)
				if err != nil {
					return err
				}
			}
			continue
		}

		count, err = r.sendGift(client, a.Room, id, a.Count)
		if err != nil {
			return err
		}
	}

	gifts2, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	log.Info().Msgf("gifts after:\n%s", gifts2.TableString())

	badges2, err := client.ListBadges()
	if err != nil {
		return fmt.Errorf("list badges error: %w", err)
	}
	log.Info().Msgf("badges after:\n%s", badges2.TableString())

	return nil
}

func (r *Renewal) sendGift(client *api.Client, room, gift, count int) (int, error) {
	log.Info().Msgf("send gift(%d) count(%d) to room(%d)", gift, count, room)
	gs, e := client.SendGift(room, gift, count)
	if e != nil {
		return 0, fmt.Errorf("send gift error: %w", e)
	}
	left := gs.Find(gift).GetCount()
	log.Info().Msgf("left gift(%d) count(%d)", gift, count)
	return left, nil
}
