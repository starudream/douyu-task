package api

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/starudream/go-lib/tablew/v2"

	"github.com/starudream/douyu-task/internal/htmlx"
)

type Badges []*Badge

func (bs Badges) TableString() string {
	if len(bs) == 0 {
		return "empty badges"
	}
	return tablew.Structs(bs)
}

// Badge 徽章
type Badge struct {
	Room     int       // 房间号
	Anchor   string    // 主播名称
	Name     string    // 名称
	Level    int       // 等级
	Intimacy float64   // 亲密度
	Rank     int       // 排名
	AccessAt time.Time // 获得时间
}

func (c *Client) ListBadges() (Badges, error) {
	resp, err := c.R().
		SetCookies(c.genAuthCookies()).
		SetHeader("referer", URL).
		Get(URL + "/member/cp/getFansBadgeList")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status())
	}

	root, err := htmlx.Parse(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	title := htmlx.NodeTitle(root)
	if !strings.Contains(title, "我的头衔") {
		return nil, fmt.Errorf("page title not match: %s", title)
	}

	table := htmlx.NodeSearch(root, func(node *htmlx.Node) bool {
		return node.Type == htmlx.ElementNode && strings.TrimSpace(node.Data) == "table" && htmlx.NodeAttrExists(node, func(attr htmlx.Attribute) bool {
			return attr.Key == "class" && strings.Contains(attr.Val, "fans-badge-list")
		})
	})
	if table == nil {
		return nil, fmt.Errorf("no table")
	}

	tbody := htmlx.NodeSearch(table, func(node *htmlx.Node) bool {
		return node.Type == htmlx.ElementNode && strings.TrimSpace(node.Data) == "tbody"
	})
	if tbody == nil {
		return nil, fmt.Errorf("no table body")
	}

	trs := htmlx.NodeChildren(tbody, "tr")
	if len(trs) == 0 {
		return nil, fmt.Errorf("no table rows")
	}

	badges := make([]*Badge, len(trs))

	for i, tr := range trs {
		badge := &Badge{}
		for _, attr := range tr.Attr {
			switch attr.Key {
			case "data-fans-room":
				badge.Room, _ = strconv.Atoi(attr.Val)
			case "data-fans-level":
				badge.Level, _ = strconv.Atoi(attr.Val)
			case "data-fans-intimacy":
				badge.Intimacy, _ = strconv.ParseFloat(attr.Val, 64)
			case "data-fans-rank":
				badge.Rank, _ = strconv.Atoi(attr.Val)
			case "data-fans-gbdgts":
				v, _ := strconv.Atoi(attr.Val)
				badge.AccessAt = time.Unix(int64(v), 0)
			}
		}
		badge.Anchor = htmlx.NodeAttrSearch(tr, func(attr htmlx.Attribute) bool { return attr.Key == "data-anchor_name" })
		badge.Name = htmlx.NodeAttrSearch(tr, func(attr htmlx.Attribute) bool { return attr.Key == "data-bn" })
		badges[i] = badge
	}

	sort.Slice(badges, func(i, j int) bool { return badges[i].Intimacy > badges[j].Intimacy })

	return badges, nil
}
