import { $t } from "@/plugins/i18n";

export default {
  path: "/admin",
  redirect: "/admin/list",
  meta: {
    icon: "ri/admin-line",
    title: $t("menus.pureAdmin"),
    rank: 1
  },
  children: [
    {
      path: "/admin/list",
      name: "AdminList",
      component: () => import("@/views/admin/index.vue"),
      meta: {
        title: $t("menus.pureAdmin")
      }
    }
  ]
} satisfies RouteConfigsTable;