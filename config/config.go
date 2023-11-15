package config

import (
	"sync"

	"github.com/starudream/go-lib/core/v2/config"
)

type Config struct {
	Cron  Cron  `json:"cron" yaml:"cron"`
	Douyu Douyu `json:"douyu" yaml:"douyu"`
}

type Cron struct {
	Refresh string `json:"refresh" yaml:"refresh"`
	Renewal string `json:"renewal" yaml:"renewal"`
}

type Douyu struct {
	Did     string   `json:"did"     yaml:"did"`
	Ltp0    string   `json:"ltp0"    yaml:"ltp0"`
	Room    int      `json:"room"    yaml:"room"`
	Assigns []Assign `json:"assigns" yaml:"assigns"`

	IgnoreExpiredCheck bool `json:"ignore_expired_check" yaml:"ignore_expired_check"`
}

type Assign struct {
	Count int  `json:"count,omitempty" yaml:"count,omitempty"`
	Room  int  `json:"room,omitempty"  yaml:"room,omitempty"`
	All   bool `json:"all,omitempty"   yaml:"all,omitempty"`
}

var (
	_c = Config{
		Cron: Cron{
			Refresh: "0 10 0 * * *",
			Renewal: "0 50 23 * * *",
		},
		Douyu: Douyu{},
	}
	_cMu = sync.Mutex{}
)

func init() {
	_ = config.Unmarshal("cron", &_c.Cron)
	_ = config.Unmarshal("douyu", &_c.Douyu)
	config.LoadStruct(_c)
}

func C() Config {
	_cMu.Lock()
	defer _cMu.Unlock()
	return _c
}
