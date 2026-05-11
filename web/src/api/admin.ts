//管理员模块接口
import { http } from "@/utils/http";

//管理员列表
export type AdminListResult = {
  code: number;
  message?: string;
  data: {
    total: number;
    list: any[];
  };
};

//新增管理员请求参数
export type AdminCreateReq = {
  avatar?: string;
  email: string;
  username: string;
  phone?: string;
  password: string;
};

//新增管理员响应
export type AdminCreateResp = {
  code: number;
  message?: string;
  data: {
    id: number;
    username: string;
  };
};

//导出管理员列表
export const getAdminList = () => {
  return http.request<AdminListResult>("get", "/admin/list");
};

//新增管理员
export const createAdmin = (data: AdminCreateReq) => {
  return http.request<AdminCreateResp>("post", "/admin/create", { data });
};