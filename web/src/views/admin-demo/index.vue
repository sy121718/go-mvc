<script setup lang="ts">
defineOptions({
  name: "Admin"
});

import { useAdmin } from "./utils/hook";

const {
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
} = useAdmin();
</script>

<template>
  <div class="main">
    <!-- 搜索表单 -->
    <el-form
      ref="formRef"
      :model="form"
      label-width="80px"
      @keyup.enter="onSearch"
    >
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="form.email"
              placeholder="邮箱搜索"
              clearable
              maxlength="100"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="姓名" prop="name">
            <el-input
              v-model="form.name"
              placeholder="姓名搜索"
              clearable
              maxlength="50"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="手机号" prop="phone">
            <el-input
              v-model="form.phone"
              placeholder="手机号搜索"
              clearable
              maxlength="20"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="状态" prop="status">
            <el-select
              v-model="form.status"
              placeholder="状态筛选"
              clearable
            >
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="2" />
              <el-option label="密码错误封禁" :value="3" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="24" style="text-align: right; margin-bottom: 12px">
          <el-button type="primary" @click="onSearch">
            搜索
          </el-button>
          <el-button @click="resetForm(formRef)">
            重置
          </el-button>
        </el-col>
      </el-row>
    </el-form>

    <!-- 表格 -->
    <RePureTableBar title="管理员列表" :columns="columns">
      <template #default="{ size, dynamicColumns }">
        <el-table
          ref="tableRef"
          :data="dataList"
          :columns="dynamicColumns"
          :size="size"
          :loading="loading"
          row-key="id"
          stripe
          @sort-change="handleSortChange"
        >
          <template v-for="col in dynamicColumns" :key="col.prop">
            <el-table-column
              v-if="col.slot === 'operation'"
              :label="col.label"
              :fixed="col.fixed"
              :width="col.width"
              v-bind="col"
            >
              <template #default="{ row }">
                <el-button type="primary" link size="small">
                  编辑
                </el-button>
                <el-button type="danger" link size="small">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </template>
        </el-table>
      </template>
    </RePureTableBar>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="form.page"
      v-model:page-size="form.limit"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      background
      style="margin-top: 16px; justify-content: flex-end"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<style scoped>
.main {
  padding: 16px;
  background: var(--el-bg-color);
  border-radius: 8px;
}
</style>