package job

import (
	"errors"
	"fmt"

	"github.com/starudream/go-lib/core/v2/slog"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/consts"
)

var ErrTodayNotExpired = errors.New("today no expired")

func Renewal() error {
	client, err := api.NewC(config.C().Douyu)
	if err != nil {
		return fmt.Errorf("init api error: %w", err)
	}

	badges1, err := client.ListBadges()
	if err != nil {
		return fmt.Errorf("list badges error: %w", err)
	}
	slog.Info("badges before:\n%s", badges1.TableString())

	gifts1, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	slog.Info("gifts before:\n%s", gifts1.TableString())

	id := gifts1.NotEmpty(consts.GiftFansGlowSticks, consts.GiftGlowSticks)
	if id == -1 {
		return fmt.Errorf("no free gift")
	}

	gift := gifts1.Find(id)

	if config.C().Douyu.IgnoreExpiredCheck {
		slog.Info("ignore expired check")
	} else if !gift.TodayExpired() {
		return ErrTodayNotExpired
	}

	count := gift.GetCount()

	assigns := config.C().Douyu.Assigns

	for i := 0; i < len(assigns); i++ {
		a := assigns[i]

		if count <= 0 {
			break
		}

		if a.All {
			count, err = SendGift(client, a.Room, id, count)
			if err != nil {
				return err
			}
			continue
		}

		if a.Room <= 0 {
			for j := 0; j < len(badges1); j++ {
				count, err = SendGift(client, badges1[j].Room, id, a.Count)
				if err != nil {
					return err
				}
			}
			continue
		}

		count, err = SendGift(client, a.Room, id, a.Count)
		if err != nil {
			return err
		}
	}

	gifts2, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	slog.Info("gifts after:\n%s", gifts2.TableString())

	badges2, err := client.ListBadges()
	if err != nil {
		return fmt.Errorf("list badges error: %w", err)
	}
	slog.Info("badges after:\n%s", badges2.TableString())

	return nil
}

func SendGift(client *api.Client, room, gift, count int) (int, error) {
	slog.Info("send gift(%d) count(%d) to room(%d)", gift, count, room)
	gs, e := client.SendGift(room, gift, count)
	if e != nil {
		return 0, fmt.Errorf("send gift error: %w", e)
	}
	left := gs.Find(gift).GetCount()
	slog.Info("left gift(%d) count(%d)", gift, count)
	return left, nil
}
