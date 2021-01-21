# golang的操作符==、!=

## 一般比较的是值

```go
// interface, struct
var a interface{} = struct{ string }{"asdf"}
var b interface{} = struct{ string }{"asdf"}
println(a == b, struct{}{} == struct{}{}) // true, true

// string
c := "asdf"
d := "asdfd"
println(c == d[0:4]) // 即使字符串底层为一个*[]byte，仍然是true

// array
println([1]{1} == [1]{1}) // true

// 指针的值是地址，所以比较的是地址
println(&struct{}{} != &struct{}{}) // true
```

## 切片、map、func只能和nil比较

## 与nil比较

```go
// interface，map，slice，pointer，func，chan默认初始化为nil
var a interface{}
var b map[string]string
var c []int
var d *int
var e func()
var f chan int
println(a == nil, b == nil, c == nil, d == nil, e == nil, f == nil) // true...
```

```go
// interface外使用= nil、:= nil初始化仍为nil, 但interface{}被赋予一个值为nil的变量后不为nil
var a interface{}
var b *string
a = b
c := b
println(a != nil, c == nil)
println(a, b) // (0x45f380,0x0) 0x0, 因为interface{}保存了类型信息
```