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
            <v-btn v-show="recoveryAccount" @click="onRecoveryAccount" small flat color="primary">Забыли пароль?</v-btn>
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
              <h4 v-if="errorAuth.show" align="center" style="color: red;">{{errorAuth.text}}</h4>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn @click ="onLogin" color="primary">Вход</v-btn>
            <v-spacer></v-spacer>
            <v-btn small flat @click = "onClosed" color="primary">Закрыть</v-btn>
            <v-btn small flat v-show = "registration" @click="onRegistration" color="info">Регистрация</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
    <recovery-account :options="recoveryAccount" :show=showRecovery @closed="showRecovery=$event;showAuth=true"></recovery-account>
    <template v-if = "registration">
      <user-registration max-width="700px" :options="registration" :show=showUserRegistration  @closed="showUserRegistration=$event">
      </user-registration>
    </template>

  </v-container>
</template>
<script>
  const axios = require('axios').default
  // import UserRegistration from 'vuetify-user-registration'
  import UserRegistration from './UserRegistration'
  import RecoveryAccount from './RecoveryAcсount'
    export default {
        name: "Login",
        components: {
          UserRegistration,
          RecoveryAccount
        },
      props: {
        url: {
            type: String
        },
        showPassword: {
            type: Boolean
        },
        recoveryAccount: {
          type: Object
        },
        registration: {
          type: Object
        }
      },
      computed: {
          validate () {
            return this.$refs.form.validate()
          }
      },
      data () {
        return {
            snackbar: {
                show: false,
                text: "",
                color: ""
            },
            showAuth:true,
            showDialog: false,
            showUserRegistration : false,
            showRecovery: false,
            errorAuth: {
              show: false,
              text: "имя пользователя или пароль введены не верно"
            } ,
            login: "",
            password: {
                text: "",
                show: ""
            }
        }
      },
      methods: {
          showSnackBar(text, color) {
              this.snackbar.show = true
              this.snackbar.text = text
              this.snackbar.color = color || "error"
          },
          switchLogin(typeAuth, login, password) {
              //для того, чтобы не было исключения для кирилических символов необходимо использовать обёртку unescape(encodeURIComponent(str)
              const credentials = btoa(unescape(encodeURIComponent(`${login}:${password}`)))

              switch (typeAuth.toLowerCase()) {
                  case 'basic':
                      return axios.post(this.url, {}, {
                          headers: {'Authorization': `Basic ${credentials}`},
                          withCredentials: true
                      })
                  default:
                      return axios.post(this.url, {}, {
                          headers: {'Authorization': `Basic ${credentials}`},
                          withCredentials: true
                      })
              }

          },
          onClosed() {
              this.showAuth = false
          },
          onRegistration() {
              this.showUserRegistration = true
          },
          onRecoveryAccount() {
              this.showRecovery = true
              this.showAuth = false
          },
          async onLogin() {
              if (!this.url) {
                  this.showSnackBar("не был передан обязательный параметр url в компонент Login")
                  return
              }
              this.showSnackBar()
              if (!this.login || !this.password || !this.password.text) {
                  this.showSnackBar("заполните имя пользователя и пароль", "warning")
                  return
              }
              this.snackbar.show = false
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
                  if (e && e.response && e.response.data && e.response.status) {
                      this.errorAuth.show = true
                      if (e.response.status === 403) {
                          this.errorAuth.text = `Пользователь ${this.login} не подтвердил адрес электронной почты. Необходимо перейти в электронную почту, указанную при регистрации и завершить регистрацию`
                      }
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
