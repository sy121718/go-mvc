<script setup lang="ts">
import { reactive, ref } from "vue"
import { createAdminRules } from "../utils/rule"
import type { AdminCreateReq } from "@/api/admin"

const formRef = ref()
const form = reactive<AdminCreateReq>({
  username: "",
  email: "",
  phone: "",
  password: "",
  remark:""
})

function getRef() {
  return formRef.value
}

function getForm() {
  return form
}

defineExpose({ getRef, getForm })
</script>

<template>
  <el-form ref="formRef" :model="form" :rules="createAdminRules" label-width="80px">
    <el-form-item label="用户名" prop="username" required>
      <el-input v-model="form.username" placeholder="请输入用户名" maxlength="50" />
    </el-form-item>
    <el-form-item label="邮箱" prop="email" required>
      <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="100" />
    </el-form-item>
    <el-form-item label="手机号" prop="phone">
      <el-input v-model="form.phone" placeholder="请输入手机号" maxlength="20" />
    </el-form-item>
    <el-form-item label="密码" prop="password" required>
      <el-input v-model="form.password" type="password" placeholder="请输入密码" maxlength="50" show-password />
    </el-form-item>
    <el-form-item label="备注" prop="remark" required>
      <el-input v-model="form.remark" type="textarea" placeholder="管理员备注" maxlength="500" />
    </el-form-item>
  </el-form>
</template>