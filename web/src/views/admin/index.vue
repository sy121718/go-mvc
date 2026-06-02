<script setup lang="ts">
defineOptions({
    name: "Admin"
})
import { useAdmin } from './utils/hook';

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
    handleSortChange,
    getAdminStatusTagType,
    getAdminStatusLabel,
    openAdd,
    openBatchDelete
} = useAdmin();
</script>

<template>

    <div class="main">
        <el-form ref="formRef" :model="form" label-position="top" class="filter-form" @keyup.enter="onSearch">
            <el-row :gutter="20">
                <el-col :span="6">
                    <el-form-item label="邮箱" prop="email">
                        <el-input v-model="form.email" placeholder="邮箱搜索" clearable maxlength="100"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="6">
                    <el-form-item label="姓名" prop="name">
                        <el-input v-model="form.name" placeholder="姓名搜索" clearable maxlength="50"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="6">
                    <el-form-item label="状态筛选" prop="status">
                        <el-select v-model="form.status" placeholder="状态筛选" clearable>
                            <el-option label="启用" :value="1" />
                            <el-option label="禁用" :value="2" />
                            <el-option label="密码错误封禁" :value="3" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="6" class="filter-actions-col">
                    <div class="filter-actions">
                        <el-button type="primary" @click="onSearch">搜索</el-button>
                        <el-button @click="resetForm(formRef)">重置</el-button>
                    </div>
                </el-col>
            </el-row>
        </el-form>

        <div class="toolbar">
            <el-button type="primary" @click="openAdd">添加管理员</el-button>
            <el-button type="danger" @click="openBatchDelete">批量删除</el-button>
        </div>

        <el-table :data="dataList" v-loading="loading" row-key="id" stripe
            @sort-change="handleSortChange" style="width: 100%">

            <el-table-column type="selection" width="55" fixed="left" />
            <el-table-column v-for="col in columns" 
                :key="(col.prop as string) || col.label"
                :prop="col.prop as string" 
                :label="col.label" 
                :width="col.width"
                :min-width="col.minWidth" 
                :fixed="col.fixed" 
                :sortable="col.sortable"
                :align="col.align" 
                :formatter="col.formatter">

                <template #default="{ row }" v-if="col.slot === 'status'">
                    <el-tag :type="getAdminStatusTagType(row.status)" effect="plain">
                        {{ getAdminStatusLabel(row.status) }}
                    </el-tag>
                </template>
                <template #default v-else-if="col.slot === 'operation'">
                    <el-button type="primary" link size="small">编辑</el-button>
                    <el-button type="danger" link size="small">删除</el-button>
                </template>

            </el-table-column>

        </el-table>

        <el-pagination
            v-if="total > 0"
            :page-size="form.limit"
            :current-page="form.page"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            background
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            style="margin-top: 16px; justify-content: flex-end;"
        />

    </div>

</template>

<style scoped>
.main {
    margin: 12px;
    background: var(--el-bg-color);
    padding: 16px;
    border-radius: 8px;
}

.filter-form :deep(.el-form-item) {
    margin-bottom: 16px;
}

.filter-form :deep(.el-form-item__label) {
    padding-bottom: 8px;
    line-height: 20px;
}

.filter-actions-col {
    display: flex;
    align-items: flex-end;
}

.filter-actions {
    width: 100%;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-bottom: 16px;
}

.toolbar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
}
</style>