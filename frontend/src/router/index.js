import Vue from 'vue'
import VueRouter from 'vue-router'
import Main from '@/components/Main'
import Login from '@/components/Login'

Vue.use(VueRouter)

const axios = require('axios').default

const baseUrl = "http://192.168.0.84"

const isAuth = async (to, from, next) => {
    const url = `${baseUrl}:27333/api/v1/auth`
    try {
        const res = await axios.post(url, {}, {withCredentials:true})
        if (res && res.data && res.data.result) {
            next()
            return
        }
        next('/login')
    } catch (e) {
        next('/login')
    }

}

const routes = [
    {
      path: '/login',
      name: 'login',
      component: Login,
      props: {
          url : `${baseUrl}:27333/api/v1/login`,
          showPassword: false,
          recoveryAccount: {
              name: 'Восстановление пароля',
              url: `${baseUrl}:27333/api/v1/recovery_password`,
              phone: '',
              email: ''
          },
          registration: {
              nameForm: "Регистрация пользователя",
              url : `${baseUrl}:27333/api/v1/add_user`,
               login: {name:"login", label:"Имя пользователя", max: 14, character:true},
               password: {name:"password", label:"Password",  min: 6},
              // surname: {name:"fam", label:"Фамилия", required: true, min:2, max:15},
              // name: {name:"fam", label:"Имя", min:2, max:10},
              // birthDay: {name:"birth_day", label:"Дата рождения"},
              email: {name:"email", label: "Электронная почта", confirm: true, required: true},
              phone: {name:"phone", label:"Мобильный телефон", confirmUrl: `${baseUrl}:27333/api/v1/confirm_phone`}
          }
      }
    },
    {
        path: '/',
        name: 'main',
        component: Main,
        beforeEnter: isAuth,
    }
]

const router = new VueRouter({
    routes,
    mode: 'history'
})

export default router
