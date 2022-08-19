# Douyu-Task

![Golang](https://img.shields.io/github/workflow/status/starudream/douyu-task/Golang/master?style=for-the-badge)
![Docker](https://img.shields.io/github/workflow/status/starudream/douyu-task/Docker/master?label=Docker&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/douyu-task?include_prereleases&style=for-the-badge)
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
```

### Global

| Name      | Type   | Comment                                        |
|-----------|--------|------------------------------------------------|
| debug     | bool   | show verbose information                       |
| log.level | string | available: DEBUG INFO WARN ERROR FATAL DISABLE |

### Douyu

| Name            | Type   | Require | Comment                             |
|-----------------|--------|---------|-------------------------------------|
| did             | string | T       | cookie, dy_did                      |
| uid             | string | T       | cookie, acf_uid                     |
| auth            | string | T       | cookie, acf_auth                    |
| cookie          | string | T       | cookie                              |
| stick.remaining | int    | F       | room id to send remaining free gift |

### Docker

![Version](https://img.shields.io/docker/v/starudream/douyu-task?style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/douyu-task/latest?style=for-the-badge)
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
