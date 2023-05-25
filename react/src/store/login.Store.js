//login
import {makeAutoObservable} from 'mobx'
import { http, roothttp, setToken, getToken, removeToken, setRootToken, removeRootToken, getRootToken } from '@/utils'


class LoginStore {
  token = getToken() || ''
  rootToken = getRootToken() || ''
  message = ''
  constructor() {
    //响应式
    makeAutoObservable(this)
  }
  getToken = async ({ mobile, code }) => {
    try {
      //调用登录接口
      const res = await http.post('/login', {
        mobile, code
      })
      //存入token
      this.token = res.data.token
      this.message = 'OK'
     
      //持久化
      setToken(this.token)
    } catch (err){
      this.message = err.response.data.message
    }
  }
  getRootToken = async ({ mobile, code }) => {
    try {
      //调用登录接口
      const res = await roothttp.post('/root/login', {
        mobile, code
      })
      //存入token
      this.rootToken = res.data.token
      this.message = 'OK'
      //持久化
      setRootToken(this.rootToken)
    } catch (err) {
      this.message = err.response.data.message
    }
  }

  logout = () => {
    this.token = ''
    removeToken()
  }
  rootlogout = () => {
    this.rootToken = ''
    removeRootToken()
  }
}

export default LoginStore