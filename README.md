# Open Tree Hole Backend

Backend of OpenTreeHole ---- Anonymous BBS for college students

## Features

- White-listed registration ---- for certain community like college students
- Anonymous: RSA encrypted personal information(email) and random identity
- Compliance: report, mute, ban, fold NSFW contents
- Push notifications: web(websocket), iOS and Android
- Balance between performance and development efficiency: stress test 300~400 qps

## Install

This installation is just for backend program. If you want to deploy the whole OpenTreeHole project, please visit [Deploy Repo](https://github.com/OpenTreeHole/deploy).

This project continuously integrates with docker. Go check it out if you don't have docker locally installed.

```shell
docker run -d -p 80:80 shi2002/open_tree_hole_backend
```

Note: this project is runnable with zero-configurations, for full configuration, visit [Configuration Doc](https://github.com/OpenTreeHole/deploy/wiki/配置文档).

## Usage

Please refer to [API Docs](https://github.com/OpenTreeHole/backend/wiki/API-文档).

## Badge

[![build](https://github.com/OpenTreeHole/backend/actions/workflows/master.yaml/badge.svg)](https://github.com/OpenTreeHole/backend/actions/workflows/master.yaml)
[![dev build](https://github.com/OpenTreeHole/backend/actions/workflows/dev.yaml/badge.svg)](https://github.com/OpenTreeHole/backend/actions/workflows/dev.yaml)

[![stars](https://img.shields.io/github/stars/OpenTreeHole/backend)](https://github.com/OpenTreeHole/backend/stargazers)
[![issues](https://img.shields.io/github/issues/OpenTreeHole/backend)](https://github.com/OpenTreeHole/backend/issues)
[![pull requests](https://img.shields.io/github/issues-pr/OpenTreeHole/backend)](https://github.com/OpenTreeHole/backend/pulls)

[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

### Powered by

![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)
![Django](https://img.shields.io/badge/django-%23092E20.svg?style=for-the-badge&logo=django&logoColor=white)
![DjangoREST](https://img.shields.io/badge/DJANGO-REST-ff1709?style=for-the-badge&logo=django&logoColor=white&color=ff1709&labelColor=gray)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/OpenTreeHole/backend/issues/new) or submit PRs.

Full contributing docs and requirements is available at [wiki](https://github.com/OpenTreeHole/backend/wiki/开发).

### Contributors

This project exists thanks to all the people who contribute.

<a href="https://github.com/OpenTreeHole/backend/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenTreeHole/backend" />
</a>

## Licence

[![license](https://img.shields.io/github/license/OpenTreeHole/backend)](https://github.com/OpenTreeHole/backend/blob/dev/LICENSE)
© OpenTreeHole
