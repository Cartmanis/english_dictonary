<template>
  <v-app>
    <v-content>
      <!--<Login :controler="controler" show-password :registration="registration" @click-login="onLogin"/>-->
      <router-view/>
    </v-content>
  </v-app>
</template>

<script>
import Login from './components/Login'
import Main from './components/Main'

export default {
  name: 'App',
  // components: {
  //   Main,
  //   Login
  // },
  data () {
    return {
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
  methods: {
    onLogin(isAuth) {
        console.log(isAuth)
        if (isAuth) {
            console.log('попали')
            this.$router.push('/')
        }
      
    }
  }
}
</script>
