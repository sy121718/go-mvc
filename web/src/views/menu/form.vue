<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import ReCol from "@/components/ReCol";
import { formRules } from "./utils/rule";
import { FormProps } from "./utils/types";
import { transformI18n } from "@/plugins/i18n";
import { IconSelect } from "@/components/ReIcon";
import Segmented from "@/components/ReSegmented";
import { getPermissionTree } from "@/api/system/permission";
import type { PermissionTreeItem } from "@/api/system/permission";
import {
  menuTypeOptions,
  statusOptions,
  hiddenOptions,
  publicOptions
} from "./utils/enums";

const props = withDefaults(defineProps<FormProps>(), {
  formInline: () => ({
    menu_code: "",
    menu_name: "",
    parent_id: 0,
    type: 2,
    title: "",
    path: "",
    component: "",
    external_url: "",
    icon: "",
    status: 1,
    is_hidden: 0,
    is_public: 0,
    sort_order: 99,
    permission_ids: [],
    higherMenuOptions: []
  })
});

const ruleFormRef = ref();
const newFormInline = ref(props.formInline);
const permissionTreeRef = ref();
const permissionTreeData = ref<PermissionTreeItem[]>([]);
const permissionLoading = ref(false);
const permissionSearchValue = ref("");

// 计算属性：菜单类型索引（Segmented组件使用索引而非值）
const menuTypeIndex = computed({
  get: () => {
    return menuTypeOptions.findIndex(
      item => item.value === newFormInline.value.type
    );
  },
  set: (index: number) => {
    newFormInline.value.type = menuTypeOptions[index]?.value || 2;
  }
});

// 计算属性：判断是否需要显示外部链接字段
const showExternalUrl = computed(() => {
  return newFormInline.value.type === 3 || newFormInline.value.type === 4;
});

// 计算属性：判断是否需要显示组件路径字段（只有菜单需要）
const showComponent = computed(() => {
  return newFormInline.value.type === 2;
});

// 计算属性：判断是否需要显示路由路径字段（菜单、外链、iframe需要）
const showPath = computed(() => {
  return [2, 3, 4].includes(newFormInline.value.type);
});

// 计算属性：状态索引
const statusIndex = computed({
  get: () => {
    return statusOptions.findIndex(
      item => item.value === newFormInline.value.status
    );
  },
  set: (index: number) => {
    newFormInline.value.status = statusOptions[index]?.value ?? 1;
  }
});

// 计算属性：隐藏状态索引
const hiddenIndex = computed({
  get: () => {
    return hiddenOptions.findIndex(
      item => item.value === newFormInline.value.is_hidden
    );
  },
  set: (index: number) => {
    newFormInline.value.is_hidden = hiddenOptions[index]?.value ?? 0;
  }
});

// 计算属性：公共状态索引
const publicIndex = computed({
  get: () => {
    return publicOptions.findIndex(
      item => item.value === newFormInline.value.is_public
    );
  },
  set: (index: number) => {
    newFormInline.value.is_public = publicOptions[index]?.value ?? 0;
  }
});

// 计算属性：判断是否需要显示权限分配（只有菜单和按钮需要权限）
const showPermissionAssign = computed(() => {
  return newFormInline.value.type === 2 || newFormInline.value.type === 5;
});

// 获取权限树数据
async function loadPermissionTree() {
  permissionLoading.value = true;
  try {
    const response = await getPermissionTree();
    if (response.success) {
      permissionTreeData.value = response.data;
    }
  } catch (error) {
    console.error("获取权限树失败:", error);
  } finally {
    permissionLoading.value = false;
  }
}

// 权限树节点过滤
function filterPermissionNode(value: string, data: PermissionTreeItem) {
  if (!value) return true;
  return data.perm_name?.includes(value) || data.perm_code?.includes(value);
}

// 监听搜索值变化
watch(permissionSearchValue, val => {
  permissionTreeRef.value?.filter(val);
});

// 权限树check事件 - 只保存真正选中的节点，不包括半选节点
function handlePermissionCheck() {
  const checkedKeys = permissionTreeRef.value?.getCheckedKeys() || [];

  newFormInline.value.permission_ids = checkedKeys as number[];
}

function getRef() {
  return ruleFormRef.value;
}

onMounted(() => {
  loadPermissionTree();
});

defineExpose({ getRef });
</script>

<template>
  <el-form
    ref="ruleFormRef"
    :model="newFormInline"
    :rules="formRules"
    label-width="100px"
  >
    <el-row :gutter="30" class="menu-form-layout">
      <re-col
        :value="showPermissionAssign ? 13 : 24"
        :xs="24"
        class="menu-form-left"
      >
        <el-row :gutter="30">
          <!-- 菜单类型 -->
          <re-col :value="24">
            <el-form-item :label="transformI18n('menu.pureMenuType')">
              <Segmented
                v-model="menuTypeIndex"
                :options="
                  menuTypeOptions.map(item => ({
                    ...item,
                    label: transformI18n(item.label)
                  }))
                "
              />
            </el-form-item>
          </re-col>

          <!-- 上级菜单 -->
          <re-col :value="24">
            <el-form-item :label="transformI18n('menu.pureParentMenu')">
              <el-cascader
                v-model="newFormInline.parent_id"
                class="w-full"
                :options="newFormInline.higherMenuOptions || []"
                :props="{
                  value: 'id',
                  label: 'menu_name',
                  emitPath: false,
                  checkStrictly: true
                }"
                clearable
                filterable
                :placeholder="transformI18n('menu.pureParentMenuPlaceholder')"
                @clear="newFormInline.parent_id = 0"
              >
                <template #default="{ node, data }">
                  <span>{{ data.menu_name }}</span>
                  <span v-if="!node.isLeaf">({{ data.children.length }})</span>
                </template>
              </el-cascader>
            </el-form-item>
          </re-col>

          <!-- 菜单标识码 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item
              :label="transformI18n('menu.pureMenuCode')"
              prop="menu_code"
            >
              <el-input
                v-model="newFormInline.menu_code"
                :disabled="!!newFormInline.id"
                clearable
                :placeholder="transformI18n('menu.pureMenuCodePlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 菜单名称 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item
              :label="transformI18n('menu.pureMenuName')"
              prop="menu_name"
            >
              <el-input
                v-model="newFormInline.menu_name"
                clearable
                :placeholder="transformI18n('menu.pureMenuNamePlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 页面标题 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureTitle')" prop="title">
              <el-input
                v-model="newFormInline.title"
                clearable
                :placeholder="transformI18n('menu.pureTitlePlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 路由路径 -->
          <re-col v-if="showPath" :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.purePath')" prop="path">
              <el-input
                v-model="newFormInline.path"
                clearable
                :placeholder="transformI18n('menu.purePathPlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 组件路径 -->
          <re-col v-if="showComponent" :value="12" :xs="24" :sm="24">
            <el-form-item
              :label="transformI18n('menu.pureComponent')"
              prop="component"
            >
              <el-input
                v-model="newFormInline.component"
                clearable
                :placeholder="transformI18n('menu.pureComponentPlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 外部链接 -->
          <re-col v-if="showExternalUrl" :value="12" :xs="24" :sm="24">
            <el-form-item
              :label="transformI18n('menu.pureExternalUrl')"
              prop="external_url"
            >
              <el-input
                v-model="newFormInline.external_url"
                clearable
                :placeholder="transformI18n('menu.pureExternalUrlPlaceholder')"
              />
            </el-form-item>
          </re-col>

          <!-- 菜单图标 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureIcon')">
              <IconSelect v-model="newFormInline.icon" class="w-full" />
            </el-form-item>
          </re-col>

          <!-- 排序 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureSortOrder')">
              <el-input-number
                v-model="newFormInline.sort_order"
                class="w-full!"
                :min="0"
                :max="9999"
                controls-position="right"
              />
            </el-form-item>
          </re-col>

          <!-- 状态 -->
          <re-col :value="10" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureStatus')">
              <Segmented
                v-model="statusIndex"
                :options="
                  statusOptions.map(item => ({
                    ...item,
                    label: transformI18n(item.label)
                  }))
                "
              />
            </el-form-item>
          </re-col>

          <!-- 是否隐藏 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureIsHidden')">
              <Segmented
                v-model="hiddenIndex"
                :options="
                  hiddenOptions.map(item => ({
                    ...item,
                    label: transformI18n(item.label)
                  }))
                "
              />
            </el-form-item>
          </re-col>

          <!-- 是否公共 -->
          <re-col :value="12" :xs="24" :sm="24">
            <el-form-item :label="transformI18n('menu.pureIsPublic')">
              <Segmented
                v-model="publicIndex"
                :options="
                  publicOptions.map(item => ({
                    ...item,
                    label: transformI18n(item.label)
                  }))
                "
              />
            </el-form-item>
          </re-col>
        </el-row>
      </re-col>

      <!-- 权限分配区域 -->
      <re-col
        v-if="showPermissionAssign"
        :value="11"
        :xs="24"
        class="permission-panel"
      >
        <el-divider content-position="left">
          {{ transformI18n("menu.purePermissionAssign") }}
        </el-divider>
        <el-form-item label="" class="permission-tree-item">
          <div class="permission-tree-wrapper">
            <el-input
              v-model="permissionSearchValue"
              :placeholder="transformI18n('menu.purePermissionSearch')"
              clearable
              class="mb-2"
            >
              <template #prefix>
                <el-icon><IconifyIconOffline icon="ri:search-line" /></el-icon>
              </template>
            </el-input>
            <el-scrollbar max-height="420px">
              <el-tree
                ref="permissionTreeRef"
                v-loading="permissionLoading"
                :data="permissionTreeData"
                :props="{
                  children: 'children',
                  label: 'perm_name'
                }"
                show-checkbox
                node-key="id"
                :default-checked-keys="newFormInline.permission_ids || []"
                :filter-node-method="filterPermissionNode"
                @check="handlePermissionCheck"
              >
                <template #default="{ data }">
                  <span class="custom-tree-node">
                    <span>{{ data.perm_name }}</span>
                    <span class="text-gray-400 text-xs ml-2">{{
                      data.perm_code
                    }}</span>
                  </span>
                </template>
              </el-tree>
            </el-scrollbar>
          </div>
        </el-form-item>
      </re-col>
    </el-row>
  </el-form>
</template>

<style scoped>
.menu-form-layout {
  width: 100%;
}

.permission-panel {
  border-left: 1px solid var(--el-border-color-light);
  padding-left: 15px;
}

.permission-tree-item :deep(.el-form-item__content) {
  display: block !important;
}

.permission-tree-wrapper {
  width: 100%;
}

.custom-tree-node {
  display: flex;
  align-items: center;
  flex: 1;
}
</style>
