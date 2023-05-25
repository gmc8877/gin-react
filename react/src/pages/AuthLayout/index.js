import { Layout, Menu, Popconfirm } from 'antd'
import {
  HomeOutlined,
  EditOutlined,
  DiffOutlined,
  LogoutOutlined
} from '@ant-design/icons'
import './index.scss'
import { Outlet, Link, useLocation, useNavigate} from 'react-router-dom'
import { useStore } from '@/store'
import { useEffect } from 'react'
import {observer} from 'mobx-react-lite'


const { Header, Sider } = Layout

const AuthLayout = () => {
  const { pathname } = useLocation()
  const { userStore, loginStore ,channelStore} = useStore()
  const navigate = useNavigate()
  useEffect(() => {
    userStore.getUserInfo()
    channelStore.loadChannelList()
    }, [userStore,channelStore]);
  //确定退出
  const onConfirm = () => {
    //删除token
    loginStore.logout()
    navigate('/login')
  }
  
  return (
    <Layout>
      <Header className="header">
        <div className="logo" />
        <div className="user-info">
          <span className="user-name">{ userStore.userInfo}</span>
          <span className="user-logout">
            <Popconfirm
              onConfirm={onConfirm}
              title="是否确认退出？"
              okText="退出"
              cancelText="取消">
              <LogoutOutlined /> 退出
            </Popconfirm>
          </span>
        </div>
      </Header>
      <Layout>
        <Sider width={150} className="sider">
          {/* 高亮原理 defaultSelectedKeys===item key*/}
          {/* 获取当前激活的path路径 */}
          <Menu
            mode="inline"
            selectedKeys={pathname}
            className="sidemenu"
            style={{ height: '100%', borderRight: 0 }}
          >
            <Menu.Item icon={<HomeOutlined />} key="/home">
              <Link to={'/home'}>文章列表</Link>
            </Menu.Item>
            <Menu.Item icon={<DiffOutlined />} key="/article">
              <Link to={'/article'}>内容管理</Link>
            </Menu.Item>
            <Menu.Item icon={<EditOutlined />} key="/publish">
              <Link to={'/publish'}>发布文章</Link>
            </Menu.Item>
          </Menu>
        </Sider>
        <Layout className="layout-content" style={{ padding: 20 }}>
          {/* 二级路由出口 */}
          <Outlet />
        </Layout>
      </Layout>
    </Layout>
  )
}

export default observer(AuthLayout)