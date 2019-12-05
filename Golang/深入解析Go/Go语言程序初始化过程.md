# Go语言程序初始化过程

```go
_rt0_amd64_darwin // 通过参数argc和argv等，确定栈的位置，得到寄存器
main
_rt0_amd64
runtime.check //检测像int8,int16,float等是否是预期的大小，检测cas操作是否正常
runtime.args  //将argc,argv设置到static全局变量中了
runtime.osinit  //osinit做的事情就是设置runtime.ncpu，不同平台实现方式不一样
runtime.hashinit  //将argc,argv设置到static全局变量中了
runtime.schedinit // 内存管理初始化，根据GOMAXPROCS设置使用的procs等等
runtime.newproc 
runtime.mstart
main.main
runtime.exit
// runtime.newproc会把runtime.main放到就绪线程队列里面。本线程继续执行runtime.mstart，m意思是machine。runtime.mstart会调用到调度函数schedule, schedule函数绝不返回，它会根据当前线程队列中线程状态挑选一个来运行。由于当前只有这一个goroutine，它会被调度，然后就到了runtime.main函数中来，runtime.main会调用用户的main函数，即main.main从此进入用户代码。
```

## main.main之前的准备

- sysmon: 在一个新的物理线程中运行sysmon函数, sysmon是一个地位非常高的后台任务，整个函数体一个死循环的形式，目前主要处理两个事件：对于网络的epoll以及抢占式调度的检测, sysmon会根据系统当前的繁忙程度睡一小段时间，然后每隔10ms至少进行一次epoll并唤醒相应的goroutine。同时，它还会检测是否有P长时间处于Psyscall状态或Prunning状态，并进行抢占式调度。
- scavenger：scavenger只是由goroutine运行的，scavenger执行的是runtime·MHeap_Scavenger函数。它将一些不再使用的内存归还给操作系统。Go是一门垃圾回收的语言，垃圾回收会在系统运行过程中被触发，内存会被归还到Go的内存管理系统中，Go的内存管理是基于内存池进行重用的，而这个函数会真正地将内存归还给操作系统。
- main.init，每个包的init函数会在包使用之前先执行。