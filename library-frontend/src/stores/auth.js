import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi, getProfileApi } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)

  // 设置Token（同时存localStorage，刷新不丢失）
  function setToken(newToken) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  // 登录：调用API → 存Token → 获取用户信息
  async function login(loginData) {
    const res = await loginApi(loginData)
    const data = res.data
    setToken(data.token)
    user.value = { id: data.id, name: data.name, email: data.email, phone: data.phone }
    return data
  }

  // 获取个人信息
  async function fetchProfile() {
    const res = await getProfileApi()
    user.value = res.data
    return res.data
  }

  // 登出
  function logout() {
    setToken('')
    user.value = null
  }

  return { token, user, isLoggedIn, setToken, login, fetchProfile, logout }
})
