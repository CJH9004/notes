# react性能优化

## 给组件的props传入jsx

```js
<Layout 
  header={<HeaderEle props1 />}
  content={<ContentEle props2 />}
/>
// props1 或 props2中的一个改变都导致Layout重render，生成了不同的Layout元素(header和content都发生了变化)，进而导致headerEle和contentEle重render

// 使用useMemo
const headerEle = useMemo(<HeaderEle props1 />, [props1])
const contentEle = useMemo(<ContentEle props2 />, [props2])

<Layout 
  header={headerEle}
  content={contentEle}
/>
// 尽管prop1导致Layout重render，但content不变，所以contentEle不会重render
```

## 正确使用context

```js
const Context = createContext({state: initState, setState: () => console.log("can not change default value")});

// 导出useProvider而不是Context，避免Context.Consumer导致的forceUpdate
export const useProvider = () => {
  return useContext(FilterContext)
}

// Context.Provider
const Provider: React.FC<Props> = ({children}) => {
  const [state, setState] = useState(initState)
  // 使用useMemo避免Provider render而导致value变化
  const value = useMemo(() => ({ state, setState }), [state]);
  return (
    <Context.Provider value={value} >
      {children}
    </Context.Provider>
  )
}

export default Provider
```