
# <p align="center">HiCHat</p>  
  
 ![Static Badge](https://img.shields.io/badge/go-1.21.6-green) 
 
HiChat 是一个基于 Go 语言开发的分布式即时通讯系统。HiChat 致力于为用户提供安全、高效、便捷的即时通讯服务，改善用户间的沟通体验。

![毕设架构图](https://github.com/HryYyR/HiChat/assets/92864176/3e1cd465-be7c-4c59-988d-f311dfec7cb4)

 
## 特性
- 基于WEBRTC的视频聊天
- RSA+AES双重加密保证数据安全
- Redis缓存即时消息，提高消息的获取速度

## 运行环境
- Go
- Redis
- Mysql
- Rabbitmq
- Consul

## 部署与安装
go >= 1.19
1. ``start your nginx``
2. ``start your consul``
3. ``start your rabbitmq``
4. ``git clone https://github.com/HryYyR/HiChat.git``
5. ``./run.sh ``



## 其他
![登录流程图](https://github.com/HryYyR/HiChat/assets/92864176/d521a456-f024-4859-82b5-e157008c8bff)
