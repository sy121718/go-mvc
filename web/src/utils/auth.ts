import { useUserStoreHook } from "@/store/modules/user";
import { isString, isIncludeAllChildren } from "@pureadmin/utils";

/** 用户信息存储 key */
export const userKey = "user-info";

/** 多标签页存储 key */
export const multipleTabsKey = "multiple-tabs";

/** localStorage 中存储的数据结构 */
export type DataInfo<T> = {
  expires?: T;
};

/** token 内存变量，不持久化到 cookie/localStorage */
let accessToken: string | null = null;

/** 获取`token` */
export function getToken(): string | null {
  return accessToken;
}

/** 设置`token`（由 Axios 响应拦截器调用） */
export function setToken(token: string): void {
  accessToken = token;
}

/** 删除`token` */
export function removeToken(): void {
  accessToken = null;
}

/** 格式化token（jwt格式） */
export const formatToken = (token: string): string => {
  return "Bearer " + token;
};

/** 是否有按钮级别的权限 */
export const hasPerms = (value: string | Array<string>): boolean => {
  if (!value) return false;
  const allPerms = "*:*:*";
  const { permissions } = useUserStoreHook();
  if (!permissions) return false;
  if (permissions.length === 1 && permissions[0] === allPerms) return true;
  const isAuths = isString(value)
    ? permissions.includes(value)
    : isIncludeAllChildren(value, permissions);
  return isAuths ? true : false;
};