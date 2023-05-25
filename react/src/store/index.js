import React from "react"
import LoginStore from "./login.Store"
import RegisterStore from "./register.Store"
import UserStore from "./user.Store"
import ChannelStore from "./chanel.Store"


class RootStore {
  constructor() {
    this.loginStore = new LoginStore()
    this.userStore = new UserStore()
    this.resStore = new RegisterStore()
    this.channelStore = new ChannelStore()
    //....
  }
}

//实例化根
//导出useStore

const rootStore = new RootStore()
const context = React.createContext(rootStore)

const useStore = () => React.useContext(context)

export {useStore}