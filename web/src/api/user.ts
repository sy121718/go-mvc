//登录和获取信息
import { http } from "@/utils/http";

// 获取验证码
export type CaptchaResult = {
  code: number;
  message?: string;
  data: {
    captcha_id: string;
    captcha: string;
  };
};

/** 登录 */
export const getLogin = (data?: object) => {
  return http.request<any>("post", "/api/admin/login", { data });
};

/** 获取验证码 */
export const getCaptcha = () => {
  return http.request<CaptchaResult>("get", "/api/captcha");
};

/** 获取当前用户信息 */
export const getProfile = () => {
  return http.request<any>("get", "/api/admin/profile");
};

