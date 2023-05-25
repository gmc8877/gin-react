//判断token是否存在
//如果不存在重定向到登录路由
import { Navigate } from 'react-router-dom'

import { getRootToken } from '@/utils'


function RootComponent ({ children }) {
  const isToken = getRootToken()
  if (isToken) {
    return <>{children}</>
  } else {
    return <Navigate to='/root_login' replace />
  }
}

export {
  RootComponent
}