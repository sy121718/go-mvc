# web — 前端项目开发约定

项目基于 vue-pure-admin / pure-admin-thin。

## 核心原则

- **优先使用框架组件。** 写代码前先翻 [COMPONENTS_SUMMARY.md](./COMPONENTS_SUMMARY.md)，确认 `ReXxx` 系列没有现成的再自己实现。
- **`.vue` 只做组件拼装，不写业务逻辑。** 把 Element-Plus、`ReXxx`、自定义组件当成乐高积木，`.vue` 只是把它们拼起来的模板。`<script setup>` 里只做三件事：import 组件、import 变量/函数（从 `.ts`/`.tsx` 导入）、在 template 中使用它们。

## 项目概览

- 框架：Vue 3 + TypeScript + Vite
- UI：Element-Plus
- 状态：Pinia
- HTTP：Axios（`src/utils/http/`）
- 国际化：vue-i18n + `locales/*.yaml`
- CSS：TailwindCSS + SCSS
- 包管理：pnpm

```bash
pnpm dev       # 启动
pnpm build     # 构建
pnpm typecheck # 类型检查
```

## src/ 目录分析

### `api/` — API 接口

每个业务模块一个文件，同时导出请求函数和类型。

```
src/api/
├── admin.ts   # 管理员（列表/创建/详情）
├── menu.ts    # 菜单管理
├── user.ts    # 登录/用户信息
├── routes.ts  # 动态路由
└── mock.ts    # Mock 接口
```

- 响应类型以 `XxxResult`，请求参数以 `XxxReq` / `XxxParams` 命名
- 视图层从 `@/api/xxx` 导入类型，不自己重复定义

### `components/` — 公共组件

`ReXxx` 命名格式，完整清单见 [COMPONENTS_SUMMARY.md](./COMPONENTS_SUMMARY.md)。

核心组件：

| 组件 | 用途 |
|---|---|
| `ReDialog` | 命令式弹窗。新增/编辑用 `addDialog()` |
| `RePureTableBar` | 表格工具栏。所有列表页表格外层包裹 |
| `ReAuth` / `RePerms` | 按钮/页面级权限控制 |
| `ReSegmented` | 分段选择器 |
| `ReIcon` / `IconSelect` | 图标渲染与选择 |
| `ReCol` | 响应式列布局（代替 el-col） |
| `ReCropper` | 图片裁剪 |
| `ReImageVerify` | 图形验证码 |

### `views/` — 视图模块

```
views/<module>/
├── index.vue              # 页面主体（筛选栏 + 表格 + 操作列）
├── components/            # 弹窗表单子组件
└── utils/
    ├── hook.tsx            # useXxx() 组合式逻辑（核心）
    ├── enums.ts            # 选项常量、状态映射
    ├── rule.ts             # el-form 校验规则
    └── types.ts            # 表单组件 props 类型（按需）
```

- `index.vue` 只做模板展示，业务逻辑全在 `hook.tsx`
- `hook.tsx` 返回 `{ form, loading, columns, dataList, onSearch, resetForm, openDialog, handleDelete }`
- 弹窗表单通过 `ReDialog` + `components/` 内的子组件实现
- 列表页用 `RePureTableBar` 包裹 `pure-table`
- `index.vue` 里直接写筛选栏和表格，不拆组件

### `store/` — 状态管理

Pinia store，按模块拆分：

```
store/modules/
├── user.ts        # 用户登录态、信息
├── app.ts         # 应用级别状态
├── permission.ts  # 路由权限
├── multiTags.ts   # 多标签页
├── settings.ts    # 系统设置
└── epTheme.ts     # 主题
```

- `store/types.ts` — store 全局类型
- `store/utils.ts` — store 工具函数

### `router/` — 路由

```
router/modules/
├── home.ts       # 首页路由
├── error.ts      # 错误页路由
└── remaining.ts  # 剩余路由
```

- 动态路由 + 静态路由
- 路由模块统一放在 `router/modules/`
- 只使用 GET 和 POST

### `layout/` — 布局

存放框架布局组件：

```
layout/
├── index.vue              # 主布局
├── frame.vue              # iframe 嵌入
├── redirect.vue           # 路由重定向
├── types.ts               # 布局类型
├── components/
│   ├── lay-sidebar/       # 侧边栏（纵向/混合/横向导航）
│   ├── lay-navbar/        # 顶部导航栏
│   ├── lay-tag/           # 标签页
│   ├── lay-content/       # 内容区
│   ├── lay-search/        # 全局搜索
│   ├── lay-notice/        # 通知
│   ├── lay-setting/       # 设置面板
│   ├── lay-footer/        # 页脚
│   └── lay-panel/         # 面板
└── hooks/                 # 布局相关 hooks
    ├── useNav.ts
    ├── useTag.ts
    ├── useLayout.ts
    └── ...
```

- 布局组件不直接放业务逻辑
- 业务页面通过 `<router-view>` 在内容区渲染

### `utils/` — 工具函数

| 模块 | 作用 |
|---|---|
| `http/` | Axios 封装（请求/响应拦截 + 自动刷新 token） |
| `localforage/` | 本地存储封装 |
| `progress/` | NProgress 进度条 |
| `auth.ts` | 登录鉴权（token 存取、登录状态判断） |
| `message.ts` | 消息提示封装 |
| `tree.ts` | 树结构工具（数据转换、查找） |
| `mitt.ts` | 事件总线 |
| `responsive.ts` | 响应式判断 |
| `propTypes.ts` | Vue 组件 prop 类型定义 |
| `globalPolyfills.ts` | 全局 polyfill |

### `directives/` — 自定义指令

| 指令 | 用途 |
|---|---|
| `auth` | 按钮级权限（v-auth） |
| `perms` | 权限标识判断 |
| `copy` | 复制到剪贴板 |
| `longpress` | 长按事件 |
| `ripple` | 水波纹效果 |
| `optimize` | 渲染优化 |

### `assets/` — 静态资源

```
assets/
├── iconfont/    # iconfont 字体图标
├── login/       # 登录页资源
├── status/      # 状态页（403/404/500 svg）
├── svg/         # 自定义 svg 图标
└── table-bar/   # 表格工具栏图标
```

### `plugins/` — 插件注册

```
plugins/
├── elementPlus.ts  # Element-Plus 全局注册
├── i18n.ts         # vue-i18n 配置
└── echarts.ts      # ECharts 全局注册
```

### `style/` — 全局样式

```
style/
├── index.scss       # 入口
├── tailwind.css     # TailwindCSS
├── element-plus.scss # Element-Plus 样式覆盖
├── dark.scss        # 暗色主题
├── reset.scss       # 样式重置
├── sidebar.scss     # 侧边栏样式
├── theme.scss       # 主题变量
├── transition.scss  # 过渡动画
└── login.css        # 登录页样式
```

### `config/` — 应用配置

`config/index.ts` 存放应用级配置项（请求地址、缓存 key 等）。

### `types/`（web 根目录） — 全局类型声明

```
types/
├── global.d.ts            # 全局类型扩展
├── index.d.ts             # 入口
├── router.d.ts            # 路由类型
├── directives.d.ts        # 指令类型
├── global-components.d.ts # 全局组件类型
├── pure-admin-components.d.ts
├── shims-vue.d.ts         # Vue 模块声明
└── shims-tsx.d.ts         # TSX 声明
```

## 国际化

- 文案 key 写 `locales/zh-CN.yaml` / `locales/en.yaml`
- 模板用 `t("xxx.yyy")`，enums.ts 的 label 用 i18n key，模板中 `transformI18n()` 转换

## 开发步骤

1. 翻 `COMPONENTS_SUMMARY.md` 看框架组件够不够
2. `src/api/` 定义接口和类型
3. 创建 `index.vue` + `utils/`
4. 依次写：`types.ts`(按需) → `rule.ts` → `enums.ts` → `hook.tsx` → `index.vue`
5. 弹窗表单通过 `ReDialog` + `components/xxxForm.vue` 实现