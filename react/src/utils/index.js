//把所有的工具函数导出的模块在这里导入
import { http } from "./http"
import { setToken, getToken, removeToken, setRootToken, getRootToken, removeRootToken } from "./token"
import { roothttp } from "./roothttp"

export {
  http,
  roothttp,
  setToken,
  getToken,
  removeToken,
  setRootToken,
  getRootToken,
  removeRootToken
}