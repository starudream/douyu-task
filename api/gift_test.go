package api

import (
	"testing"

	"github.com/starudream/go-lib/testx"

	"github.com/starudream/douyu-task/consts"
)

func TestClient_SendGift(t *testing.T) {
	resp, err := client.SendGift(consts.RoomYYF, consts.GiftGlowSticks, 1)
	testx.P(t, err, resp)
}

func TestClient_ListGift(t *testing.T) {
	resp, err := client.ListGifts()
	testx.P(t, err, resp.TableString())
}

func TestGift_TodayExpired(t *testing.T) {
	resp, err := client.ListGifts()
	testx.P(t, err, resp.TableString(), resp.Find(consts.GiftGlowSticks).TodayExpired())
}
