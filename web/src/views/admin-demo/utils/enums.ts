export const adminStatusOptions = [
  { value: 1, label: "启用", tagType: "success" },
  { value: 2, label: "禁用", tagType: "danger" },
  { value: 3, label: "密码错误封禁", tagType: "warning" }
];

export const adminIsAdminOptions = [
  { value: 0, label: "否" },
  { value: 1, label: "是" }
];

export const getAdminStatusTagType = (status: number) => {
  const opt = adminStatusOptions.find(o => o.value === status);
  return opt?.tagType ?? "info";
};

export const getAdminStatusLabel = (status: number) => {
  const opt = adminStatusOptions.find(o => o.value === status);
  return opt?.label ?? "未知";
};