package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/fmtutil"
	"github.com/starudream/go-lib/core/v2/utils/sliceutil"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/consts"
	"github.com/starudream/douyu-task/job"
)

var (
	runCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "run"
		c.Short = "Run douyu job manually"
	})

	runRefreshCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "refresh"
		c.Short = "Simulate login to Douyu and get free gifts"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			return job.Refresh()
		}
	})

	runGiftCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "gift"
		c.Short = "Manage your gifts"
	})

	runGiftListCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "list"
		c.Short = "List your gifts"
		c.Aliases = []string{"ls"}
		c.RunE = func(cmd *cobra.Command, args []string) error {
			client, err := api.NewC(config.C().Douyu)
			if err != nil {
				return fmt.Errorf("new client error: %w", err)
			}
			gifts, err := client.ListGifts()
			if err != nil {
				return fmt.Errorf("list gifts error: %w", err)
			}
			fmt.Println(gifts.TableString())
			return nil
		}
	})

	runGiftSendCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "send <room id> <gift id> <count>"
		c.Short = "Send gift to room"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			var (
				roomId = atoi(args, 0, strconv.Itoa(consts.RoomYYF))
				giftId = atoi(args, 1, strconv.Itoa(consts.GiftGlowSticks))
				count  = atoi(args, 2, "1")
			)

			client, err := api.NewC(config.C().Douyu)
			if err != nil {
				return fmt.Errorf("new client error: %w", err)
			}

			badges, err := client.ListBadges()
			if err != nil {
				return fmt.Errorf("list badges error: %w", err)
			}
			fmt.Println(badges.TableString())

			gifts, err := client.ListGifts()
			if err != nil {
				return fmt.Errorf("list gifts error: %w", err)
			}
			fmt.Println(gifts.TableString())

			confirm := fmtutil.Scan(fmt.Sprintf("send gift to room id %d, gift id %d, count %d, confirm? (y/n)", roomId, giftId, count))
			if s := strings.ToLower(strings.TrimSpace(confirm)); s != "y" && s != "yes" {
				fmt.Println("canceled")
				return nil
			}

			_, err = client.SendGift(roomId, giftId, count)
			if err != nil {
				return fmt.Errorf("send gift error: %w", err)
			}

			gifts, err = client.ListGifts()
			if err != nil {
				return fmt.Errorf("list gifts error: %w", err)
			}
			fmt.Println(gifts.TableString())

			badges, err = client.ListBadges()
			if err != nil {
				return fmt.Errorf("list badges error: %w", err)
			}
			fmt.Println(badges.TableString())

			return nil
		}
	})

	runBadgeCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "badge"
		c.Short = "Manage your badges"
	})

	runBadgeListCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "list"
		c.Short = "List your badges"
		c.Aliases = []string{"ls"}
		c.RunE = func(cmd *cobra.Command, args []string) error {
			client, err := api.NewC(config.C().Douyu)
			if err != nil {
				return fmt.Errorf("new client error: %w", err)
			}
			badges, err := client.ListBadges()
			if err != nil {
				return fmt.Errorf("list badges error: %w", err)
			}
			fmt.Println(badges.TableString())
			return nil
		}
	})
)

func init() {
	runCmd.AddCommand(runRefreshCmd)

	runGiftCmd.AddCommand(runGiftListCmd)
	runGiftCmd.AddCommand(runGiftSendCmd)
	runCmd.AddCommand(runGiftCmd)

	runBadgeCmd.AddCommand(runBadgeListCmd)
	runCmd.AddCommand(runBadgeCmd)

	rootCmd.AddCommand(runCmd)
}

func atoi(vs []string, idx int, def ...string) int {
	v, _ := sliceutil.GetValue(vs, idx, def...)
	i, _ := strconv.Atoi(v)
	return i
}
