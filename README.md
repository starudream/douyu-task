# Douyu-Task

![Golang](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/golang.yml?label=golang&style=for-the-badge)
![Release](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/release.yml?label=release&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/douyu-task?include_prereleases&sort=semver&style=for-the-badge)
![License](https://img.shields.io/github/license/starudream/douyu-task?style=for-the-badge)

## Config

### Full

```yaml
debug: false
startup: false
log:
  level: debug
douyu:
  room: 9999
  did: xxx
  stk: xxx
  uid: xxx
  ltkid: xxx
  username: xxx
  ltp0: xxx
  stick:
    remaining: 9999
dingtalk:
  token: xxx
  secret: xxx
telegram:
  token: xxx
  chat_id: xxx
cron:
  refresh: "0 0 1 * * *"
  send_gift: "0 0 12 * * 0"
```

### Global

| Name      | Type   | Comment                                        |
|-----------|--------|------------------------------------------------|
| debug     | bool   | show verbose information                       |
| log.level | string | available: DEBUG INFO WARN ERROR FATAL DISABLE |

### Douyu

#### Send Gift（通过 http 接口送免费礼物）

| Name            | Type   | Require | Comment                             | Cookie                    |
|-----------------|--------|---------|-------------------------------------|---------------------------|
| did             | string | T       |                                     | dy_did (douyu.com)        |
| uid             | string | T       |                                     | acf_uid (douyu.com)       |
| auth            | string | T/F     | expire in 7 days                    | acf_auth (douyu.com)      |
| ltp0            | string | T/F     | to refresh auth token               | LTP0 (passport.douyu.com) |
| stick.remaining | int    | F       | room id to send remaining free gift |                           |

#### Daily Refresh（通过 websocket 获取每日荧光棒）

| Name     | Type   | Require | Comment | Cookie                   |
|----------|--------|---------|---------|--------------------------|
| stk      | string | T       |         | acf_stk (douyu.com)      |
| ltkid    | string | T       |         | acf_ltkid (douyu.com)    |
| username | string | T       |         | acf_username (douyu.com) |

- `chrome://settings/cookies/detail?site=douyu.com`
- `chrome://settings/cookies/detail?site=passport.douyu.com`

### Dingtalk

| Name   | Type   | Comment               |
|--------|--------|-----------------------|
| token  | string | access token          |
| secret | string | secret (not required) |

### Telegram

| Name    | Type   | Comment |
|---------|--------|---------|
| token   | string | token   |
| chat_id | string |         |

### Docker

![Version](https://img.shields.io/docker/v/starudream/douyu-task?sort=semver&style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/douyu-task?sort=semver&style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/douyu-task?style=for-the-badge)

```bash
docker pull starudream/douyu-task
```

```bash
docker run -d \
    --name douyu-task \
    --restart always \
    -v /opt/docker/douyu/config.yaml:/config.yaml \
    -e CONFIG_PATH=/config.yaml \
    starudream/douyu-task
```

## License

[Apache License 2.0](./LICENSE)
