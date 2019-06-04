<template>
  <v-app>
    <v-content>
      <Login show-password :registration="registration" @click-login="onLogin"/>
    </v-content>
  </v-app>
</template>

<script>
import Login from './components/Login'

export default {
  name: 'App',
  components: {
    Login
  },
  data () {
    return {
      registration: {
        login: {max: 14, name:"login"},
        password: {min: 7, name:"password"},
        email: {required: true, name:"email"},
        phone: {required: false, name:"phone"}
      }
    }
  },
  methods: {
    onLogin(login, password) {
      console.log(login)
      console.log(password)
    }
  }
}
</script>
