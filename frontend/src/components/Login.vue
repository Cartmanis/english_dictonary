<template>
  <v-container fluid fill-height>
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="10000"
    >{{snackbar.text}}</v-snackbar>
    <v-layout align-center justify-center>
      <v-flex xs10 sm8 md5 lg4>
        <v-card v-show="showAuth"  class="elevation-12">
          <v-toolbar>
            <v-toolbar-title>Авторизация</v-toolbar-title>
            <v-spacer></v-spacer>
          </v-toolbar>
          <v-card-text>
            <v-form>
              <v-text-field  prepend-icon="person" label="Имя пользователя"
                             v-model = "login"
              ></v-text-field>
              <v-text-field v-if = "showPassword"
                      v-model='password.text' prepend-icon="lock" label="Пароль"
                      :type= "password.show? 'text' :'password'"
                      :append-icon="password.show ? 'visibility' : 'visibility_off'"
                      @click:append="password.show=!password.show">
              </v-text-field>
              <v-text-field v-else type ="password" v-model='password.text'
                            prepend-icon="lock" label="Пароль"></v-text-field>
              <h4 v-if="error" align="center" style="color: red;">имя пользователя или пароль введены не верно</h4>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn @click ="onLogin" color="primary">Вход</v-btn>
            <v-spacer></v-spacer>
            <v-btn flat @click = "onClosed" color="primary">Закрыть</v-btn>
            <v-btn flat v-show = "registration" @click="onRegistration" color="info">Регистрация</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
    <user-registration width="620px" :registration="registration"
                       :show=userRegistration  @closed="userRegistration=$event"></user-registration>
  </v-container>
</template>
<script>
  const axios = require('axios').default
  import UserRegistration from './UserRegistration'
    export default {
        name: "Login",
        components: {
          UserRegistration
        },
      data () {
        return {
            snackbar: {
                show: false,
                text: "",
                color: ""
            },
            showAuth:true,
            userRegistration : false,
            error: false,
            login: "",
            password: {
                text: "",
                show: ""
            },
        }
      },
      props: {
          showPassword: {
            type: Boolean
          },
          registration: {
            type: Object
          },
          controler: {
              type: Object
          }
      },
      methods: {
        showSnackBar(text, color) {
            this.snackbar.show = true
            this.snackbar.text = text
            this.snackbar.color = color || "error"
        },
        onClosed () {
          this.showAuth = false
        },
        onRegistration () {
          this.userRegistration = true
        },
        async onLogin () {
          if (!this.controler) {
              return
          }
          const url = `http://${this.controler.ip}:${this.controler.port}/${this.controler.url}`
            auth = {
              login: this.login,
              password: this.password
            }
            const data = new FormData()
            data.append('login', this.login)
            data.append('password', this.password.text)
          try {
              const res = await axios.post(`${url}`, data)
              console.log(res)
          } catch (e) {
              this.showSnackBar(`не удалось произвести авторизацию. Ошибка:${e}`)
              return
          }
          this.$emit('click-login', this.login, this.password.text)
        },
        async logIn () {
          const res = await controller.login(this.auth.login, this.auth.password)
          console.log(res.data.error, ' - ', res.data.is_login)
          if (!res.data.error && res.data.is_login === 'true') {
            window.localStorage.setItem('login', true)
            this.$router.push('/admin')
          } else {
            this.error = true
          }
        }
      }
    }
</script>

<style scoped>

</style>
