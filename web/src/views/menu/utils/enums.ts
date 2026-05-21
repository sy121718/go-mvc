import type { OptionsType } from "@/components/ReSegmented";

// 菜单类型选项：1=目录、2=菜单、3=外链、4=iframe、5=按钮
const menuTypeOptions: Array<OptionsType> = [
  {
    label: "menu.pureMenuTypeDirectory",
    value: 1
  },
  {
    label: "menu.pureMenuTypeMenu",
    value: 2
  },
  {
    label: "menu.pureMenuTypeLink",
    value: 3
  },
  {
    label: "menu.pureMenuTypeIframe",
    value: 4
  },
  {
    label: "menu.pureMenuTypeButton",
    value: 5
  }
];

// 状态选项：0=禁用 1=启用
const statusOptions: Array<OptionsType> = [
  {
    label: "menu.pureStatusEnable",
    value: 1
  },
  {
    label: "menu.pureStatusDisable",
    value: 0
  }
];

// 是否隐藏选项：0=显示 1=隐藏
const hiddenOptions: Array<OptionsType> = [
  {
    label: "menu.pureIsHiddenShow",
    value: 0
  },
  {
    label: "menu.pureIsHiddenHide",
    value: 1
  }
];

// 是否公共选项：0=否 1=是（公共菜单无角色用户也能访问）
const publicOptions: Array<OptionsType> = [
  {
    label: "menu.pureIsPublicNo",
    value: 0
  },
  {
    label: "menu.pureIsPublicYes",
    value: 1
  }
];

export { menuTypeOptions, statusOptions, hiddenOptions, publicOptions };
