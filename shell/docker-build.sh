#!/bin/bash

# 手动定义每个项目的路径
PROJECT_DIRS=(
    "../hichat-ws-service"
    "../hichat-static-service"
    "../hichat-file-service"
    "../hichat-mq-service"
    "../hichat-streammedia-service"
)

# Docker镜像的仓库前缀，例如：yourusername 或 yourregistry.com/yourusername
REPO_PREFIX="hyyyh"

# 遍历手动定义的项目路径
for PROJECT_DIR in "${PROJECT_DIRS[@]}"; do
    # 提取项目名称（文件夹名称）
    PROJECT_NAME=$(basename "$PROJECT_DIR")

    # 构建Docker镜像
    echo "Building Docker image for project: $PROJECT_NAME"
    docker build -t "$REPO_PREFIX/$PROJECT_NAME:latest" "$PROJECT_DIR"
    
    # 检查构建是否成功
    if [ $? -eq 0 ]; then
    else
        echo "Failed to build image for project: $PROJECT_NAME"
    fi
done
