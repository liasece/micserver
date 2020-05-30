# MICServer

![](https://github.com/liasece/micserver/workflows/Go/badge.svg)

micserver 是一个为分布式系统设计的服务器框架，以模块为服务的基本单位。底层模块间使用 TCP 及自定义二进制编码格式通信，底层实现优先考虑时间效率。

使用 ROC(Remote Object Call)远程对象调用作为模块间的主要通信接口，不关心服务模块本身，而是将上层业务作为对象注册到 micserver 中，寻址路由均由底层维护，只需要知道目标调用对象的 ID 即可调用，不需要关心目标所在的模块。得益于 ROC 抽象，使所有业务状态都可以与服务本身解耦，可以轻松将各个 ROC 对象在模块间转移或加载，实现分布式系统的**无状态**/**热更**/**容灾冗余**等特性。

你可以在[示例程序](https://github.com/liasece/micchaos)中了解 micserver 的基本使用方法。

目前 micserver 不需要任何第三方包。

# 适用于谁?

本框架设计初衷，是个人开发者或者低成本团队中能以最简单步骤建立一个能用于生产的服务器集群。

至少适用于以下业务场景:

1. 服务端主动驱动业务，如游戏服务器，互联网等应用中的主动推送业务，主动式监控。
2. 物联网服务。

# 安装

    go get github.com/liasece/micserver

# 官方文档

[GoDoc](https://godoc.org/github.com/liasece/micserver)
