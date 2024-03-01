**音视频通讯的具体实现**



本功能采用SFU (Selective Forwarding Unit）架构

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/4/24/16a4ed4d53fb54a1~tplv-t2oaga2asx-jj-mark:3024:0:0:0:q75.png)

服务器在此处的作用主要是转发数据流,当然也可以对数据进行一些操作



*步骤*

1. 通过 navigator.mediaDevices.getUserMedia打开摄像头

2. new MediaRecorder记录数据

3. 监听ondataavailable拿到stream,并设置数据格式 (video/webm),通过ws.send发送至服务器

4. 服务器转发数据流

5. 接收客户端将数据流重新转化为video

   通过mseAPI中的MediaSource , MediaSource , addSourceBuffer() , appendBuffer() 将视频实时的延长以达到即时的效果 , 这也是国内直播平台普遍使用的方法



扩展:

1. 通过MediaRecorder得到的数据只能是单一数据流:只能为音频或者视频,否则不能直接转化为video的src使用,需要分离数据.但是我们可以创建多个MediaRecorder同时记录各自不同数据,这又会遇到另一个问题:

2. 音画不同步和客户显示不同步,两条数据流合并必然会因为网络问题,导致其中有落后数据,需要通过进度追赶的技术手段解决.

3. 一条websocket只能传输一个MediaRecorder的数据流,否则需要数据流分流问题.如果同时开启多条websocket来传输数据呢?

   情况会有点类似webrtc的mesh架构和Mixer架构的结合,区别在于路由架构的主要负载在服务器这边的数据流转发,而mesh结构主要是客户自己负载流量的接收,需要要求客户有良好的网络情况,否则当同时在线的人数过多时一定会出现网络卡顿的情况,而SFU架构取决于服务器的网络带宽.

4. Simulcast模式让客户发送不同分辨率的数据流,其他客户可以根据自己的网络情况进行切换







其他架构:

MCU（Multipoint Conferencing Unit）:由一个服务器和多个终端组成一个[星形结构](https://www.zhihu.com/search?q=星形结构&search_source=Entity&hybrid_search_source=Entity&hybrid_search_extra={"sourceType"%3A"answer"%2C"sourceId"%3A2450913782})。接收每个共享端的音视频流,经过合并后发送给房间里的所有人。

Mesh :即多个终端之间两两进行连接，形成一个网状结构



参考:

https://developer.mozilla.org/zh-CN/docs/Web/API/MediaSource



https://developer.mozilla.org/zh-CN/docs/Web/API/SourceBuffer



https://www.w3schools.com/tags/ref_av_dom.asp

