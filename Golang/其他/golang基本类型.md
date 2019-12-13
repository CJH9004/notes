# golang基本类型

## 数组，string，slice

- string存储的是一个字节数组指针和长度

```go
//  reflect.StringHeader
type StringHeader struct {
  Data uintptr 
  Len int
}
```

- []byte(string) 不会产生运行时开销, 可以直接便利字节数组，for range []byte(string)
- []rune 是 []int32 的别名，[]rune 和 string底层数据类型[]byte不一致，所以转换可能引起开销：O(n)的时间复杂度和内存分配
- 数组定义[...]int{1: 1, 0: 0, 2}
- 不同长度的数组和数组指针是不同类型
- for range比fori快，var times [5][0]int的大小是0，可通过for range times快速循环5次，[0]int可用作不占内存的值，如：c <- [0]int{}, <- c, 这里并不关心c的值，不过更倾向于使用空的匿名结构体：c <- struct{}{}
- slice的实现

```go
// reflect.SliceHeader
type SliceHeader struct { 
  Data uintptr 
  Len int 
  Cap int 
}
```

- slice：array[start : end : cap]
- 切片可以和 nil 进行比较，只有当切片底层数据指针为空时切 片本身为 nil ，这时候切片的长度和容量信息将是无效的。
- 在切片首部添加元素, 避免生成中间数组

```go
a = append(a, 0) // length + 1
copy(a[i+1:], a[i:]) // >> 1
a[i] = x  // set a[i]

// insert slice x
a = append(a, x...) // length + 1
copy(a[i+len(x):], a[i:]) // >> len x
copy(a[i:], x)
```