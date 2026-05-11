# vue-pure-admin 组件库总结

> 基于 `pure-admin/vue-pure-admin` 框架，所有组件以 `Re` 前缀命名，共约 25 个组件。
https://github.com/pure-admin/vue-pure-admin/tree/main/src/components
下载到本地的地址：/home/sky/project/raw/vue-pure-admin





---

## 一、组件索引

| 组件 | 类型 | 用途 |
|------|------|------|
| ReCol | 布局 | 封装 el-col 统一响应式 |
| ReText | 基础 | 文本组件 |
| ReFlicker | 效果 | 闪烁动画 |
| ReFlop | 效果 | 数字翻牌器 |
| ReCountTo | 效果 | 数字滚动动画（normal / rebound） |
| ReTypeit | 效果 | 打字机效果 |
| ReIcon | 图标 | 统一图标体系（iconfont / svg / iconify） |
| ReAuth | 权限 | 按钮级权限控制 |
| RePerms | 权限 | 页面级权限控制 |
| ReImageVerify | 表单 | 图形验证码（SVG） |
| ReBarcode | 表单 | 条形码生成 |
| ReQrcode | 表单 | 二维码生成 |
| ReCropper | 表单 | 图片裁剪 |
| ReCropperPreview | 表单 | 裁剪结果预览 |
| ReDialog | 弹窗 | 命令式弹框（addDialog / closeDialog） |
| ReDrawer | 弹窗 | 命令式抽屉 |
| ReSelector | 选择 | 高级选择器 |
| ReSegmented | 选择 | 分段选择器 |
| ReSplitPane | 布局 | 分割面板（可拖拽） |
| ReSeamlessScroll | 滚动 | 无缝滚动 |
| ReAnimateSelector | 动画 | 动画选择器 |
| ReMap | 地图 | 高德地图封装 |
| RePureTableBar | 表格 | 表格工具栏（列设置/刷新/密度/全屏） |
| ReVxeTableBar | 表格 | VxeTable 工具栏 |
| ReFilterForm | 表单 | 筛选表单容器 |
| ReTreeLine | 树 | 树形连接线 |

---

## 二、组件详解

### 1. ReCol — 响应式列布局

封装 `el-col`，统一设置响应式断点。

```vue
<ReCol :value="12">
  <el-form-item label="名称">
    <el-input />
  </el-form-item>
</ReCol>
```

`value` 默认 24，等价于设置 xs/sm/md/lg/xl 全部为同一值。

---

### 2. ReIcon — 统一图标体系

核心图标方案，支持三种来源，通过 `useRenderIcon(icon)` 自动判别：

| 来源 | 格式 | 示例 |
|------|------|------|
| iconify 在线 | 含 `:` | `"ri/user-3-fill"` |
| iconify 离线 | 对象 | `import User from "~icons/ri/user-3-fill"` |
| iconfont | `"IF-xxx"` | `"IF-iconfont icon-class"` |
| 自定义 SVG | 组件函数 | SVG 组件 |

```vue
<!-- 三种用法 -->
<el-input :prefix-icon="useRenderIcon(User)" />
<IconifyIconOffline icon="ep:check" />
<IconifyIconOnline icon="ri:user-3-fill" />
```

---

### 3. ReAuth — 按钮级权限

配合 `hasAuth()` 函数控制元素显隐，无权限时不渲染。

```vue
<Auth :value="['btn:add']">
  <el-button>新增</el-button>
</Auth>
```

- `value` 传权限标识数组
- 内部调用 `hasAuth()` 判断

---

### 4. RePerms — 页面级权限

与 ReAuth 类似，用于页面级权限判断。

---

### 5. ReDialog — 命令式弹框

核心的弹框管理方案，通过状态数组管理多个弹框。

```typescript
import { addDialog, closeDialog, updateDialog, closeAllDialog } from "@/components/ReDialog";

// 打开弹框
addDialog({
  title: "提示",
  contentRenderer: () => h("div", "内容"),
  beforeSure: (done, { closeLoading }) => {
    // 异步操作
    done();
  }
});

// 关闭所有
closeAllDialog();
```

**DialogOptions 完整配置：**
- `props` — 内容区组件的 props（defineProps 接收）
- `headerRenderer` / `contentRenderer` / `footerRenderer` — 自定义渲染器
- `footerButtons` — 自定义按钮组
- `beforeCancel` / `beforeSure` — 取消/确定前回调（可异步阻塞关闭）
- `open` / `close` / `closeCallBack` — 生命周期回调
- `sureBtnLoading` — 确定按钮 loading
- `fullscreenIcon` — 全屏切换按钮

---

### 6. ReDrawer — 命令式抽屉

用法与 ReDialog 类似，通过命令式 API 控制抽屉的打开/关闭/更新。

---

### 7. RePureTableBar — 表格工具栏

表格头部操作栏，提供统一功能：

```vue
<RePureTableBar title="用户列表" :columns="columns" @refresh="getList">
  <template #buttons>
    <el-button>导出</el-button>
  </template>
  <template #default="{ size, dynamicColumns }">
    <el-table :size="size" :columns="dynamicColumns" />
  </template>
</RePureTableBar>
```

**功能清单：**
- 刷新（loading 动画）
- 密度切换（宽松/默认/紧凑）
- 列显隐控制（Checkbox 列表）
- 列拖拽排序（SortableJS）
- 全屏切换
- 树形表格展开折叠
- 插槽自定义按钮区

---

### 8. ReImageVerify — 图形验证码

SVG 图形验证码组件，支持后端 API 获取验证码并渲染为 SVG 图片。

```vue
<ReImageVerify v-model:code="captchaCode" :width="120" :height="40" />
```

- 自动调用后端 `GET /captcha` 接口
- 点击刷新验证码
- 暴露 `captchaKey`（验证码 ID）和 `getImgCode` 刷新方法
- 内置干扰线+噪点+旋转字体

---

### 9. ReCropper — 图片裁剪

基于图片裁剪库，提供裁剪功能。

```vue
<ReCropper v-model:file="file" :width="200" :height="200" />
```

---

### 10. ReCountTo / ReFlop — 数字动画

| 组件 | 用途 |
|------|------|
| ReCountTo (normal) | 普通数字滚动动画 |
| ReCountTo (rebound) | 回弹式数字动画 |
| ReFlop | 翻牌器效果 |

---

### 11. ReTypeit — 打字机效果

基于 TypeIt 库，逐字输出文本。

---

### 12. ReSplitPane — 分割面板

可拖拽调整左右/上下区域大小的面板组件。

---

### 13. ReSeamlessScroll — 无缝滚动

列表或表格数据的无缝循环滚动。

---

### 14. ReBarcode / ReQrcode

| 组件 | 用途 |
|------|------|
| ReBarcode | 条形码生成 |
| ReQrcode | 二维码生成（带中间 logo） |

---

### 15. ReSelector — 高级选择器

封装 el-select，支持更多定制选项。

---

### 16. ReSegmented — 分段选择器

Segment 风格的选择器，替代 Radio 组或 Tab。

---

### 17. ReFilterForm — 筛选表单

封装搜索/筛选表单布局，通常放在表格上方。

---

### 18. ReMap — 高德地图

高德地图组件封装，使用：

```vue
<ReMap :center="[116.4, 39.9]" :zoom="12" />
```

---

### 19. ReText — 文本组件

用于文本展示，支持动态样式绑定。

---

### 20. ReTreeLine — 树形连接线

为树形结构添加可视化连接线。

---

### 21. ReAnimateSelector — 动画选择器

选择入场/离场动画的组件。

---

## 三、组件设计模式

### 3.1 目录结构

```
ReXxx/
├── index.ts          # 入口：导出 + withInstall 注册
├── index.vue         # 组件实现（或用 TSX）
├── type.ts           # 类型定义（可选）
├── src/
│   └── ...           # 子模块
└── README.md         # 文档（可选）
```

### 3.2 注册方式

使用 `withInstall` 工具进行组件注册，支持全局/按需引入。

```typescript
// index.ts
import reDialog from "./index.vue";
import { withInstall } from "@pureadmin/utils";
const ReDialog = withInstall(reDialog);
export { ReDialog };
```

### 3.3 命令式 API

ReDialog、ReDrawer 采用命令式 API，通过全局状态数组管理实例：

```
addDialog(options)    → 打开弹框
closeDialog(...)      → 关闭弹框
updateDialog(...)     → 更新弹框属性
closeAllDialog()      → 关闭所有
```

### 3.4 权限体系

```
ReAuth           → 按钮级：v-if 控制 DOM 显隐
RePerms          → 页面级
hasAuth(value)   → 工具函数：判断是否有权限
hasPerms(value)  → 按钮级：判断按钮权限标识
```

### 3.5 图标体系

```
useRenderIcon(icon) → 自动路由到对应图标组件
IconifyIconOffline  → 离线 iconify 图标
IconifyIconOnline   → 在线 iconify 图标
FontIcon            → iconfont 图标
```

---

## 四、引用

- 官方文档：[https://pure-admin.cn](https://pure-admin.cn)
- GitHub：[https://github.com/pure-admin/vue-pure-admin](https://github.com/pure-admin/vue-pure-admin)