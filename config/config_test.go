package config

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/codec/yaml"
	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func Test(t *testing.T) {
	testutil.Log(t, "\n"+yaml.MustMarshalString(config.Raw()))
}
