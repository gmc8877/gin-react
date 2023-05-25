import { Layout, Menu, Popconfirm } from 'antd'
import {
  HomeOutlined,
  DiffOutlined,
  EditOutlined,
  LogoutOutlined
} from '@ant-design/icons'
import './index.scss'
import { Outlet, Link, useLocation, useNavigate} from 'react-router-dom'
import { useStore } from '@/store'
import { useEffect } from 'react'
import {observer} from 'mobx-react-lite'


const { Header, Sider } = Layout

const GeekLayout = () => {
  const { pathname } = useLocation()
  const { userStore, loginStore ,channelStore} = useStore()
  const navigate = useNavigate()
  useEffect(() => {
    userStore.getRootInfo()
    channelStore.loadChannelList()
    }, [userStore,channelStore]);
  //确定退出
  const onConfirm = () => {
    //删除token
    loginStore.rootlogout()
    navigate('/root_login')
  }
  
  return (
    <Layout>
      <Header className="header">
        <div className="logo" />
        <div className="user-info">
          <span className="user-name">{ userStore.rootInfo}</span>
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
        <Sider width={200} className="site-layout-background">
          {/* 高亮原理 defaultSelectedKeys===item key*/}
          {/* 获取当前激活的path路径 */}
          <Menu
            mode="inline"
            theme="dark"
            selectedKeys={pathname}
            style={{ height: '100%', borderRight: 0 }}
          >
            <Menu.Item icon={<HomeOutlined />} key="/root/home">
              <Link to={'/root/home_root'}>文章列表</Link>
            </Menu.Item>
            <Menu.Item icon={<DiffOutlined />} key="/root/article">
              <Link to={'/root/article_root'}>内容管理</Link>
            </Menu.Item>
            <Menu.Item icon={<EditOutlined />} key="/root/publish">
              <Link to={'/root/publish_root'}>发布文章</Link>
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

export default observer(GeekLayout)