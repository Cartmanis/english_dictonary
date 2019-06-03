<template>
  <v-app>
    <v-content>
      <Login show-password @click-login="onLogin"/>
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
      //
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
