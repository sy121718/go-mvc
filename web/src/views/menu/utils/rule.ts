import { reactive } from "vue";
import type { FormRules } from "element-plus";

/** 自定义表单规则校验 */
export const formRules = reactive(<FormRules>{
  menu_code: [
    { required: true, message: "菜单标识码为必填项", trigger: "blur" },
    {
      pattern: /^[a-zA-Z_]+$/,
      message: "菜单标识码只能包含英文字母和下划线",
      trigger: "blur"
    }
  ],
  menu_name: [{ required: true, message: "菜单名称为必填项", trigger: "blur" }],
  title: [{ required: true, message: "页面标题为必填项", trigger: "blur" }],
  path: [
    { required: true, message: "路由路径为必填项", trigger: "blur" },
    {
      pattern: /^\//,
      message: "路由路径必须以/开头",
      trigger: "blur"
    }
  ],
  component: [{ required: true, message: "组件路径为必填项", trigger: "blur" }],
  external_url: [
    {
      pattern: /^https?:\/\/.+/,
      message: "请输入完整的URL地址",
      trigger: "blur"
    }
  ]
});
