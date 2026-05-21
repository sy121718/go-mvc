interface FormItemProps {
  /** 菜单ID */
  id?: number;
  /** 菜单唯一标识码（必填，英文+下划线，全局唯一） */
  menu_code: string;
  /** 菜单名称（必填，中文名称，显示在侧边栏） */
  menu_name: string;
  /** 父菜单ID（必填，0=顶级菜单，>0=子菜单的父ID） */
  parent_id: number;
  /** 菜单类型（必填，1=目录、2=菜单、3=外链、4=iframe、5=按钮） */
  type: number;
  /** 菜单备注（选填，用于说明菜单功能，可在鼠标悬浮时显示） */
  title: string;
  /** 路由路径（必填，前端路由地址，以/开头） */
  path: string;
  /** 组件路径（必填，前端组件文件路径，不需要.vue后缀） */
  component: string;
  /** 外部链接（选填，type=3外链或type=4的iframe时必填，完整URL） */
  external_url: string;
  /** 图标（选填，图标类名） */
  icon: string;
  /** 状态（必填，0=禁用 1=启用，默认1） */
  status: number;
  /** 是否隐藏（选填，0=显示 1=隐藏，默认0） */
  is_hidden: number;
  /** 是否公共（选填，0=否 1=是，公共菜单无角色用户也能访问，默认0） */
  is_public: number;
  /** 排序（选填，数字越小越靠前，默认0） */
  sort_order: number;
  /** 关联权限ID数组（选填，绑定的权限字典ID，用于权限控制） */
  permission_ids: number[];

  // 用于表单组件的辅助字段
  higherMenuOptions?: Record<string, unknown>[];
}

interface FormProps {
  formInline: FormItemProps;
}

export type { FormItemProps, FormProps };
