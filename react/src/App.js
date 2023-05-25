
import { unstable_HistoryRouter as HistoryRouter, Route, Routes } from 'react-router-dom'
import { history } from './utils/history'
import { AuthComponent } from '@/components/AuthComponent'
import { RootComponent } from '@/components/RootComponent'
import { lazy, Suspense } from 'react'
import './App.css'

const Login = lazy(() => import('./pages/Login'))
const Layout = lazy(() => import('./pages/Layout'))
const Home = lazy(() => import('./pages/Home'))
const Article = lazy(() => import('./pages/Article'))
const Publish = lazy(() => import('./pages/Publish'))
const Register = lazy(() => import('./pages/Register'))
const ShowHome = lazy(() => import('./pages/ShowHome'))
const Content = lazy(() => import('./pages/ShowContent'))
const RootLogin = lazy(() => import('./pages/RootLogin'))
const AuthLayout = lazy(() => import('./pages/AuthLayout'))


function App () {

  return (
    <>
      {/* 注册路由 */}
      <HistoryRouter history={history}>
        <div className="App">
          <Suspense
            fallback={
              <div
                style={{
                  textAlign: 'center',
                  marginTop: 200
                }}
              >
                loading...
              </div>
            }
          >
            <Routes>
              <Route path="/" element={
                <AuthComponent>
                  <AuthLayout />
                </AuthComponent>
              }>
                {/* 二级路由默认页面 */}

                <Route path="/home" element={<Home />} />
                <Route path="/article" element={<Article />} />
                <Route path="/publish" element={<Publish />} />
              
              </Route>

              <Route path="/login" element={<Login />}></Route>
              <Route path="/root_login" element={<RootLogin />}></Route>
              <Route path="/register" element={<Register />}></Route>
              <Route path="/show_home" element={<ShowHome />} />

              <Route path="/content" element={
                <AuthComponent>
                  <Content />
                </AuthComponent>
              }>
              </Route>

              <Route path="/root" element={
                <RootComponent>
                  <Layout />
                </RootComponent>
              }>

                <Route path="home_root" element={<Home />} />
                <Route path="article_root" element={<Article />} />
                <Route path="publish_root" element={<Publish />} />
                
              </Route>
            </Routes>




          </Suspense>
        </div>
      </HistoryRouter>
    </>

  )

}

export default App