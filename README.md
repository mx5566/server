#### 分布式框架开发开发中
#### 采用actor模型开发服务器，避免多线程锁问题，通过消息进行交互，所有的状态转化都在entity实体的内部
#### 通过rpc进行远程通信，不需要写函数与消息的绑定关系，自动的去查找对应函数去调用
#### 使用etcd作为服务的注册与发现
#### 增加nats作为服务器之间交互的消息队列
