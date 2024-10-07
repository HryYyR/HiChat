#!/bin/bash

# 检查是否提供了更新信息
if [ -z "$1" ]; then
  echo "错误: 请提供提交的信息。"
  echo "用法: ./git_push.sh \"提交信息\""
  exit 1
fi

# 获取提交信息
commit_message=$1

# 执行 git add
git add .

# 执行 git commit，并使用提供的提交信息
git commit -m "$commit_message"

# 推送到远程仓库
git push

# 检查是否成功
if [ $? -eq 0 ]; then
  echo "代码提交并推送成功！"
else
  echo "提交或推送失败，请检查！"
fi
