import { makeAutoObservable } from 'mobx'
import { http } from '@/utils'

class RegisterStore {
  errmsg = 'OK'
  constructor() {
    //响应式
    makeAutoObservable(this)
  }
  getCaptcha = async ({ mobile, email }) => {
    try {
      await http.post('/register/captcha', { mobile, email })
      this.errmsg = 'OK'
    } catch (err) {
      this.errmsg = err.response.data.message
    }
  }
  getRegister = async ({ mobile, password, email, captcha }) => {
    try {
      await http.post('/register', { mobile, password, email, captcha })
      this.errmsg = 'OK'
    } catch (err) {
      this.errmsg = err.response.data.message
    }
  }
}
export default RegisterStore