package api

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"

	"github.com/starudream/douyu-task/config"
)

func TestClient_ListBadges(t *testing.T) {
	resp, err := MustNew(config.C().Douyu).ListBadges()
	testutil.LogNoErr(t, err, resp.TableString())
}
