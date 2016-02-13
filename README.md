# pombridge

![](https://avatars1.githubusercontent.com/u/2346051?v=3&s=120)

简介
===
一个简单的 tcp/udp 隧道,客户端监听一个本地端口,并通过多条tcp链路和服务端
连接.当本地尝试与客户端监听的端口连接的时候,服务端对应的与远端的某个本地端口
连接.客户端与服务端的多个tcp连接可以横跨不同的网络,一个简单的例子是,两台
机器通过两张百兆有线网卡连接,如果各在一张网卡上建立一个客户端到服务端的tcp
链路,则我们可以代理让客户端到服务端的单条tcp连接达到两百兆的传送速度.

下面是一张结构图

```
TCP Conn1   |                          | BusLine1
----------> | Client   SendBus      ---|---------->
   ...      | listen  ------------>    |              INTERNET
TCP ConnN   | port     Split into      | BusLine2
----------> |          Messages     ---|---------->


BusLine1    |                           | TCP Conn1  |
----------> | Server    RecvBus       --|----------->| Server
            | listen  -------------->   |     ...    | dst
BusLine2    | ports     FlowControl     | TCP ConnN  | port
----------> |                         --|----------->|

```
