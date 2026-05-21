import { reactive, ref, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";
import { isAllEmpty } from "@pureadmin/utils";
import { message } from "@/utils/message";
//提交请求
import type { AdminListReq } from "@/api/admin";
// 获取请求
import { getAdminList } from "@/api/admin";
//获取定义的枚举
import { getAdminStatusTagType, getAdminStatusLabel } from "./enums";

//export 用于导出到外部了，use开头表示导出，admin约定的list列表或者主数据
export function useAdmin() {
  //先定义空容器响应式
  const formRef = ref();         // 空 → 模板 ref="formRef" 赋值
  const dataList = ref([]);      // 空数组 → 请求后塞数据
  const total = ref(0);          // 0 → 请求后更新总量
  const loading = ref(true);     // true → 初始加载态，请求完变 false
  //使用国际化
  const { t } = useI18n();
  //给请求接口默认挂几个参数
  const form: AdminListReq = reactive({
    page: 1,
    limit: 10,
    email: "",
    name: "",
    phone: "",
    status: undefined,
    sort_field: undefined,
    sort_order: undefined
  });

  const columns: TableColumnList = [

    {
      label: "ID",
      prop: "id",
      minWidth: 80,
      align: "center",
      // 或固定在左侧
      fixed: "left",
      //允许排序
      sortable: "custom",
    }, {
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
      label: "操作",
      fixed: "right",
      width: 180,
      slot: "operation"
    }

  ];
  async function onSearch() {
    // 每次刷新加载动画
    loading.value = true;
    try {
      const params: AdminListReq = {
        page: form.page,
        limit: form.limit
      }
      // 如果form的值不是空的，那就赋值给params用于发起请求
      if (!isAllEmpty(form.email)) params.email = form.email;
      if (!isAllEmpty(form.name)) params.name = form.name;
      if (!isAllEmpty(form.phone)) params.phone = form.phone;
      if (form.status !== undefined && form.status !== null) params.status = form.status;
      if (!isAllEmpty(form.sort_field)) params.sort_field = form.sort_field;
      if (!isAllEmpty(form.sort_order)) params.sort_order = form.sort_order;
      //发起请求并且用res存储响应的值
      const res = await getAdminList(params);
      if (res.code === 200) {
        // 类似ajax的sussces
        dataList.value = res.data.list || [];
        total.value = res.data.total || 0;
      } else {
        message(res.message, { type: "error" })

      }

    } catch (error) {
      //ajax的error
      console.error("获取管理员列表失败", error)
    }


  }

  /**
   * 重置搜索
   * 1.手动
   * 2.调用resetFields
   */
  function resetForm(formEl) {
    if (!formEl) return;
    // 这里的form是手动清理
    form.page = 1;
    form.limit = 10;
    form.email = "";
    form.name = "";
    form.phone = "";
    form.status = undefined;
    form.sort_field = undefined;
    form.sort_order = undefined;
    //resetfields是elemelt-plus的默认方法，把所有的字段恢复到初始值。用来补充全清兜底
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
      form,//表单字段
      formRef,//表单id
      loading,//加载动画
      columns,//列表
      dataList,//存储表格数据
      total,//总数，分页组件的
      onSearch,//搜索+刷新
      resetForm,//重置
      handleSizeChange,//limit条数
      handleCurrentChange,//分页
      handleSortChange//排序
    };




}

