import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('../layouts/AppLayout.vue'),
    redirect: '/books',
    children: [
      {
        path: 'books',
        name: 'Books',
        component: () => import('../views/Books.vue'),
        meta: { title: '图书浏览' }
      },
      {
        path: 'books/manage',
        name: 'BookManage',
        component: () => import('../views/BookManage.vue'),
        meta: { title: '图书管理', requiresAuth: true }
      },
      {
        path: 'borrows',
        name: 'Borrows',
        component: () => import('../views/Borrows.vue'),
        meta: { title: '我的借阅', requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/Profile.vue'),
        meta: { title: '个人中心', requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫：需要登录的页面自动跳转
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 图书管理系统` : '图书管理系统'

  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.isLoggedIn) {
      ElMessage.warning('请先登录')
      next('/login')
      return
    }

    // 图书管理页面需要管理员权限
    if (to.path === '/books/manage' && auth.user?.role !== 'admin') {
      ElMessage.error('权限不足，需要管理员身份')
      next('/books')
      return
    }
  }

  next()
})

export default router
