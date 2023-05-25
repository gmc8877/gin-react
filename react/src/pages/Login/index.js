import { Card, Button, Checkbox, Form, Input, message, Col } from 'antd'
import './index.scss'
import { useStore } from '@/store'
import { useNavigate, Link } from 'react-router-dom'


function Login () {
  const { loginStore } = useStore()
  const navigate = useNavigate()


  async function onFinish (values) {
    await loginStore.getToken({
      mobile: values.username,
      code: values.password
    })
    if (loginStore.message !== 'OK') {
      message.error(loginStore.message)
    } else {
      navigate('/', { replace: true })
      message.success("登录成功")
    }
   
  }

  return (
    <div className="login">
      <Card className="login-container">
        {/* <img className="login-logo" src={logo} alt="" /> */}
        {/* 登录表单 */}
        <Form
          initialValues={{
            remember: true,
          }}
          onFinish={onFinish}

        >
          <Form.Item
            // label="Username"
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
            <Input size='large' placeholder='请输入手机号' />
          </Form.Item>

          <Form.Item
            // label="Password"
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
          <Form.Item>
            <Form.Item
              name='remember'
              valuePropName='checked'
              style={{
                display: 'inline-block',
               
              }}
            >
              <Checkbox className='login-checkbox-label'>我已阅读并同意</Checkbox>
            </Form.Item>
            <Form.Item
              style={{
                display: 'inline-block',
                position: 'absolute',
                right:'10%'
               
              }}
            >
              <Col>
                <Link className='register-link' to={"/register"}>注册新用户</Link>
              </Col>
            </Form.Item>
          </Form.Item>
          <Form.Item>
            <Button className='login-button' type="primary" htmlType="submit" size='large' block>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>

    </div>
  )
}

export default Login