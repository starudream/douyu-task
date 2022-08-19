package api

import (
	"testing"

	"github.com/starudream/go-lib/testx"
)

func TestClient_ListBadges(t *testing.T) {
	resp, err := client.ListBadges()
	testx.P(t, err, resp)
}
