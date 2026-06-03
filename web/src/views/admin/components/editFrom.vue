<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import { editAdminRules } from "../utils/rule"
import type { AdminDetailResp } from "@/api/admin"

interface EditFormData {
  id: number | null;
  username: string;
  email: string;
  phone: string;
  remark: string;
}

const props = defineProps<{
  detail?: AdminDetailResp["data"];
}>()

const editFromRef = ref()
const form = reactive<EditFormData>({
  id: null,
  username: "",
  email: "",
  phone: "",
  remark: ""
})

watch(() => props.detail, (val) => {
  if (!val) return
  form.id = val.id
  form.username = val.username
  form.email = val.email
  form.phone = val.phone || ""
  form.remark = val.remark || ""
}, { immediate: true })

function getRef() {
  return editFromRef.value
}

function getForm() {
  return form
}

defineExpose({ getRef, getForm })
</script>



<template>

  <el-form ref="editFromRef" :model="form" :rules="editAdminRules" label-width="80px">

    <el-row :gutter="16">
      <el-col :span="12">
        <el-form-item label="邮箱" type="email" prop="email" required>
          <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="100" />
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item label="用户名" prop="username" required>
          <el-input v-model="form.username" placeholder="请输入用户名" maxlength="100" />
        </el-form-item>
      </el-col>

      <el-col :span="12">
        <el-form-item label="手机号码" prop="phone" >
          <el-input v-model="form.phone" placeholder="请输入手机号码" maxlength="11" />
        </el-form-item>
      </el-col>


      <el-col :span="24">
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" maxlength="200" />
        </el-form-item>
      </el-col>
    </el-row>


  </el-form>


</template>