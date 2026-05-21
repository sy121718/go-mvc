import editForm from "../form.vue";
import { useI18n } from "vue-i18n";
import { message } from "@/utils/message";
import { addDialog } from "@/components/ReDialog";
import { reactive, ref, onMounted, h, computed } from "vue";
import type { FormItemProps } from "../utils/types";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { cloneDeep, isAllEmpty, deviceDetection } from "@pureadmin/utils";
import {
  getMenuTree,
  getFullMenuTree,
  getMenuDetail,
  createMenu,
  updateMenu,
  deleteMenu
} from "@/api/system/menu";
import type { MenuTreeItem, MenuQueryParams } from "@/api/system/menu";
import { hasAuth } from "@/router/utils";
import { hasPerms } from "@/utils/auth";

export function useMenu() {
  const { t } = useI18n();
  const form = reactive<MenuQueryParams>({
    search: "",
    parent_id: undefined,
    type: undefined,
    status: undefined,
    is_public: undefined,
    create_time_min: "",
    create_time_max: "",
    update_time_min: "",
    update_time_max: "",
    id: "desc" // 默认按ID降序
  });

  const formRef = ref();
  const dataList = ref([]);
  const loading = ref(true);

  // 计算是否显示操作列：只要有任一操作权限就显示
  const showOperationColumn = computed(() => {
    return (
      (hasAuth("menu_update") && hasPerms("menu:update")) ||
      (hasAuth("menu_create") && hasPerms("menu:create")) ||
      (hasAuth("menu_batchDelete") && hasPerms("menu:batchDelete"))
    );
  });

  const getMenuType = (type: number, text = false) => {
    switch (type) {
      case 1:
        return text ? t("menuManagement.directory") : "primary";
      case 2:
        return text ? t("menuManagement.menu") : "success";
      case 3:
        return text ? t("menuManagement.link") : "danger";
      case 4:
        return text ? t("menuManagement.iframe") : "warning";
      case 5:
        return text ? t("menuManagement.button") : "info";
      default:
        return text ? t("menuManagement.unknown") : "";
    }
  };

  const columns: TableColumnList = [
    {
      label: t("menuManagement.menuName"),
      prop: "menu_name",
      align: "left",
      minWidth: 200,
      cellRenderer: ({ row }) => (
        <>
          <span class="inline-block mr-1">
            {h(useRenderIcon(row.icon), {
              style: { paddingTop: "1px" }
            })}
          </span>
          <span>{row.menu_name}</span>
        </>
      )
    },
    {
      label: t("menuManagement.menuRemark"),
      prop: "title",
      minWidth: 150
    },
    {
      label: t("menuManagement.menuCode"),
      prop: "menu_code",
      minWidth: 150
    },
    {
      label: t("menuManagement.menuType"),
      prop: "type",
      width: 100,
      cellRenderer: ({ row, props }) => (
        <el-tag size={props.size} type={getMenuType(row.type)} effect="plain">
          {getMenuType(row.type, true)}
        </el-tag>
      )
    },
    {
      label: t("menuManagement.routePath"),
      prop: "path",
      minWidth: 150
    },
    {
      label: t("menuManagement.componentPath"),
      prop: "component",
      minWidth: 180,
      formatter: ({ component }) => component || "-"
    },
    {
      label: t("menuManagement.sort"),
      prop: "sort_order",
      width: 80,
      sortable: "custom"
    },
    {
      label: t("menuManagement.status"),
      prop: "status",
      width: 80,
      cellRenderer: ({ row, props }) => (
        <el-tag
          size={props.size}
          type={row.status === 1 ? "success" : "danger"}
          effect="plain"
        >
          {row.status === 1
            ? t("menuManagement.enabled")
            : t("menuManagement.stopped")}
        </el-tag>
      )
    },
    {
      label: t("menuManagement.hidden"),
      prop: "is_hidden",
      formatter: ({ is_hidden }) =>
        is_hidden === 1 ? t("menuManagement.yes") : t("menuManagement.no"),
      width: 80
    },
    {
      label: t("menuManagement.public"),
      prop: "is_public",
      cellRenderer: ({ row, props }) => (
        <el-tag
          size={props.size}
          type={row.is_public === 1 ? "warning" : "info"}
          effect="plain"
        >
          {row.is_public === 1
            ? t("menuManagement.yes")
            : t("menuManagement.no")}
        </el-tag>
      ),
      width: 80
    },
    {
      label: t("menuManagement.createTime"),
      prop: "create_time",
      minWidth: 160,
      sortable: "custom"
    },
    {
      label: t("menuManagement.updateTime"),
      prop: "update_time",
      minWidth: 160,
      sortable: "custom"
    },
    {
      label: t("menuManagement.operation"),
      fixed: "right",
      width: 210,
      slot: "operation",
      hide: !showOperationColumn.value
    }
  ];

  function handleSelectionChange(val) {
    console.log("handleSelectionChange", val);
  }

  function resetForm(formEl) {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  }

  async function onSearch() {
    loading.value = true;
    try {
      const params: MenuQueryParams = {};

      // 搜索关键词
      if (!isAllEmpty(form.search)) {
        params.search = form.search;
      }

      // 父级菜单
      if (form.parent_id !== undefined && form.parent_id !== null) {
        params.parent_id = form.parent_id;
      }

      // 菜单类型
      if (form.type !== undefined && form.type !== null) {
        params.type = form.type;
      }

      // 状态
      if (form.status !== undefined && form.status !== null) {
        params.status = form.status;
      }

      // 是否公共
      if (form.is_public !== undefined && form.is_public !== null) {
        params.is_public = form.is_public;
      }

      // 创建时间范围
      if (!isAllEmpty(form.create_time_min)) {
        params.create_time_min = form.create_time_min;
      }
      if (!isAllEmpty(form.create_time_max)) {
        params.create_time_max = form.create_time_max;
      }

      // 更新时间范围
      if (!isAllEmpty(form.update_time_min)) {
        params.update_time_min = form.update_time_min;
      }
      if (!isAllEmpty(form.update_time_max)) {
        params.update_time_max = form.update_time_max;
      }

      // 排序
      if (form.id) {
        params.id = form.id;
      } else if (form.sort_order) {
        params.sort_order = form.sort_order;
      } else if (form.create_time) {
        params.create_time = form.create_time;
      } else if (form.update_time) {
        params.update_time = form.update_time;
      }

      const response = await getMenuTree(params);
      if (response.success) {
        dataList.value = response.data;
      } else {
        message(
          `${t("menuManagement.getMenuListFailed")}: ${response.message}`,
          { type: "error" }
        );
      }
    } catch (error) {
      console.error("获取菜单列表失败:", error);
    }

    setTimeout(() => {
      loading.value = false;
    }, 500);
  }

  /** 生成上级菜单选项：排除当前菜单及其子节点，避免形成环 */
  function formatHigherMenuOptions(treeList, excludeId?: number) {
    if (!treeList || !treeList.length) return [];
    const result = [];
    for (let i = 0; i < treeList.length; i++) {
      const node = treeList[i];
      if (excludeId && node.id === excludeId) continue;
      if (node.children && node.children.length) {
        node.children = formatHigherMenuOptions(node.children, excludeId);
      }
      result.push(node);
    }
    return result;
  }

  function normalizeMenuType(type?: number | string) {
    if (typeof type === "number") return type;
    switch (type) {
      case "directory":
        return 1;
      case "menu":
        return 2;
      case "link":
        return 3;
      case "iframe":
        return 4;
      case "button":
        return 5;
      default:
        return 2;
    }
  }

  async function openDialog(title = "新增", row?: MenuTreeItem) {
    // 为选择"上级菜单"单独拉取完整树（不受当前表格筛选影响）
    let higherOptions = [];
    let menuDetail = row;
    try {
      const fullTreeRes = await getFullMenuTree({ status: 1 });
      if (fullTreeRes?.success) {
        higherOptions = formatHigherMenuOptions(
          cloneDeep(fullTreeRes.data || []),
          row?.id
        );
      } else {
        higherOptions = [];
      }
    } catch {
      // 回退到当前表格数据
      higherOptions = formatHigherMenuOptions(
        cloneDeep(dataList.value || []),
        row?.id
      );
    }

    // 确保 higherOptions 始终是一个数组
    if (!Array.isArray(higherOptions)) {
      higherOptions = [];
    }

    if (row?.id) {
      try {
        const detailRes = await getMenuDetail(row.id);
        if (detailRes?.success && detailRes.data) {
          const permissionsFromDetail =
            detailRes.data.permission_ids ??
            detailRes.data.permissions?.map((item: any) => item.id) ??
            [];
          menuDetail = {
            ...detailRes.data,
            type: normalizeMenuType(detailRes.data.type),
            permission_ids: permissionsFromDetail
          };
        } else {
          message(detailRes?.message || "获取菜单失败", { type: "error" });
          return;
        }
      } catch (error) {
        console.error("获取菜单详情失败:", error);
        message("获取菜单详情失败，请稍后重试", { type: "error" });
        return;
      }
    }

    addDialog({
      title: `${title}菜单`,
      props: {
        formInline: {
          id: menuDetail?.id ?? undefined,
          menu_code: menuDetail?.menu_code ?? "",
          menu_name: menuDetail?.menu_name ?? "",
          parent_id: menuDetail?.parent_id ?? 0,
          type: menuDetail?.type ?? 2,
          title: menuDetail?.title ?? "",
          path: menuDetail?.path ?? "",
          component: menuDetail?.component ?? "",
          external_url: menuDetail?.external_url ?? "",
          icon: menuDetail?.icon ?? "",
          status: menuDetail?.status ?? 1,
          is_hidden: menuDetail?.is_hidden ?? 0,
          is_public: menuDetail?.is_public ?? 0,
          sort_order: menuDetail?.sort_order ?? 99,
          permission_ids: menuDetail?.permission_ids ?? [],
          // 上级菜单选项（完整树），排除当前菜单，允许修改或清空上级
          higherMenuOptions: higherOptions
        }
      },
      width: "60%",
      draggable: true,
      fullscreen: deviceDetection(),
      fullscreenIcon: true,
      closeOnClickModal: false,
      contentRenderer: () => h(editForm, { ref: formRef, formInline: null }),
      beforeSure: (done, { options }) => {
        const FormRef = formRef.value.getRef();
        const curData = options.props.formInline as FormItemProps;

        FormRef.validate(async valid => {
          if (valid) {
            try {
              const menuData = {
                menu_code: curData.menu_code,
                menu_name: curData.menu_name,
                parent_id: curData.parent_id || 0,
                type: curData.type,
                title: curData.title,
                path: curData.path,
                component: curData.component,
                external_url: curData.external_url,
                icon: curData.icon,
                status: curData.status,
                is_hidden: curData.is_hidden,
                is_public: curData.is_public,
                sort_order: curData.sort_order,
                permission_ids: curData.permission_ids || []
              };

              let result;
              if (title === "新增") {
                result = await createMenu(menuData);
              } else {
                result = await updateMenu(curData.id!, menuData);
              }

              if (result.success) {
                message("操作成功", { type: "success" });
                done(); // 关闭弹框
                onSearch(); // 刷新表格数据
              } else {
                message(result.message || "操作失败", { type: "error" });
              }
            } catch (error: any) {
              console.error("操作失败:", error);
            }
          }
        });
      }
    });
  }

  async function handleDelete(row) {
    try {
      const result = await deleteMenu([row.id]);
      if (result.success) {
        message("删除成功", { type: "success" });
        onSearch();
      } else {
        message(result.message || "删除失败", { type: "error" });
      }
    } catch (error: any) {
      console.error("删除失败:", error);
    }
  }

  // 处理表格排序
  function handleSortChange({ prop, order }) {
    // 清空所有排序
    form.id = undefined;
    form.sort_order = undefined;
    form.create_time = undefined;
    form.update_time = undefined;

    if (order) {
      const sortValue = order === "ascending" ? "asc" : "desc";

      // 根据点击的列设置对应的排序
      if (prop === "sort_order") {
        form.sort_order = sortValue;
      } else if (prop === "id") {
        form.id = sortValue;
      } else if (prop === "create_time") {
        form.create_time = sortValue;
      } else if (prop === "update_time") {
        form.update_time = sortValue;
      }
    } else {
      // 如果取消排序，恢复默认ID降序
      form.id = "desc";
    }

    onSearch();
  }

  onMounted(() => {
    onSearch();
  });

  return {
    form,
    loading,
    columns,
    dataList,
    /** 搜索 */
    onSearch,
    /** 重置 */
    resetForm,
    /** 新增、修改菜单 */
    openDialog,
    /** 删除菜单 */
    handleDelete,
    handleSelectionChange,
    /** 表格排序 */
    handleSortChange,
    /** 是否显示操作列 */
    showOperationColumn
  };
}
