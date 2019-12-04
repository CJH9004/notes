# XSS攻击

## 原理

XSS攻击基本上是基于js脚本的（在flash中是ActionScript），利用页面上的bug输入一些js或sql脚本并执行（获取cookie，发送数据，删库等）

## 防范

- 对所有用户输入的内容转义
- 设置cookie httponly