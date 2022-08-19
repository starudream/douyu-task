package api

import (
	"testing"

	"github.com/starudream/go-lib/testx"
)

func TestClient_SendGift(t *testing.T) {
	resp, err := client.SendGift(RoomYYF, GiftGlowSticks, 1)
	testx.P(t, err, resp)
}

func TestClient_ListGift(t *testing.T) {
	resp, err := client.ListGifts()
	testx.P(t, err, resp)
}
