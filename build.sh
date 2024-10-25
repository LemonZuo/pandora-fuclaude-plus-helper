#!/bin/bash

# 记录构建开始时间
start_time=$(date +%s)

# 获取最新的标签名
TAG=$(git describe --tags --abbrev=0)

# 检查是否获取到标签名
if [ -z "$TAG" ]; then
  echo "No Git tags found."
  exit 1
fi

# 转义标签名中的特殊字符
ESCAPED_TAG=$(printf '%s\n' "$TAG" | sed 's/[\/&]/\\&/g')

# 获取操作系统名称
OS_NAME=$(uname)

# 提前判断操作系统类型，不是 Linux 或 macOS 提前退出脚本
if [ "$OS_NAME" != "Darwin" ] && [ "$OS_NAME" != "Linux" ]; then
    # 其他操作系统，提示不支持
    echo "Unsupported operating system: $OS_NAME"
    exit 1
fi

# 更新 .env 文件中的 VERSION 值
if [ -f ".env" ]; then
    if [ "$OS_NAME" = "Darwin" ]; then
        # macOS 使用 BSD sed
        sed -i '' "s/^VERSION=.*/VERSION=$ESCAPED_TAG/" .env
    elif [ "$OS_NAME" = "Linux" ]; then
        # Linux 使用 GNU sed
        sed -i "s/^VERSION=.*/VERSION=$ESCAPED_TAG/" .env
    fi
else
    echo ".env file not found"
    exit 1
fi

# 从 .env 文件中导入环境变量
export $(cat .env | sed 's/#.*//g' | xargs)


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

# 更新 VERSION 文件
echo ${VERSION} > VERSION

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

# 恢复 VERSION 文件
git restore VERSION

# 记录构建结束时间
end_time=$(date +%s)
# 计算总耗时
elapsed_time=$(( end_time - start_time ))

# 判断运行的操作系统，分别采用兼容的日期命令
if [ "$OS_NAME" = "Darwin" ]; then
    # Mac OS 使用的date命令
    start_fmt=$(date -r $start_time '+%Y-%m-%d %H:%M:%S')
    end_fmt=$(date -r $end_time '+%Y-%m-%d %H:%M:%S')
elif [ "$OS_NAME" = "Linux" ]; then
    # Linux 使用的date命令
    start_fmt=$(date -d @$start_time '+%Y-%m-%d %H:%M:%S')
    end_fmt=$(date -d @$end_time '+%Y-%m-%d %H:%M:%S')
fi

echo "Build started at: $start_fmt"
echo "Build finished at: $end_fmt"
echo "Total build time: $elapsed_time seconds"