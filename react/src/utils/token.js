const key = 'pc-key'
const rootkey = 'root-key'

const setToken = (token) => {
  return window.localStorage.setItem(key, token)
}

const setRootToken = (token) => {
  return window.localStorage.setItem(rootkey, token)
}

const getToken = () => {
  return window.localStorage.getItem(key)
}

const getRootToken = () => {
  return window.localStorage.getItem(rootkey)
}

const removeToken = () => {
  return window.localStorage.removeItem(key)
}

const removeRootToken = () => {
  return window.localStorage.removeItem(rootkey)
}

export {
  setToken,
  getToken,
  removeToken,
  setRootToken,
  getRootToken,
  removeRootToken
}