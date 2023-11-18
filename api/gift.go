package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/tablew/v2"

	"github.com/starudream/douyu-task/consts"
)

type Gifts struct {
	List []*Gift `json:"list"`
}

func (gs *Gifts) TableString() string {
	if len(gs.List) == 0 {
		return "empty gifts"
	}
	return tablew.Structs(gs.List)
}

func (gs *Gifts) Find(id int) *Gift {
	for i := range gs.List {
		if gs.List[i].Id == id {
			return gs.List[i]
		}
	}
	return nil
}

func (gs *Gifts) NotEmpty(ids ...int) int {
	for i := 0; i < len(ids); i++ {
		if gs.Find(ids[i]).GetCount() > 0 {
			return ids[i]
		}
	}
	return -1
}

// Gift 礼物
type Gift struct {
	Id       int    `json:"id"`       // id
	Name     string `json:"name"`     // 名称
	Count    int    `json:"count"`    // 现有数量
	Exp      int    `json:"exp"`      // 经验
	Intimate int    `json:"intimate"` // 亲密度
	Met      int    `json:"met"`      // 过期时间

	Price     int `json:"price"     table:",ignored"` // 价值
	PriceType int `json:"priceType" table:",ignored"` // 价值类型（不确定）2-免费礼物
	PropType  int `json:"propType"  table:",ignored"` // 礼物类型（不确定）2-免费礼物 5-等级礼包 6-分区喇叭
}

func (gift *Gift) GetCount() int {
	if gift == nil {
		return 0
	}
	return gift.Count
}

func (gift *Gift) TodayExpired() bool {
	if gift == nil {
		return false
	}
	y, m, d := time.Now().Date()
	t := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
	return t.Unix() >= int64(gift.Met)
}

type SendGiftResp struct {
	CommonResp
	Data *Gifts `json:"data"`
}

func (c *Client) SendGift(room, gift, count int) (*Gifts, error) {
	roomId := strconv.Itoa(room)

	resp, err := c.R().
		SetCookies(c.genAuthCookies()).
		SetHeader("referer", URL+"/"+roomId).
		SetFormData(gh.MS{
			"roomId":    roomId,
			"propId":    strconv.Itoa(gift),
			"propCount": strconv.Itoa(count),
		}).
		SetResult(&SendGiftResp{}).
		Post(URL + "/japi/prop/donate/mainsite/v1")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status())
	}

	res, ok := resp.Result().(*SendGiftResp)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}

	if res.Error != 0 {
		return nil, fmt.Errorf("douyu error(%d): %s", res.Error, res.Msg)
	}

	return res.Data, nil
}

type ListGiftResp struct {
	CommonResp
	Data *Gifts `json:"data"`
}

func (c *Client) ListGifts(rooms ...int) (*Gifts, error) {
	room := func() int {
		if len(rooms) > 0 && rooms[0] != 0 {
			return rooms[0]
		}
		return consts.RoomYYF
	}()

	resp, err := c.R().
		SetCookies(c.genAuthCookies()).
		SetHeader("referer", URL).
		SetResult(&ListGiftResp{}).
		SetQueryParam("rid", strconv.Itoa(room)).
		Get(URL + "/japi/prop/backpack/web/v1")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http status: %s", resp.Status())
	}

	res, ok := resp.Result().(*ListGiftResp)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}

	if res.Error != 0 {
		return nil, fmt.Errorf("douyu error(%d): %s", res.Error, res.Msg)
	}

	return res.Data, nil
}
