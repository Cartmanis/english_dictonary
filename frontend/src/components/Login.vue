<template>
  <v-container fluid fill-height>
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :multi-line=true
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
        switchLogin(typeAuth, login, password) {
            const url = `http://${this.controler.ip}:${this.controler.port}/${this.controler.url}`
            //для того, чтобы не было исключения для кирилических символов необходимо использовать обёртку unescape(encodeURIComponent(str)
            const credentials = btoa(unescape(encodeURIComponent(`${login}:${password}`)))

            switch (typeAuth.toLowerCase()) {
              case 'basic':
                return axios.post(url, {}, {
                  headers: {'Authorization': `Basic ${credentials}`},
                  withCredentials: true
                })
              default:
                return axios.post(url, {}, {
                  headers: {'Authorization': `Basic ${credentials}`},
                  withCredentials: true
                })
            }

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
          this.error = false
          try {
              const res = await this.switchLogin('basic', this.login, this.password.text)
              if (res && res.data && res.data.result) {
                  this.$emit('click-login', true)
              } else {
                  this.error = true
                  this.$emit('click-login', false)
              }
              
          } catch (e) {
              if (e && e.response && e.response.data) {
                  this.error = true
                  this.$emit('click-login', false)
                  return
              }
              this.showSnackBar(`не удалось произвести авторизацию. ${e}`)
          }
        }
      }
    }
</script>

<style scoped>

</style>
