
# <p align="center">HiCHat</p>  

 # <p align="center">![Static Badge](https://img.shields.io/badge/go-1.21.6-green) </p>

# <p align="center">![icon256](https://github.com/user-attachments/assets/f3baae05-2335-40e8-9f6d-fabcd5447395)</p>  

HiChat 是一个基于 Go 语言开发的分布式即时通讯系统。HiChat 致力于为用户提供安全、高效、便捷的即时通讯服务，改善用户间的沟通体验。

## 客户端
### 源码地址：[HiChatClient](https://github.com/HryYyR/HiChatClient)
### 下载地址：http://203.195.163.23/

## 特性
- 基于非对称和对称双重加密算法实现消息的加密通信以及消息鉴权
- 基于Redis 缓存 聊天记录、用户群聊数据等信息、提高数据响应时间、降低数据库压力
- 基于注册中心实现掉线重连机制、服务器宕机时自动选择合适服务器
- 基于WEBRTC 实现 1V1 的视频聊天功能

## 运行环境
- Go
- Redis
- Mysql
- Rabbitmq
- Consul
- NebulaGraph(可选)

## 部署与安装

### 1.容器部署(推荐)
1. ``docker-compose up -d --build``

### 2.手动安装
go >= 1.19
1. ``start your nginx``
2. ``start your consul``
3. ``start your rabbitmq``
4. ``git clone https://github.com/HryYyR/HiChat.git``
5. ``./run.sh ``

## 联系作者

QQ：2452719312@qq.com

## 其他

![系统架构图](https://github.com/HryYyR/HiChat/assets/92864176/3e1cd465-be7c-4c59-988d-f311dfec7cb4)

![登录流程图](https://github.com/HryYyR/HiChat/assets/92864176/d521a456-f024-4859-82b5-e157008c8bff)

![分布式高可用架构图](https://github.com/HryYyR/HiChat/assets/92864176/1db3de2a-348c-4e71-a093-dd525e429121)

