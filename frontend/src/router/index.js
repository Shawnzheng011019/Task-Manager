import { createRouter, createWebHistory } from "vue-router";
import TaskList from "../views/TaskList.vue";
import TaskCreate from "../views/TaskCreate.vue";
import TaskDetail from "../views/TaskDetail.vue";
import Settings from "../views/Settings.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      redirect: "/tasks",
    },
    {
      path: "/tasks",
      name: "tasks",
      component: TaskList,
    },
    {
      path: "/tasks/create",
      name: "task-create",
      component: TaskCreate,
    },
    {
      path: "/tasks/:id",
      name: "task-detail",
      component: TaskDetail,
      props: true,
    },
    {
      path: "/settings",
      name: "settings",
      component: Settings,
    },
  ],
});

export default router;
