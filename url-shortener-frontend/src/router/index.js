import { createRouter, createWebHistory } from 'vue-router';
import FormComponent from '@/components/FormComponent.vue';
import AnalyticsComponent from '@/components/AnalyticsComponent.vue';
import DashboardComponent from '@/components/DashboardComponent.vue';

const routes = [
  {
    path: '/',
    name: 'home',
    component: FormComponent
  },
  {
    path: '/analytics/:shortUrl',
    name: 'analytics',
    component: AnalyticsComponent // Component for analytics route
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardComponent
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
