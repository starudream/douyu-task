# Douyu-Task

![Golang](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/golang.yml?label=golang&style=for-the-badge)
![Docker](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/docker.yml?label=docker&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/douyu-task?include_prereleases&sort=semver&style=for-the-badge)
![License](https://img.shields.io/github/license/starudream/douyu-task?style=for-the-badge)

## Config

### Full

```yaml
debug: false
log:
  level: debug
douyu:
  did: abc
  uid: 123
  auth: xyz
  stick:
    remaining: 9999
dingtalk:
  token: qwe
  secret: zxc
telegram:
  token: okm
  chat_id: 123
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

| Name            | Type   | Require | Comment                             | Cookie                    |
|-----------------|--------|---------|-------------------------------------|---------------------------|
| did             | string | T       |                                     | dy_did (douyu.com)        |
| uid             | string | T       |                                     | acf_uid (douyu.com)       |
| auth            | string | T/F     | expire in 7 days                    | acf_auth (douyu.com)      |
| ltp0            | string | T/F     | to refresh auth token               | LTP0 (passport.douyu.com) |
| stick.remaining | int    | F       | room id to send remaining free gift |                           |

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
