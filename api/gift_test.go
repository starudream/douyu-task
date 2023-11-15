package api

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"

	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/consts"
)

func TestClient_SendGift(t *testing.T) {
	resp, err := MustNew(config.C().Douyu).SendGift(consts.RoomYYF, consts.GiftGlowSticks, 1)
	testutil.LogNoErr(t, err, resp)
}

func TestClient_ListGift(t *testing.T) {
	resp, err := MustNew(config.C().Douyu).ListGifts()
	testutil.LogNoErr(t, err, resp.TableString())
}

func TestGift_TodayExpired(t *testing.T) {
	resp, err := MustNew(config.C().Douyu).ListGifts()
	testutil.LogNoErr(t, err, resp.TableString(), resp.Find(consts.GiftGlowSticks).TodayExpired())
}
