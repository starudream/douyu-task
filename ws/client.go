package ws

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/starudream/go-lib/codec/json"
	"github.com/starudream/go-lib/log"
	"github.com/starudream/go-lib/seq"

	"github.com/starudream/douyu-task/consts"
	"github.com/starudream/douyu-task/internal/cryptox"
)

const (
	// URL websocket url
	// https://www.douyu.com/lapi/live/gateway/web/9999?isH5=1
	URL = "wss://wsproxy.douyu.com:6671"

	loginTimeout = 10 * time.Second
)

type Client struct {
	conn *websocket.Conn
	msg  chan map[string]string
	wMu  *sync.Mutex
	done chan struct{}
}

type LoginParams struct {
	Room     int
	Stk      string
	Ltkid    string
	Username string
}

func Login(p LoginParams) error {
	if p.Room <= 0 {
		p.Room = consts.RoomYYF
	}

	conn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		return err
	}

	client := &Client{
		conn: conn,
		msg:  make(chan map[string]string, 16),
		wMu:  &sync.Mutex{},
		done: make(chan struct{}),
	}

	go client.listen()

	err = client.login(p)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	time.Sleep(time.Second)

	return nil
}

func (c *Client) login(p LoginParams) error {
	devId := seq.UUIDShort()
	rt := strconv.FormatInt(time.Now().Unix(), 10)
	vk := cryptox.MD5Hex(rt + consts.WSLoginHash + devId)

	c.write(
		"type", "loginreq",
		"ver", "20220825", "aver", "218101901",
		"biz", "1", "stk", p.Stk, "ltkid", p.Ltkid, "username", p.Username,
		"roomid", strconv.Itoa(p.Room),
		"devid", devId, "rt", rt, "vk", vk,
	)

	for {
		select {
		case msg := <-c.msg:
			if v, ok := msg["type"]; ok && v == "loginres" {
				if msg["username"] == p.Username {
					log.Info().Msgf("[websocket] login success: %s", json.MustMarshal(msg))
					return nil
				} else {
					return fmt.Errorf("secret keys maybe expired")
				}
			}
		case <-c.done:
			return fmt.Errorf("connection closed")
		case <-time.After(loginTimeout):
			return fmt.Errorf("timeout")
		}
	}
}

func (c *Client) write(kv ...string) {
	if c.conn == nil {
		return
	}
	msg := map[string]string{}
	for i := 0; i < len(kv); i += 2 {
		msg[kv[i]] = kv[i+1]
	}
	c.wMu.Lock()
	defer c.wMu.Unlock()
	err := c.conn.WriteMessage(websocket.BinaryMessage, Encode(kv...))
	if err != nil {
		log.Debug().Msgf("[websocket] write message error: %s", err)
	} else {
		log.Debug().Msgf("[websocket] write message: %s", json.MustMarshal(msg))
	}
}

func (c *Client) listen() {
	if c.conn == nil {
		return
	}
	defer close(c.done)
	for {
		_, bs, err := c.conn.ReadMessage()
		if err != nil {
			log.Debug().Msgf("[websocket] read message error: %s", err)
			break
		}
		msg := Decode(bs)
		log.Debug().Msgf("[websocket] read message: %s", json.MustMarshal(msg))
		c.msg <- msg
	}
}
