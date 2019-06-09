import Vue from 'vue'
import VueRouter from 'vue-router'
import Main from '@/components/Main'
import Login from '@/components/Login'

Vue.use(VueRouter)

const isAuth = (to, from, next) => {
    const authUser = JSON.parse(window.localStorage.getItem('login'))

    if (authUser) {
        next()
        return
    }
    next('/login')
}

const routes = [
    {
      path: '/login',
      name: 'login',
      component: Login,
      props: {
          controler: {
              ip: window.location.hostname,
              port: 27333,
              url: "api/v1/login"
          },
          registration: {
              controler :{
                  ip: window.location.hostname,
                  port: 27333,
                  url: "api/v1/add_user"
              },
              login: {max: 14, character:true, name:"login"},
              password: {min: 6, name:"password"},
              email: {required: true, name:"email"},
              phone: {required: false, name:"phone"}
          }
      }
    },
    {
        path: '/',
        name: 'main',
        component: Main,
        beforeEnter: isAuth
    }
]

const router = new VueRouter({
    routes,
    mode: 'history'
})

export default router
