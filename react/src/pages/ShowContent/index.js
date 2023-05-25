
import { Avatar, Card } from 'antd'

import { observer } from 'mobx-react-lite'
import './index.scss'
// import 'react-quill/dist/quill.snow.css'
import { useLocation } from 'react-router-dom'
import { useEffect,useState } from 'react'
import { http } from '@/utils'


const { Meta } = Card;

const Content= () => {
  const [data, setData] = useState({})
  const [url, setUrl] = useState()
  const location = useLocation()
  const { state: { id } } = location
  useEffect(() => {
    const loadDetail = async () => {
      const res = await http.get(`/content/articles/${id}`)
      setData(res.data)
      setUrl(res.data.cover.images[0])
    }
    loadDetail()
    
  }, [id])
  const { content } = data
  return (
    <Card
      className='card'
      cover={
        <img 
          className='cardimg'
          alt="example"
          src={url}
        />
      }
      
    >
      <Meta
        className='cardmeta'
        avatar={<Avatar src="https://xsgames.co/randomusers/avatar.php?g=pixel" />}
        title={data.title}
      />
      <div className='content' dangerouslySetInnerHTML={{__html:content}}/>
    </Card>
  );
  
}

export default observer(Content)