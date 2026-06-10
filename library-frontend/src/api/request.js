import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建Axios实例
const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

// 请求拦截器：自动附加JWT Token
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器：统一处理错误
request.interceptors.response.use(
  response => {
    const res = response.data
    // GoFrame返回 code=0 表示成功
    if (res.code !== 0 && res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message))
    }
    return res
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        localStorage.removeItem('token')
        window.location.hash = '#/login'
        ElMessage.error('登录已过期，请重新登录')
      } else {
        ElMessage.error(data?.message || `请求错误 (${status})`)
      }
    } else {
      ElMessage.error('网络错误，请检查后端是否启动')
    }
    return Promise.reject(error)
  }
)

export default request
