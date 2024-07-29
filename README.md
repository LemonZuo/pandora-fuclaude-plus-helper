
## docker-compose 部署
```
version: '3.8'
services:
  pandora-helper:
    image: raspberrycheese/pandora-fuclaude-plus-helper
    ports:
      - "8900:5000"
  restart: always
```

## 重要链接
- [Linux.do](https://linux.do)
- [fuclaude](https://github.com/wozulong/fuclaude)
- [PandoraHelper](https://github.com/nianhua99/PandoraHelper)

