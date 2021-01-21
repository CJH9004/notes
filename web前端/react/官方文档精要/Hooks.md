# Hooks

## 动机

- Hook 使你在无需修改组件结构的情况下复用状态逻辑

## State Hook

- 通过在函数组件里调用useState来给组件添加一些内部 state。React 会在重复渲染时保留这个 state。
- React 假设当你多次调用 useState 的时候，你能保证每次渲染时它们的调用顺序是不变的。
- todo：调用 State Hook 的更新函数并传入当前的 state 时，React 将跳过子组件的渲染及 effect 的执行。（React 使用 Object.is 比较算法 来比较 state。）

## Effect Hook

-  React 组件中执行数据获取、订阅或者手动修改过 DOM 称为“副作用”
- 可以在 effect 中获取最新的 state 的值，而不用担心其过期的原因是每次我们重新渲染，都会生成新的 effect，替换掉之前的。一旦 effect 的依赖发生变化，它就会被重新创建。
- 如果你传入了一个空数组（[]），effect 内部的 props 和 state 就会一直拥有其初始值。
- 不应在函数中执行阻塞浏览器更新屏幕的操作。
- 虽然 useEffect 会在浏览器绘制后延迟执行，但会保证在任何新的渲染前执行。React 将在组件更新前刷新上一轮渲染的 effect。
- React 会等待浏览器完成画面渲染之后才会延迟调用 useEffect，因此会使得额外操作很方便。

## useContext

- 调用了 useContext 的组件总会在 context 值变化时重新渲染。如果重渲染组件的开销较大，你可以 通过使用 memoization 来优化。

## useReducer

- React 会确保 dispatch 函数的标识是稳定的，并且不会在组件重新渲染时改变。这就是为什么可以安全地从 useEffect 或 useCallback 的依赖列表中省略 dispatch。
- 如果 Reducer Hook 的返回值与当前 state 相同，React 将跳过子组件的渲染及副作用的执行。（React 使用 Object.is 比较算法 来比较 state。）

```js
const initialState = {count: 0};

function reducer(state, action) {
  switch (action.type) {
    case 'increment':
      return {count: state.count + 1};
    case 'decrement':
      return {count: state.count - 1};
    default:
      throw new Error();
  }
}

function Counter() {
  const [state, dispatch] = useReducer(reducer, initialState);
  return (
    <>
      Count: {state.count}
      <button onClick={() => dispatch({type: 'increment'})}>+</button>
      <button onClick={() => dispatch({type: 'decrement'})}>-</button>
    </>
  );
}
```

## useCallback

- 返回一个 memoized 回调函数。
- useCallback(fn, deps) 相当于 useMemo(() => fn, deps)。

```js
const memoizedCallback = useCallback(
  () => {
    doSomething(a, b);
  },
  [a, b],
);
```

## useMemo

- 返回一个 memoized 值。
- 传入 useMemo 的函数会在渲染期间执行。请不要在这个函数内部执行与渲染无关的操作

## useRef

- useRef 返回一个可变的 ref 对象，其 .current 属性被初始化为传入的参数（initialValue）。返回的 ref 对象在组件的整个生命周期内保持不变。
- useRef() 比 ref 属性更有用。它可以很方便地保存任何可变值，其类似于在 class 中使用实例字段的方式
- useRef() 和自建一个 {current: ...} 对象的唯一区别是，useRef 会在每次渲染时返回同一个 ref 对象。
- 当 ref 对象内容发生变化时，useRef 并不会通知你。变更 .current 属性不会引发组件重新渲染。

```js
function TextInputWithFocusButton() {
  const inputEl = useRef(null);
  const onButtonClick = () => {
    // `current` 指向已挂载到 DOM 上的文本输入元素
    inputEl.current.focus();
  };
  return (
    <>
      <input ref={inputEl} type="text" />
      <button onClick={onButtonClick}>Focus the input</button>
    </>
  );
}

function Timer() {
  const intervalRef = useRef();

  useEffect(() => {
    const id = setInterval(() => {
      // ...
    });
    intervalRef.current = id;
    return () => {
      clearInterval(intervalRef.current);
    };
  });

  // ...
}
```

## useDebugValue

```js
function useFriendStatus(friendID) {
  const [isOnline, setIsOnline] = useState(null);

  // ...

  // 在开发者工具中的这个 Hook 旁边显示标签
  // e.g. "FriendStatus: Online"
  useDebugValue(isOnline ? 'Online' : 'Offline');

  return isOnline;
}
```

## 自定义 Hook