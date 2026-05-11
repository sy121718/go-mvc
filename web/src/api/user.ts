//登录和获取信息
import { http } from "@/utils/http";

export type UserResult = {
  code: number;
  message: string;
  data: {
    /** `token` */
    accessToken: string;
    /** 用于兼容本地存储结构的刷新标识 */
    refreshToken: string;
    /** `accessToken`的过期时间（格式'xxxx/xx/xx xx:xx:xx'） */
    expires: string;
    /** 用户名 */
    username: string;
    /** 昵称 */
    nickname: string;
    /** 邮箱 */
    email: string;
    /** 头像 */
    avatar: string;
  };
};
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
  return http.request<UserResult>("post", "/admin/login", { data });
};

/** 获取验证码 */
export const getCaptcha = () => {
  return http.request<CaptchaResult>("get", "/api/captcha");
};

