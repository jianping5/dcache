# 分布式缓存
## 项目介绍
简易的分布式缓存项目

## 项目要点
1. 为解决内存资源限制的问题，实现了 LRU 缓存淘汰算法；
2. 引入 sync.Mutex 互斥锁，实现 LRU 缓存的并发控制；
3. 为解决节点扩展问题，实现一致性哈希算法，并引入虚拟节点，解决数据倾斜的问题；
4. 实现 HTTP 服务端，并创建 HTTP 客户端，实现多节点间的通信；
5. 为解决缓存击穿问题，实现 singleflight，将并发请求合并成一个请求，以减少对下层服务的压力；
6. 使用 protobuf 库，在传输前后对消息进行编解码，优化节点间通信的性能。
   
## 参考
[https://github.com/geektutu/7days-golang](https://github.com/geektutu/7days-golang)