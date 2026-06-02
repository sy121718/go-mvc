<script setup lang="ts">
defineOptions({
    name: "Admin"
})
import { useAdmin } from './utils/hook';

//提前导出解构全局使用
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
        <!-- 把这个 el-form 组件实例赋给你在 <script setup> 里定义的 const formRef = ref()，之后就能在 JS 里调它的方法。 -->
        <el-form ref="formRef" :model="form" label-width="80px" @keyup.enter="onSearch">
            <el-row :gutter="20">
                <el-col :span="6">
                    <el-form-item label="邮箱" prop="email">
                        <!-- clearable 清空输入框 -->
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
                <el-table ref="tableRef" :data="dataList" :columns="dynamicColumns" :size="size" row-key="id" stripe
                    @sort-change="handleSortChange">

                    <template v-for="col in dynamicColumns" :key="col.prop">
                        <el-table-column v-if="col.slot === 'operation'" :label="col.label" :fixed="col.fixed"
                            :width="col.width" v-bind="col">
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


        <el-pagination />

    </div>



</template>