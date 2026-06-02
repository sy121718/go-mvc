//管理员模块接口
import { http } from "@/utils/http";

// ──────── 列表 ────────

//列表请求参数
export type AdminListReq = {
  page?: number;
  limit?: number;
  email?: string;
  name?: string;
  status?: number;
  sort_field?: string;
  sort_order?: string;
};

//列表响应
export type AdminListResp = {
  code: number;
  message?: string;
  data: {
    total: number;
    list: any[];
  };
};

// ──────── 新增 ────────

//新增请求参数
export type AdminCreateReq = {
  avatar?: string;
  email: string;
  username: string;
  phone?: string;
  password: string;
};

//新增响应
export type AdminCreateResp = {
  code: number;
  message?: string;
  data: {
    id: number;
    username: string;
  };
};

// ──────── 详情 ────────

//详情响应
export type AdminDetailResp = {
  code: number;
  message?: string;
  data: {
    id: number;
    username: string;
    nickname: string;
    avatar: string;
    email: string;
    phone: string;
    status: number;
    is_admin: number;
    roles: any[];
    menus: any[];
    register_ip: string;
    register_location: string;
    last_login_ip: string;
    last_login_location: string;
    last_login_time: string;
    create_by: number;
    create_time: string;
  };
};

// ──────── API 函数 ────────

export const getAdminList = (params?: AdminListReq) => {
  return http.request<AdminListResp>("get", "/api/admin/list", { params });
};

export const createAdmin = (data: AdminCreateReq) => {
  return http.request<AdminCreateResp>("post", "/api/admin/create", { data });
};

export const getAdminDetail = (id: number) => {
  return http.request<AdminDetailResp>("get", "/api/admin/detail", {
    params: { id }
  });
};