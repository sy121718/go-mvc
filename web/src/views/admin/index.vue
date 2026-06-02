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
    getAdminStatusLabel
} = useAdmin();
</script>

<template>

    <div class="main">
        <el-form ref="formRef" :model="form" label-width="80px" @keyup.enter="onSearch">
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
                    <el-form-item label="手机号" prop="phone">
                        <el-input v-model="form.phone" placeholder="手机号搜索" clearable maxlength="100"></el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="6">
                    <el-form-item label="状态筛选" prop="status">
                        <el-select v-model="form.status" placeholder="状态筛选">
                            <el-option label="启用" :value="1" />
                            <el-option label="禁用" :value="2" />
                            <el-option label="密码错误封禁" :value="3" />
                        </el-select>
                    </el-form-item>
                </el-col>

                <el-col :span="24" style="text-align: right;margin-bottom: 12px;">
                    <el-button type="primary" @click="onSearch">搜索</el-button>
                    <el-button @click="resetForm(formRef)">重置</el-button>
                </el-col>
            </el-row>
        </el-form>

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
                <template #default="{ row }" v-else-if="col.slot === 'operation'">
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
</style>