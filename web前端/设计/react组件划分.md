# 组件划分

## 组件分类

### 1.容器组件

- 容器组件从获取、存储数据，维护一个较大状态，来源于父组件的数据较少
- 例如：一个页面组件，初始化时从后端获取整个页面所需的或者多个子组件依赖的数据，并维护一个状态来提供需要多个子组件的交互

```js
function Page(){
  const data = useApi("some url")

  return (
    <Framework>
      <Functions data={data} />
    </Framework>
  )
}
```

### 2.框架组件

- 框架组件用于布局、容纳其他组件，是较为通用的样式组件，包含的业务较少，自身数据与外界关联较少
- 例如：Card, Row, Col, Modal

```js
<Card>
  <Row>
    <Func1 data={data} />
  </Row>
  <Row>
    <Func2 data={data} />
  </Row>
</Card>
```

### 3.功能组件

- 功能组件用于实现一个功能、业务或交互
- 例如：Button, Input, Form

## 自定义hook分类

> 使用自定义hook实现UI与逻辑分离

### 1.数据hook

- 对useState，useReducer，useContext，useRef的拓展，提供对数据的操作

```js
import * as React from 'react';

interface state {
  current: number, 
}

interface action {
  type: string,
  payload: {total: number, to?: number}
}

const initialState: state = {current: 0};

function reducer(state: state, action: action) {
  const {total} = action.payload
  const current = state.current

  switch (action.type) {
    case 'prev':
      if(current > 0){
        return {current: current - 1};
      }
      return {current: total - 1}
    case 'next':
      if(current < total - 1){
        return {current: current + 1};
      }
      return {current: 0}
    case 'jumpTo':
      const to = action.payload.to
      if(to >= 0 && to < total){
        return {current: to};
      }
      return state
    default:
      throw new Error("invalid type");
  }
}

export default function useRange(total: number): [number, () => void, () => void, (to: number) => void] {
  if(total < 1){
    total = 1
  }

  const [{current}, dispatch] = React.useReducer(reducer, initialState);
  
  const prev = React.useCallback(() => {
    dispatch({type: "prev", payload: {total}})
  }, [total])

  const next = React.useCallback(() => {
    dispatch({type: "next", payload: {total}})
  }, [total])

  const jumpTo = React.useCallback((to: number) => {
    dispatch({type: "jumpTo", payload: {total, to}})
  }, [total])

  return [current, prev, next, jumpTo]
}
```

### 2.副作用hook

- 对useEffect的拓展，自成一体，依赖和返回较少的数据

```js
import * as React from "react"

export default function useKeys(keyMap : {[propName: string]: () => void}){
  React.useEffect(() => {
    const onkeydown = (e) => {
      if(keyMap[e.code]){
        keyMap[e.code]()
      }
    }

    window.addEventListener("keydown", onkeydown)

    return () => {
      window.removeEventListener("keydown", onkeydown)
    }
  }, (<any>Object).values(keyMap))
}
```

```js
import * as React from "react"

export default function useLoop(action: () => void, gap: number) {
  const id: {current: number} = React.useRef();

  React.useEffect(() => {
    id.current = setInterval(() => {
      action()
    }, gap)

    return () => {
      clearInterval(id.current)
    }
  }, [action, gap])
}
```

### 3.组合hook

- 综合所有原生hooks的复杂自定义hooks
- 一般可以拆分