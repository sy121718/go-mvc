<script setup lang="ts">
import { useI18n } from "vue-i18n";
import Motion from "./utils/motion";
import TypeIt from "@/components/ReTypeit";
import { ReImageVerify } from "@/components/ReImageVerify";
import { useRouter } from "vue-router";
import { message } from "@/utils/message";
import { loginRules } from "./utils/rule";
import { ref, reactive, toRaw, watch } from "vue";
import { debounce } from "@pureadmin/utils";
import { useNav } from "@/layout/hooks/useNav";
import { useEventListener } from "@vueuse/core";
import type { FormInstance } from "element-plus";
import { useLayout } from "@/layout/hooks/useLayout";
import { useUserStoreHook } from "@/store/modules/user";
import { initRouter, getTopMenu } from "@/router/utils";
import { bg, avatar, illustration } from "./utils/static";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { useTranslationLang } from "@/layout/hooks/useTranslationLang";
import { useDataThemeChange } from "@/layout/hooks/useDataThemeChange";

import dayIcon from "@/assets/svg/day.svg?component";
import darkIcon from "@/assets/svg/dark.svg?component";
import globalization from "@/assets/svg/globalization.svg?component";
import Lock from "~icons/ri/lock-fill";
import Check from "~icons/ep/check";
import User from "~icons/ri/user-3-fill";
import Info from "~icons/ri/information-line";
import Keyhole from "~icons/ri/shield-keyhole-line";

defineOptions({
  name: "Login"
});

const rememberDays = 7;
const router = useRouter();
const loading = ref(false);
const checked = ref(false);
const disabled = ref(false);
const ruleFormRef = ref<FormInstance>();
const captchaRef = ref<{ captchaKey?: string; getImgCode?: () => void } | null>(
  null
);

const { initStorage } = useLayout();
initStorage();

const { t } = useI18n();
const { dataTheme, overallStyle, dataThemeChange } = useDataThemeChange();
dataThemeChange(overallStyle.value);
const { title, getDropdownItemStyle, getDropdownItemClass } = useNav();
const { locale, translationCh, translationEn } = useTranslationLang();

const ruleForm = reactive({
  username: "",
  password: "",
  verifyCode: ""
});

watch(checked, value => {
  useUserStoreHook().SET_ISREMEMBERED(value);
});
useUserStoreHook().SET_LOGINDAY(rememberDays);

const trimInput = (field: keyof typeof ruleForm) => {
  ruleForm[field] = ruleForm[field].trim();
};

const onLogin = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;

  trimInput("username");
  trimInput("password");
  trimInput("verifyCode");

  const valid = await formEl
    .validate()
    .then(() => true)
    .catch(() => false);

  if (!valid) return;

  loading.value = true;
  disabled.value = true;

  try {
    const res = await useUserStoreHook().loginByUsername({
      username: ruleForm.username,
      password: ruleForm.password,
      captcha_id: captchaRef.value?.captchaKey || "",
      captcha: ruleForm.verifyCode,
      remember_me: checked.value
    });

    if (res.code !== 200) {
      message(res.message || t("login.pureLoginFail"), { type: "error" });
      captchaRef.value?.getImgCode?.();
      return;
    }

    await initRouter();
    await router.push(getTopMenu(true).path);
    message(t("login.pureLoginSuccess"), { type: "success" });
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : t("login.pureLoginFail");
    message(errorMessage, { type: "error" });
    captchaRef.value?.getImgCode?.();
  } finally {
    loading.value = false;
    disabled.value = false;
  }
};

const immediateDebounce: any = debounce(
  formRef => onLogin(formRef),
  1000,
  true
);

useEventListener(document, "keydown", ({ code }) => {
  if (
    ["Enter", "NumpadEnter"].includes(code) &&
    !disabled.value &&
    !loading.value
  ) {
    immediateDebounce(ruleFormRef.value);
  }
});
</script>

<template>
  <div class="select-none">
    <img :src="bg" class="wave" />
    <div class="flex-c absolute right-5 top-3">
      <el-switch
        v-model="dataTheme"
        inline-prompt
        :active-icon="dayIcon"
        :inactive-icon="darkIcon"
        @change="dataThemeChange"
      />
      <el-dropdown trigger="click">
        <globalization
          class="hover:text-primary hover:bg-[transparent]! ml-1.5 h-[20px] w-[20px] cursor-pointer outline-hidden duration-300"
        />
        <template #dropdown>
          <el-dropdown-menu class="translation">
            <el-dropdown-item
              :style="getDropdownItemStyle(locale, 'zh')"
              :class="['dark:text-white!', getDropdownItemClass(locale, 'zh')]"
              @click="translationCh"
            >
              <IconifyIconOffline
                v-show="locale === 'zh'"
                class="check-zh"
                :icon="Check"
              />
              简体中文
            </el-dropdown-item>
            <el-dropdown-item
              :style="getDropdownItemStyle(locale, 'en')"
              :class="['dark:text-white!', getDropdownItemClass(locale, 'en')]"
              @click="translationEn"
            >
              <span v-show="locale === 'en'" class="check-en">
                <IconifyIconOffline :icon="Check" />
              </span>
              English
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <div class="login-container">
      <div class="img">
        <component :is="toRaw(illustration)" />
      </div>

      <div class="login-box">
        <div class="login-form">
          <avatar class="avatar" />
          <Motion>
            <h2 class="outline-hidden">
              <TypeIt :text="title" :speed="100" />
            </h2>
          </Motion>

          <el-form
            ref="ruleFormRef"
            :model="ruleForm"
            :rules="loginRules"
            size="large"
          >
            <Motion :delay="100">
              <el-form-item prop="username">
                <el-input
                  v-model="ruleForm.username"
                  clearable
                  :placeholder="t('login.pureUsername')"
                  :prefix-icon="useRenderIcon(User)"
                  @blur="trimInput('username')"
                />
              </el-form-item>
            </Motion>

            <Motion :delay="150">
              <el-form-item prop="password">
                <el-input
                  v-model="ruleForm.password"
                  clearable
                  show-password
                  :placeholder="t('login.purePassword')"
                  :prefix-icon="useRenderIcon(Lock)"
                  @blur="trimInput('password')"
                />
              </el-form-item>
            </Motion>

            <Motion :delay="200">
              <el-form-item prop="verifyCode">
                <el-input
                  v-model="ruleForm.verifyCode"
                  clearable
                  :placeholder="t('login.pureVerifyCode')"
                  :prefix-icon="useRenderIcon(Keyhole)"
                  @blur="trimInput('verifyCode')"
                >
                  <template #append>
                    <ReImageVerify ref="captchaRef" />
                  </template>
                </el-input>
              </el-form-item>
            </Motion>

            <Motion :delay="250">
              <el-form-item>
                <div class="flex h-[20px] w-full items-center justify-between">
                  <el-checkbox v-model="checked">
                    <span class="flex items-center">
                      {{ rememberDays }}{{ t("login.pureRemember") }}
                      <IconifyIconOffline
                        v-tippy="{
                          content: t('login.pureRememberInfo'),
                          placement: 'top'
                        }"
                        :icon="Info"
                        class="ml-1"
                      />
                    </span>
                  </el-checkbox>
                </div>
                <el-button
                  class="mt-4! w-full"
                  size="default"
                  type="primary"
                  :loading="loading"
                  :disabled="disabled"
                  @click="onLogin(ruleFormRef)"
                >
                  {{ t("login.pureLogin") }}
                </el-button>
              </el-form-item>
            </Motion>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import url("@/style/login.css");
</style>

<style lang="scss" scoped>
:deep(.el-input-group__append, .el-input-group__prepend) {
  padding: 0;
  background-color: transparent;
}

.translation {
  ::v-deep(.el-dropdown-menu__item) {
    padding: 5px 40px;
  }

  .check-zh {
    position: absolute;
    left: 20px;
  }

  .check-en {
    position: absolute;
    left: 20px;
  }
}
</style>
