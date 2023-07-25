package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/flag"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/consts"
	"github.com/starudream/douyu-task/internal/osx"
)

var runCmd = &flag.Command{
	Use:   "run",
	Short: "just run once for test/debug",
}

var runLoginCmd = &flag.Command{
	Use:   "login",
	Short: "login douyu through websocket to get glow sticks",
	Run: func(cmd *flag.Command, args []string) {
		osx.PE(NewRefresh().do())
	},
}

var runGiftCmd = &flag.Command{
	Use:   "gift",
	Short: "douyu gift",
}

var runGiftListCmd = &flag.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list gift you have",
	Run: func(cmd *flag.Command, args []string) {
		c, e := api.NewFromEnv()
		osx.PE(e)
		gs, e := c.ListGifts(config.GetInt("douyu.room"))
		osx.PE(e)
		osx.PA(gs.TableString())
	},
}

var runGiftSendCmd = &flag.Command{
	Use:   "send",
	Short: "send gift you have",
	Args:  flag.ExactArgs(1),
	Run: func(cmd *flag.Command, args []string) {
		count, e := strconv.Atoi(args[0])
		osx.PE(e)

		c, e := api.NewFromEnv()
		osx.PE(e)

		bs1, e := c.ListBadges()
		osx.PE(e, bs1.TableString())

		gs1, e := c.ListGifts(config.GetInt("douyu.room"))
		osx.PE(e, gs1.TableString())

		fmt.Printf("send gift(%d) count(%d) to room(%d), please confirm(y/n): ", config.GetInt("douyu.gift"), count, config.GetInt("douyu.room"))

		var s string
		_, e = fmt.Scanf("%s", &s)
		osx.PE(e)
		if xx := strings.ToLower(s); xx != "y" && xx != "yes" {
			osx.PA("abort")
		}

		_, e = c.SendGift(config.GetInt("douyu.room"), config.GetInt("douyu.gift"), count)
		osx.PE(e)

		gs2, e := c.ListGifts(config.GetInt("douyu.room"))
		osx.PE(e, gs2.TableString())

		bs2, e := c.ListBadges()
		osx.PE(e, bs2.TableString())
	},
}

var runBadgeCmd = &flag.Command{
	Use:   "badge",
	Short: "douyu badge",
}

var runBadgeListCmd = &flag.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list badge you have",
	Run: func(cmd *flag.Command, args []string) {
		c, e := api.NewFromEnv()
		osx.PE(e)
		bs, e := c.ListBadges()
		osx.PE(e)
		osx.PA(bs.TableString())
	},
}

func init() {
	runGiftCmd.PersistentFlags().Int("room-id", consts.RoomYYF, "room id")
	osx.PE(config.BindPFlag("douyu.room", runGiftCmd.PersistentFlags().Lookup("room-id")))

	runGiftSendCmd.PersistentFlags().Int("gift-id", consts.GiftGlowSticks, "gift id (268 or 2358)")
	osx.PE(config.BindPFlag("douyu.gift", runGiftSendCmd.PersistentFlags().Lookup("gift-id")))

	runGiftCmd.AddCommand(runGiftListCmd)
	runGiftCmd.AddCommand(runGiftSendCmd)

	runBadgeCmd.AddCommand(runBadgeListCmd)

	runCmd.AddCommand(runLoginCmd)
	runCmd.AddCommand(runGiftCmd)
	runCmd.AddCommand(runBadgeCmd)
}
