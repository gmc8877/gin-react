import { Layout} from 'antd';
import './index.scss'
import Home from '../Home';

const { Header, Footer, Content } = Layout

const ShowHome = () => {
  return (
    <Layout>
      <Header className="header"></Header>
      <Content className="content">
        <Home/>
      </Content>
      <Footer className="footerStyle"></Footer>
    </Layout>
  )
}
export default ShowHome