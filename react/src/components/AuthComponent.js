//判断token是否存在
//如果不存在重定向到登录路由
import { Navigate } from 'react-router-dom'

import { getToken } from '@/utils'


function AuthComponent ({ children }) {
  const isToken = getToken()
  if (isToken) {
    return <>{children}</>
  } else {
    return <Navigate to='/login' replace />
  }
}

export {
  AuthComponent
}