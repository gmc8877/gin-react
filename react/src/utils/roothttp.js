//封装axios
import axios from 'axios'
import { getRootToken } from './token'
import { history } from './history'

const roothttp = axios.create({
  baseURL: '/api',
  timeout: 5000,
})
// 添加请求拦截器
roothttp.interceptors.request.use((config) => {
  const token = getRootToken()
  if (token) {
    config.headers.Token = token
  }
  return config
}, (error) => {
  return Promise.reject(error)
})

// 添加响应拦截器
roothttp.interceptors.response.use((response) => {
  // 2xx 范围内的状态码都会触发该函数。
  // 对响应数据做点什么
  return response.data
}, (error) => {
  // 超出 2xx 范围的状态码都会触发该函数。
  // 对响应错误做点什么

  if (error.response.status === 401) {
    //跳回桌面
    history.push('/root_login')
  }

  return Promise.reject(error)
})

export { roothttp }