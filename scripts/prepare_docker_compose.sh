#!/bin/bash

#
# prepare_docker_compose.sh
#
# 功能: 为 docker-compose 类型的离线包准备文件。
#
# 参数:
#   $1 - STAGING_DIR: 文件的暂存目录
#   $2 - APP_NAME: 应用名称
#   $3 - APP_VERSION: 应用版本
#   $4 - ARCH: 架构 (x86_64, arm)
#

prepare_package() {
  local STAGING_DIR=$1
  local APP_NAME=$2
  local APP_VERSION=$3
  local ARCH=$4

  echo "--- 准备 Docker Compose 包 ---"

  # --- 1. 创建 docker-compose.yml ---
  # 这个示例假设应用包含一个 web 服务和一个数据库
  cat > "$STAGING_DIR/docker-compose.yml" <<EOF
version: '3.8'

services:
  web:
    image: ${APP_NAME}:${APP_VERSION}
    ports:
      - "8080:80"
    volumes:
      - ./config:/app/config
      - ./data/web:/app/data
      - ./logs/web:/app/logs
    restart: always

  database:
    image: postgres:13
    environment:
      POSTGRES_DB: ${APP_NAME}_db
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - ./data/db:/var/lib/postgresql/data
    restart: always
EOF

  echo "已创建 docker-compose.yml"

  # --- 2. 拉取并保存 Docker 镜像 ---
  # 为了演示，我们将创建一个假的 "应用" 镜像，并拉取一个真实的 postgres 镜像

  # 2.1 创建一个假的 Dockerfile 用于演示
  cat > "$STAGING_DIR/Dockerfile.fake" <<EOF
FROM busybox:latest
LABEL app.name="${APP_NAME}"
LABEL app.version="${APP_VERSION}"
CMD ["echo", "这是一个 ${APP_NAME} v${APP_VERSION} 的模拟镜像 for ${ARCH}"]
EOF

  # 2.2 构建并保存应用镜像
  docker build -t "${APP_NAME}:${APP_VERSION}" -f "$STAGING_DIR/Dockerfile.fake" .
  echo "正在保存镜像: ${APP_NAME}:${APP_VERSION}"
  docker save -o "$STAGING_DIR/${APP_NAME}-${APP_VERSION}.tar" "${APP_NAME}:${APP_VERSION}"

  # 2.3 拉取并保存依赖的镜像 (例如 postgres)
  echo "正在拉取并保存镜像: postgres:13"
  docker pull postgres:13
  docker save -o "$STAGING_DIR/postgres-13.tar" "postgres:13"

  # --- 3. 创建必要的目录 ---
  mkdir -p "$STAGING_DIR/config"
  mkdir -p "$STAGING_DIR/data/web"
  mkdir -p "$STAGING_DIR/data/db"
  mkdir -p "$STAGING_DIR/logs/web"

  # --- 4. 创建一个示例配置文件 ---
  cat > "$STAGING_DIR/config/app.conf" <<EOF
# ${APP_NAME} 的示例配置
# 时间: $(date)
LOG_LEVEL=INFO
FEATURE_FLAG_X=true
EOF

  # --- 5. 清理临时文件 ---
  rm "$STAGING_DIR/Dockerfile.fake"

  echo "--- Docker Compose 包准备完成 ---"
}
