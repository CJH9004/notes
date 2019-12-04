# 使用typescript和storybook开发react组件库

## 1.脚手架

### 1.初始化

1. 创建git repository
2. `npm init`
3. `mkdir src`，src作为源码路径

### 2.安装和配置react, typescript

1. `yarn add typescript react @types/react -D`
2. 创建`tsconfig.json`, outDir为输出路径

```json
{
  "compilerOptions": {
    "sourceMap": true,
    "target": "es6",
    "jsx": "react",
    "declaration": true,
    "outDir": "lib",
    "moduleResolution": "node",
  },
  "include": [
    "src"
  ]
}
```

3. 修改package.json

```json
{
  //...
  "scripts": {
    //...
    "build": "tsc -p .",
    "start": "npm run build -- -w",
    //...
  },
  // ...
}
```

### 3.安装和配置storybook, jest

1. install

```sh
npx -p @storybook/cli sb init --type react
yarn add -D awesome-typescript-loader
yarn add -D @storybook/addon-info react-docgen-typescript-loader # optional but recommended
yarn add -D jest "@types/jest" ts-jest #testing
```

2. 在.storybook目录下增加webpack.config.js文件

```js
module.exports = ({ config }) => {
config.module.rules.push({
  test: /\.(ts|tsx)$/,
  use: [
    {
      loader: require.resolve('awesome-typescript-loader'),
    },
    // Optional
    {
      loader: require.resolve('react-docgen-typescript-loader'),
    },
  ],
});
config.resolve.extensions.push('.ts', '.tsx');
return config;
};
```

3. 更改.storybook/config.js, 使用addon-info

```js
import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';

addDecorator(withInfo); 

// automatically import all files ending in *.stories.js
configure(require.context('../stories', true, /\.stories\.js$/), module);
```

4. `yarn ts-jest config:init`


5. 修改package.json

```json
{
  //...
  "scripts": {
    //...
    "test": "jest",
    //...
  },
  // ...
}
```

### 4. 使用storybook插件

1. 安装notes 和 source

```sh
yarn add @storybook/addon-notes @storybook/addon-storysource
```

2. 更改.storybook/addons.js

```js
import '@storybook/addon-storysource/register';
import '@storybook/addon-notes/register-panel';
```

3. 更改.storybook/webpack.config.js，配置addon source

```js
module.exports = ({ config }) => {
  config.module.rules.push({
    test: /\.(ts|tsx)$/,
    use: [
      {
        loader: require.resolve('awesome-typescript-loader'),
      },
      // Optional
      {
        loader: require.resolve('react-docgen-typescript-loader'),
      },
    ],
  });
  config.resolve.extensions.push('.ts', '.tsx');

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
  config.module.rules.push({
    test: /\.stories\.js?$/,
    loaders: [
      {
        loader: require.resolve('@storybook/source-loader'),
        options: { parser: 'typescript' },
      },
    ],
    enforce: 'pre',
  });
// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

  return config;
};
```

### 5.使用typedoc

1. `yarn add -D typedoc typedoc-plugin-markdown`
2. 修改package.json

```json
{
  //...
  "scripts": {
    //...
    "doc": "npx typedoc --plugin typedoc-plugin-markdown --out typedocs src"
    //...
  },
  // ...
}
```

## 2. 开发

### 1. src下除index.ts的每个文件或目录代表一个组件

### 2. src/index.ts 导出所有组件

### 3. 每个组件有自己的测试目录__test__

## 3. 测试

### 1. 测试自定义hooks

1. `yarn add -D @testing-library/react-hooks react-test-renderer`
2. https://react-hooks-testing-library.com/usage/basic-hooks