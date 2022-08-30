package api

import (
	"testing"

	"github.com/starudream/go-lib/testx"
)

func TestClient_Refresh(t *testing.T) {
	err := client.Refresh()
	testx.P(t, err)
}
