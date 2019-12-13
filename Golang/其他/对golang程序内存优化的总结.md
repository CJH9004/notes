# 对golang程序内存优化的总结

## 1.问题描述

需求是开发一个运行在嵌入式设备上的web程序，内存要求严格。程序包括 gin server 作为后端 + thrift server 接收拎一个程序数据 + 两个http client调用第三方接口传出接收到的数据。由于两个http client的数据字段和thrift字段不一样，所以同一份数据保存了三份，导致一次传输进行的内存分配达到几十兆，程序运行过程中内存一致上涨。

## 2.分析工具

### 1) pprof

```go
import (
  _ "net/http/pprof"
  "net/http"
  "log"
)

go func() {
  log.Println(http.ListenAndServe("0.0.0.0:9999", nil))
}()
```

```sh
# 使用的内存
go tool pprof -inuse_space http://127.0.0.1:6060/debug/pprof/heap?debug=2
# 分配的内存
go tool pprof -alloc_space http://127.0.0.1:6060/debug/pprof/heap?debug=2
```

### 2) GODEBUG

```sh
GODEBUG='gctrace=1' ./main
```

## 3. 原因

go GC的内存不会立即归还给系统，GC有延迟，thrift高并发导致内存高峰

## 4. 解决办法

强制GC和归还内存，减少thrift并发

```golang
import "runtime/debug"

// 强制GC和释放内存
debug.FreeOSMemory()
```

## 5. 总结

### 1) 栈内存和堆内存的区别

**栈（stack）**是由编译器自动分配和释放的一块内存区域，主要用于存放一些基本类型（如int、float等）的变量、指令代码、常量及对象句柄（也就是对象的引用地址）。栈内存的操作方式类似于数据结构中的栈（仅在表尾进行插入或删除操作的线性表）。栈的优势在于，它的存取速度比较快，仅此于寄存器，栈中的数据还可以共享。其缺点表现在，存在栈中的数据大小与生存期必须是确定的，缺乏灵活性。

**堆（heap）**是一个程序运行动态分配的内存区域。堆内存在使用完毕后，是由垃圾回收（Garbage Collection,GC）器“隐式”回收的。堆的优势是在于动态地分配内存大小，可以“按需分配”，其生存期也不必事先告诉编译器，在使用完毕后，垃圾收集器会自动收走这些不再使用的内存块。其缺点为，由于要在运动时才动态分配内存，相比于栈内存，它的存取速度较慢。

### 2) golang中栈内存和堆内存分配原则

在程序中一个数据没有指针指向他，那么使用栈内存，如果有指针指向他，那么用堆内存；slice本身包括一个数组指针和cap、len属性，map都是堆内存；struct是栈内存；函数传值时，栈内存为值传递，堆内存为引用传递。

### 3) golang GC 和内存使用

golang GC使用标记清除，即标记需要清除的数据，在合适的时间清除；golang分配的内存在GC后不会立即归还系统，会在后续分批归还