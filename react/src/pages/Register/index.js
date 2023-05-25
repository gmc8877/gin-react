import './index.scss'
import { Card,Row, Col, Button, Checkbox, Form, Input, message } from 'antd'

import { useState } from 'react'
import { useStore } from '@/store'
import { useRef } from 'react'
import { useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { observer } from 'mobx-react-lite'

function Register () {
  const [form] = Form.useForm();
  const navigate = useNavigate()
  const { resStore } = useStore()
  const [loadings, setLoadings] = useState([])
  const timeRef = useRef()//设置延时器
  const [loadingtime, setLoadingtime] = useState(60)
 
  const enterLoading = (index) => {
      form.validateFields(['username','password','email']).then(async (v) => {
      await resStore.getCaptcha({
        mobile: v.username,
        email: v.email
      })
        if (resStore.errmsg !== 'OK') {
          message.error(resStore.errmsg)
        return
      }
      setLoadings((prevLoadings) => {
        const newLoadings = [...prevLoadings]
        newLoadings[index] = true
        return newLoadings
      })
      setLoadingtime(60)
      setTimeout(() => {
        setLoadings((prevLoadings) => {
          const newLoadings = [...prevLoadings]
          newLoadings[index] = false
          return newLoadings
        })
      }, 60000)
    })
  
  }
  useEffect(() => {
    //如果设置倒计时且倒计时不为0
    if (loadingtime && loadingtime !== 0)
      timeRef.current = setTimeout(() => {
        setLoadingtime(time => time - 1)
      }, 1000)
    //清楚延时器
    return () => {
      clearTimeout(timeRef.current)
    }
  }, [loadingtime])
  // useEffect(() => {
  //   if (resStore.errmsg !== 'OK') message.error(resStore.errmsg)
  // },[resStore])


  async function onFinish (values) {
    await resStore.getRegister({
      mobile: values.username,
      password: values.password,
      email: values.email,
      captcha:values.captcha
    })
    if (resStore.errmsg !== 'OK') {
      message.error(resStore.errmsg)
      return
    }
    navigate('/', { replace: true })
    message.success("注册成功")
  }
  
  return (
    <div className="register">
      <Card className="register-container">
        
        {/* 登录表单 */}
        <Form
          form={form}
          labelCol={{
            span: 4,
          }}
          initialValues={{
            remember: true,
          }}
          onFinish={onFinish}

        >
          <Form.Item
            label="手机号"
            name="username"
            rules={[
              {
                required: true,
                message: '请输入手机号',
              },
              {
                pattern: /^1[3-9]\d{9}$/,
                message: '请输入正确的手机号',
              },
            ]}
          >
            <Input size='large' placeholder='请输入手机号'/>
          </Form.Item>

          <Form.Item
            label="密码"
            name="password"
            rules={[
              {
                required: true,
                message: '请输入密码',
              },
              {
                len: 6,
                message: '请输入6位密码',
              },
            ]}
          >
            <Input size='large' placeholder='请输入密码' />
          </Form.Item>


          <Form.Item
            label="邮箱"
            name="email"
            rules={[
              {
                required: true,
                message: '请输入邮箱',
              },
              {
                // pattern: /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/,
                // required: true,
                type: 'email',
                message: '请输入正确的邮箱',
              },
            ]}
          >
            <Input size='large' placeholder='请输入邮箱'/>
          </Form.Item>

          <Form.Item label="验证码" >
            <Row gutter={8}>
              <Col span={12}>
                <Form.Item
                  name="captcha"
                  noStyle
                  rules={[
                    {
                      required: true,
                      message: '请输入验证码',
                    },
                    {
                      len: 6,
                      message: '请输入6位验证码',
                    },
                  ]}
                >
                  <Input placeholder='输入验证码' />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Button loading={loadings[0]} onClick={() => enterLoading(0)}>{!loadings[0] ? "发送验证码" : loadingtime + 's'}</Button>
              </Col>
            </Row>
          </Form.Item>

          <Form.Item>
            <Form.Item
              name='remember'
              valuePropName='checked'
              style={{
                display: 'inline-block',

              }}
            >
              <Checkbox className='register-checkbox-label'>我已阅读并同意</Checkbox>
            </Form.Item>
            <Form.Item
              style={{
                display: 'inline-block',
                position: 'absolute',
                right: '10%'
                
              }}
            >
              <Col>
                <Link className='register-link' to={"/login"}>返回登录页</Link>
              </Col>
            </Form.Item>
          </Form.Item>
          <Form.Item>
            <Button className='register-button' type="primary" htmlType="submit" size='large' block>
              注册
            </Button>

          </Form.Item>
        </Form>
      </Card>

    </div>
  )
}

export default observer(Register)