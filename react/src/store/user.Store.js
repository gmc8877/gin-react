import { makeAutoObservable } from 'mobx'
import { http, roothttp } from '@/utils'


class UserStore {
  userInfo = ''
  rootInfo = ''
  constructor() {
    makeAutoObservable(this)
  }
  getUserInfo = async () => {
    //调用接口获得的数据
    const res = await http.get('/userinfo')
    this.userInfo = res.name
  }
  getRootInfo = async () => {
    //调用接口获得的数据
    const res = await roothttp.get('/root/rootinfo')

    this.rootInfo = res.name
  }
}

export default UserStore