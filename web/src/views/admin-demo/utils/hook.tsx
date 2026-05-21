/**
 * admin-demo 页面的 hook（核心组合式逻辑）
 *
 * ═══════════════════════════════════════════════════════════
 * 学习重点：hook 是 .vue 大脑，.vue 只做组件拼装。
 * ═══════════════════════════════════════════════════════════
 *
 * 这个文件是整个页面的"业务逻辑中心"。
 * index.vue 只负责：import 这个 hook → 解构返回值 → 在 template 里绑数据和方法。
 *
 * 整体模式：
 *   useAdmin() 返回 { 响应式数据, 方法 }
 *   → index.vue 解构后直接用在 el-form/el-table/el-pagination 上
 */

// ────── Vue 组合式 API ──────
// reactive: 对象深响应式（form 里的字段变更会触发模板重渲染）
// ref:      基本类型/数组响应式（.value 读写）
// onMounted: 组件挂载后执行一次（首次加载数据）
import { reactive, ref, onMounted, computed } from "vue";
// 国际化（本项目用 vue-i18n），模板里用 t("xxx.yyy") 显示对应语言文案
import { useI18n } from "vue-i18n";
// pure-admin 工具：判断是否"全空"（空字符串/undefined/null/空数组 都算空）
import { isAllEmpty } from "@pureadmin/utils";
// 封装好的 Element-Plus ElMessage，传 message(text, { type }) 即可
import { message } from "@/utils/message";
// 从 src/api/ 导入接口函数和类型
import { getAdminList } from "@/api/admin";
import type { AdminListReq } from "@/api/admin";
// 从本模块的 enums.ts 导入常量/映射函数
import { getAdminStatusTagType, getAdminStatusLabel } from "./enums";

export function useAdmin() {
  const { t } = useI18n();

  // ────── form：搜索表单数据（绑定到模板 <el-form :model="form">）──────
  // reactive 包裹，form.xxx 直接修改即可触发响应
  // 类型 AdminListReq 来自 src/api/admin.ts，保持前后端一致
  const form: AdminListReq = reactive({
    page: 1,            // 当前页码 → 绑定 el-pagination v-model:current-page
    limit: 10,          // 每页条数 → 绑定 el-pagination v-model:page-size
    email: "",          // 搜索条件：邮箱
    name: "",           // 搜索条件：姓名
    phone: "",          // 搜索条件：手机号
    status: undefined,  // 搜索条件：状态（下拉框，undefined 表示"不限"）
    sort_field: undefined,  // 排序列名 → handleSortChange 赋值
    sort_order: undefined   // 排序方向 "asc" / "desc"
  });

  // ────── 响应式状态 ──────
  const formRef = ref();         // 模板上 ref="formRef"，用来调 resetFields()
  const dataList = ref([]);      // 表格数据，请求成功后赋值 → 模板 :data="dataList"
  const total = ref(0);          // 总记录数 → 分页组件 :total="total"
  const loading = ref(true);     // 加载中 → 模板 :loading="loading"，el-table 显示骨架

  // ────── 表格列定义 ──────
  // TableColumnList 是 pure-admin 提供的全局类型，不需要 import
  // 模板中通过 <RePureTableBar> 传递 columns → 自动渲染 el-table-column
  const columns: TableColumnList = [
    // sortable: "custom" 表示"手动排序"，点表头触发 @sort-change 事件
    // 排序逻辑在 handleSortChange 里处理
    {
      label: "ID",
      prop: "id",
      width: 80,
      sortable: "custom"          // 点表头排序，数据由后端排序
    },
    {
      label: "用户名",
      prop: "username",
      minWidth: 120
    },
    {
      label: "姓名",
      prop: "name",
      minWidth: 100,
      formatter: (_: any, __: any, value: string) => value || "-"
    },
    {
      label: "邮箱",
      prop: "email",
      minWidth: 160,
      formatter: (_: any, __: any, value: string) => value || "-"
    },
    {
      label: "手机号",
      prop: "phone",
      minWidth: 130,
      formatter: (_: any, __: any, value: string) => value || "-"
    },
    {
      label: "状态",
      prop: "status",
      width: 120,
      // cellRenderer 是 TSX 渲染函数，替代 slot
      // getAdminStatusTagType / getAdminStatusLabel 从 enums.ts 读取
      cellRenderer: ({ row, props }) => (
        <el-tag size={props.size} type={getAdminStatusTagType(row.status)} effect="plain">
          {getAdminStatusLabel(row.status)}
        </el-tag>
      )
    },
    {
      label: "超管",
      prop: "is_admin",
      width: 80,
      formatter: (_: any, __: any, value: number) => (value === 1 ? "是" : "否")
    },
    {
      label: "创建时间",
      prop: "create_time",
      minWidth: 170,
      sortable: "custom"
    },
    {
      // slot: "operation" → 模板中通过 v-if="col.slot === 'operation'" 自定义渲染
      // 这里只占位，实际按钮在 index.vue 的 template 里手写
      label: "操作",
      fixed: "right",
      width: 180,
      slot: "operation"
    }
  ];

  // ────── 核心方法 ──────

  /**
   * 搜索/刷新
   * 1. 组装 params（只传非空字段，减少请求体）
   * 2. 调 getAdminList 发 GET 请求
   * 3. 成功后塞 dataList + total，失败弹错误提示
   * 4. 延迟 500ms 关 loading（避免闪烁）
   */
  async function onSearch() {
    loading.value = true;
    try {
      // 必传的分页参数直接取 form
      const params: AdminListReq = {
        page: form.page,
        limit: form.limit
      };

      // 可选参数：非空才传，让后端统一处理"传了才筛"
      if (!isAllEmpty(form.email)) params.email = form.email;
      if (!isAllEmpty(form.name)) params.name = form.name;
      if (!isAllEmpty(form.phone)) params.phone = form.phone;
      if (form.status !== undefined && form.status !== null) params.status = form.status;
      if (!isAllEmpty(form.sort_field)) params.sort_field = form.sort_field;
      if (!isAllEmpty(form.sort_order)) params.sort_order = form.sort_order;

      // 发请求，类型 AdminListResp 已在 src/api/admin.ts 定义
      const res = await getAdminList(params);
      // 兼容 code 0 和 200（不同后端风格）
      if (res.code === 0 || res.code === 200) {
        dataList.value = res.data.list || [];
        total.value = res.data.total || 0;
      } else {
        message(res.message , { type: "error" });
      }
    } catch (error) {
      console.error("获取管理员列表失败:", error);
    }
    setTimeout(() => {
      loading.value = false;
    }, 500);
  }

  /**
   * 重置搜索条件
   * formEl 从模板传入：resetForm(formRef)
   * 1. 手动重置所有字段
   * 2. 调 el-form 的 resetFields() 清空校验状态
   * 3. 重新搜索
   */
  function resetForm(formEl) {
    if (!formEl) return;
    form.page = 1;
    form.limit = 10;
    form.email = "";
    form.name = "";
    form.phone = "";
    form.status = undefined;
    form.sort_field = undefined;
    form.sort_order = undefined;
    formEl.resetFields();
    onSearch();
  }

  /** 切换每页条数（el-pagination @size-change） */
  function handleSizeChange(val: number) {
    form.limit = val;
    form.page = 1;   // 切换条数后回到第一页
    onSearch();
  }

  /** 切换页码（el-pagination @current-change） */
  function handleCurrentChange(val: number) {
    form.page = val;
    onSearch();
  }

  /**
   * 表头排序（el-table @sort-change）
   * prop: 列名，order: "ascending" / "descending" / null
   * 转为后端识别的 sort_field + sort_order
   * 取消排序时清空，交给后端默认排序
   */
  function handleSortChange({ prop, order }) {
    if (order) {
      form.sort_field = prop;
      form.sort_order = order === "ascending" ? "asc" : "desc";
    } else {
      form.sort_field = undefined;
      form.sort_order = undefined;
    }
    form.page = 1;   // 排序变更也回到第一页
    onSearch();
  }

  // 页面加载时自动触发首次查询，vue框架自带的页面初始化
  onMounted(() => {
    onSearch();
  });

  // ────── 返回给 index.vue 使用 ──────
  // 模板中绑定：form → el-form :model，dataList → el-table :data，loading → :loading
  // 方法绑定：onSearch → 搜索按钮 @click，resetForm → 重置按钮 @click
  return {
    form,
    formRef,
    loading,
    columns,
    dataList,
    total,
    onSearch,
    resetForm,
    handleSizeChange,
    handleCurrentChange,
    handleSortChange
  };
}