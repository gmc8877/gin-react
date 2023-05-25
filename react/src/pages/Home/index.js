import img404 from '@/assets/error.png'
import './index.scss'
import { Link} from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { List, Avatar, Card, Breadcrumb, Form, Button, DatePicker, Select } from 'antd'
import 'moment/locale/zh-cn'
import locale from 'antd/es/date-picker/locale/zh_CN'
import React from 'react'
import { useEffect, useState } from 'react'
import { useStore } from '@/store'
import { http } from '@/utils'

const { Option } = Select
const { RangePicker } = DatePicker



const Home = () => {

  const { channelStore} = useStore()

  // 文章列表管理 统一管理数据 将来修改给setArticleData传对象
  const [articleData, setArticleData] = useState({
    list: [],// 文章列表
    count: 0 // 文章数量
  })

  // 文章参数管理
  const [params, setParams] = useState({
    page: 1,
    per_page: 10
  })
  // 获取文章列表
  useEffect(() => {
    const loadList = async () => {
      const res = await http.get('/home/articles', { params })
      const { results, total_count } = res.data
      setArticleData({
        list: results,
        count: total_count
      })
    }

    loadList()
  }, [params])
  

  /* 表单筛选功能实现 */
  const onFinish = (values) => {
    const { channel_id, date, status } = values

    // 数据处理
    const _params = {}
    // 格式化status
    if (status >= 0) {
      _params.status = status

    } else {
      delete params.status
    }

    // 初始化频道
    if (channel_id >= 0) {
      _params.channel_id = channel_id
    } else {
      delete params.channel_id
    }


    // 初始化时间
    if (date) {
      _params.begin_pubdate = date[0].format('YYYY-MM-DD')
      _params.end_pubdate = date[1].format('YYYY-MM-DD')
    } else {
      delete params.begin_pubdate
      delete params.end_pubdate
    }

    // 修改params数据 引起接口的重新发送 对象的合并是一个整体覆盖 改了对象的整体引用
    setParams({
      ...params,
      ..._params
    })
  }
  // 翻页实现
  const pageChange = (page) => {
    setParams({
      ...params,
      page
    })
  }
 
  return (
    <div>
      {/* 筛选区域 */}
      <Card
        title={
          <Breadcrumb separator=">">
            <Breadcrumb.Item>
              <Link to="/">首页</Link>
            </Breadcrumb.Item>
            <Breadcrumb.Item>文章</Breadcrumb.Item>
          </Breadcrumb>
        }
        style={{ marginBottom: 20 }}
      >
        <Form
          onFinish={onFinish}
        >
          <Form.Item label="频道" name="channel_id">
            <Select
              placeholder="请选择文章频道"
              style={{ width: 120 }}
            >
              <Option key={-1} value={-1}>全部</Option>
              {channelStore.channelList.map(channel => <Option key={channel.id} value={channel.id}>{channel.name}</Option>)}
            </Select>
          </Form.Item>

          <Form.Item label="日期" name="date">
            {/* 传入locale属性 控制中文显示*/}
            <RangePicker locale={locale}></RangePicker>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ marginLeft: 80 }}>
              筛选
            </Button>
          </Form.Item>
        </Form>
      </Card>
      {/* 文章列表区域 */}
      <Card title={`根据筛选条件共查询到 ${articleData.count} 条结果：`}>
        <List
          itemLayout="vertical"
          size='small'
          pagination={{
            onChange: pageChange,
            pageSize: 10,
          }}
          dataSource={articleData.list}
          footer={
            <div>
              
            </div>
          }
          renderItem={(item) => (
            <div className='contentlist'>
            <Link to="/content" state={
              {
                id: item.id,
              }}>
            <List.Item
              key={item.id}
              extra={
                <img
                  width={272}
                  alt="logo"
                  src={item.image}
                />
              }
            >
              <List.Item.Meta
                avatar={<Avatar src={img404} />}
                title={item.name}
                description={item.title}
              />
              {item.content}
            </List.Item>
            </Link>
            
            </div>
          )}
        />
      </Card>
    </div>
  )
}

export default observer(Home)