import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/LoginPage.vue'),
    meta: { requiresAuth: false, layout: 'auth' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/pages/RegisterPage.vue'),
    meta: { requiresAuth: false, layout: 'auth' },
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/boards',
    name: 'Boards',
    component: () => import('@/pages/BoardsPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/boards/:id',
    name: 'Board',
    component: () => import('@/pages/BoardPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/teams',
    name: 'Teams',
    component: () => import('@/pages/TeamPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/reports',
    name: 'Reports',
    component: () => import('@/pages/ReportsPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/pages/ProfilePage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFoundPage.vue'),
    meta: { requiresAuth: false, layout: 'auth' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.requiresAuth && !token) {
    return { name: 'Login' }
  }
  if (!to.meta.requiresAuth && token && (to.name === 'Login' || to.name === 'Register')) {
    return { name: 'Dashboard' }
  }
})

export default router
