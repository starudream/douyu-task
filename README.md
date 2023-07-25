# Douyu-Task

![Golang](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/golang.yml?label=golang&style=for-the-badge)
![Release](https://img.shields.io/github/actions/workflow/status/starudream/douyu-task/release.yml?label=release&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/douyu-task?include_prereleases&sort=semver&style=for-the-badge)
![License](https://img.shields.io/github/license/starudream/douyu-task?style=for-the-badge)

## Config

rename [config.yaml.template](./config.yaml.template) to `config.yaml`.

- `cookie` use [EditThisCookie](https://www.editthiscookie.com/) extension.

  ![cookies.png](./docs/cookies.png)

## Usage

```shell
douyu-task run gift list
douyu-task run gift send --gift-id 268 --room-id 9999 1

douyu-task run badge list

douyu-task run login
```

## Docker

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
    -v $(pwd)/config.yaml:/config.yaml \
    starudream/douyu-task
```

## License

[Apache License 2.0](./LICENSE)
