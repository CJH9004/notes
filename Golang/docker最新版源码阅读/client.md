# client

client包括dockerd程序和docker程序

## client

> docker/cli/cli
使用cobra库解析命令行参数，调用http api。
client主体为struct DockerCli，包含一个api client和一些基本信息获取函数

### http api

> moby/client
http api 主体为struct Client，将各种http请求封装成方法，包含一个http.client

