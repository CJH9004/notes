# golang的泛型函数

## sort包使用纯接口

## 使用类型断言

```go
func add (a interface{}, b interface{}) interface{} {
  switch a.(type) {
    // ...
  }
}

add(1, 2).(int)
```

## 封装类型断言

```go
package main

import (
	"errors"
	"fmt"
)

type AddFunc func(a interface{}, b interface{}) interface{}

func (f AddFunc) Int(a int, b int) int {
	return f(a, b).(int)
}

func (f AddFunc) Float64(a float64, b float64) float64 {
	return f(a, b).(float64)
}

var Add AddFunc = func(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case int:
		return a.(int) + b.(int)
	case float64:
		return a.(float64) + b.(float64)
	default:
		panic(errors.New("not support type"))
	}
}

func main() {
  // 避免使用者做类型断言
  Add.Int(1, 2)
  Add.Float64(1.2, 1.2)
}
```

## 按类型实现多个函数

```go
func addInt(a int, b int) int {
  return a + b
}
func addFloat64(a float64, b float64) float64 {
  return a + b
}

addInt(1,2)
addFloat64(1.2,1.3)
```