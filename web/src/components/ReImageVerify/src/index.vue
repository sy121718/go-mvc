<script setup lang="ts">
import { watch } from "vue";
import { useImageVerify } from "./hooks";

defineOptions({
  name: "ReImageVerify"
});

interface Props {
  code?: string;
}

interface Emits {
  (e: "update:code", code: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  code: ""
});

const emit = defineEmits<Emits>();

const { domRef, imgCode, captchaKey, loading, setImgCode, getImgCode } =
  useImageVerify();

watch(
  () => props.code,
  newValue => {
    setImgCode(newValue);
  }
);
watch(imgCode, newValue => {
  emit("update:code", newValue);
});

defineExpose({ getImgCode, captchaKey });
</script>

<template>
  <div class="captcha-container">
    <img
      ref="domRef"
      width="120"
      height="40"
      :class="[
        'captcha-image',
        {
          'cursor-pointer': !loading,
          'cursor-not-allowed': loading,
          'opacity-50': loading
        }
      ]"
      alt="验证码"
      @click="getImgCode"
    />
    <div v-if="loading" class="loading-overlay">
      <span class="loading-text">加载中...</span>
    </div>
  </div>
</template>

<style scoped>
.captcha-container {
  position: relative;
  display: inline-block;
  margin-left: 8px;
}

.captcha-image {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #ffffff;
  display: block;
  object-fit: contain;
  transition: all 0.3s ease;
}

.captcha-image:hover:not(.cursor-not-allowed) {
  border-color: #409eff;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  color: white;
  font-size: 12px;
}

.loading-text {
  font-weight: 500;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}
</style>
