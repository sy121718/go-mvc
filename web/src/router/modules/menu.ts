import { $t } from "@/plugins/i18n";

export default {
  path: "/menu",
  redirect: "/menu/list",
  meta: {
    icon: "ri/menu-line",
    title: $t("menus.pureMenu"),
    rank: 2
  },
  children: [
    {
      path: "/menu/list",
      name: "MenuList",
      component: () => import("@/views/menu/index.vue"),
      meta: {
        title: $t("menus.pureMenu")
      }
    }
  ]
} satisfies RouteConfigsTable;