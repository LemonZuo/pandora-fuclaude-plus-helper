
## docker-compose 部署
```
services:
  helper:
    image: raspberrycheese/pandora-fuclaude-plus-helper:latest
    container_name: helper
    restart: unless-stopped
    environment:
      # 时区
      - TZ=Asia/Shanghai
      # API端口
      - HTTP_PORT=8181
      # 反代OPENAI端口
      - OPENAI_PORT=8182
      # 反代CLAUDE端口
      - CLAUDE_PORT=8183
      # 数据库驱动，可选sqlite、mysql，默认sqlite
      - DATABASE_DRIVER=mysql
      # 数据库DSN，sqlite时注释，mysql时必填
      - DATABASE_DSN=****:********@tcp(127.0.0.1:3306)/db_pandora_plus_helper?parseTime=true&loc=Asia%2FShanghai
      # 管理员密码
      - ADMIN_PASSWORD=**************
      # 后台登录加密密钥
      - SECRET=**********************************
      # 反代OPENAI地址，默认https://new.oaifree.com
      - OPENAI_SITE=https://new.oaifree.com
      # OPENAI跳转地址，默认https://new.oaifree.com,如需修改访问地址需要nginx反向代理到OPENAI_PORT
      - OPENAI_AUTH_SITE=https://new.oaifree.com
      # 反代CLAUDE地址，默认https://demo.fuclaude.com, fuclaude站点地址,可以是内网地址，确保容器之间可以通信
      - CLAUDE_SITE=https://demo.fuclaude.com
      # CLAUDE跳转地址，默认https://demo.fuclaude.com,如需修改访问地址需要nginx反向代理到CLAUDE_PORT
      - CLAUDE_AUTH_SITE=https://demo.fuclaude.com
      # 内容审查地址，不开启内容审查时注释
      - MODERATION_ENDPOINT=https://api.openai.com
      # 内容审查API密钥，不开启内容审查时注释
      - MODERATION_API_KEY=sk-********************
      # 内容审查消息提示
      - MODERATION_MESSAGE=***********************
      # 是否隐藏openai邮箱信息，默认false
      - HIDDEN_USER_INFO=false
      # 是否开启定时刷新，默认true
      - ENABLE_TASK=true
    volumes:
      # 数据驱动为sqlite时，数据存储位置
      - ./data:/data
```

## 重要链接
- [Linux.do](https://linux.do)
- [Fuclaude](https://github.com/wozulong/fuclaude)
- [PandoraHelper](https://github.com/nianhua99/PandoraHelper)

