
## docker-compose 部署
```
version: '3.8'
services:
  pandora-fuclaude-plus-helper:
    image: raspberrycheese/pandora-fuclaude-plus-helper:latest
    container_name: pandora-plus-helper
    restart: unless-stopped
    ports:
      - "8182:5000"
    environment:
      - TZ=Asia/Shanghai
      # 管理员密码
      - ADMIN_PASSWORD=123456
      # oaifree站点地址
      - SHARE_TOKEN_AUTH=https://new.oaifree.com
      # fuclaude站点地址
      - FUCLAUDE_LOGIN_AUTH=https://demo.fuclaude.com
    volumes:
      - ./data:/data
```

## 重要链接
- [Linux.do](https://linux.do)
- [fuclaude](https://github.com/wozulong/fuclaude)
- [PandoraHelper](https://github.com/nianhua99/PandoraHelper)

