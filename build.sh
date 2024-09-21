#!/bin/bash

# 记录构建开始时间
start_time=$(date +%s)

# 从 .env 文件中导入环境变量
if [ -f ".env" ]; then
    export $(cat .env | sed 's/#.*//g' | xargs)
else
    echo ".env file not found"
    exit 1
fi

# 使用环境变量中的用户名和密码尝试登录Docker Hub
docker login -u="${HUB_USER}" -p="${HUB_PASS}"
status=$?

# 检查登录命令的退出状态
if [ $status -ne 0 ]; then
    echo "Docker login failed, exiting..."
    exit $status
else
    echo "Docker login successful."
fi

# 创建并使用一个新的 Buildx 构建器实例，如果已存在则使用现有的
BUILDER_NAME=multi-platform-build
docker buildx create --name ${BUILDER_NAME} --use || true
docker buildx use ${BUILDER_NAME}
docker buildx inspect --bootstrap

# 使用 Docker Buildx 构建镜像，同时标记为 latest 和 VERSION，支持多架构
docker buildx build \
  --no-cache \
  --platform linux/amd64,linux/arm64 \
  -t ${HUB_USER}/${HUB_REPO}:${VERSION} \
  -t ${HUB_USER}/${HUB_REPO} . \
  --push \
  --progress=plain

# 登出 Docker Hub
docker logout

# 记录构建结束时间
end_time=$(date +%s)
# 计算总耗时
elapsed_time=$(( end_time - start_time ))

# 判断运行的操作系统，分别采用兼容的日期命令
if [[ "$(uname)" == "Darwin" ]]; then
    # Mac OS 使用的date命令
    start_fmt=$(date -r $start_time '+%Y-%m-%d %H:%M:%S')
    end_fmt=$(date -r $end_time '+%Y-%m-%d %H:%M:%S')
else
    # Linux 使用的date命令
    start_fmt=$(date -d @$start_time '+%Y-%m-%d %H:%M:%S')
    end_fmt=$(date -d @$end_time '+%Y-%m-%d %H:%M:%S')
fi

echo "Build started at: $start_fmt"
echo "Build finished at: $end_fmt"
echo "Total build time: $elapsed_time seconds"