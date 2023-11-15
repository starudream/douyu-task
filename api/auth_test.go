package api

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"

	"github.com/starudream/douyu-task/config"
)

func TestClient_Refresh(t *testing.T) {
	err := MustNew(config.C().Douyu).Refresh()
	testutil.Nil(t, err)
}
