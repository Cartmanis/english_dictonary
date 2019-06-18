<template>
  <v-container fluid fill-height>
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :multi-line=true
      :timeout="10000"
    >{{snackbar.text}}</v-snackbar>

    <v-dialog v-model="showDialog" persistent max-width="700">
      <v-card>
        <v-card-title>
          <span class="headline">Отправить пароль</span>
        </v-card-title>
        <v-card-text>
          <v-form ref="form">
            <v-label small text-color="primary">Для восстановления пароля, необходимо заполнить электронную почту,
              которая была указана при регистрации и нажать на кнопку "Отправить пароль". На эту почту будет выслан новый пароль.</v-label>
            <v-flex xs12>
              <v-text-field prepend-icon="email"
                            v-model="email"
                            label="Электронная почта"
                            :rules="[rules.required, rules.email]">
              </v-text-field>
            </v-flex>
            <v-btn @click="onRecoveryPassword">Отправить пароль</v-btn>
            <v-btn @click="showDialog=false;snackbar.show=false">Закрыть</v-btn>
            <template>
              <v-flex xs12 v-show = "recoveryPassword.showText">
                <v-label text-color ="success">На ваш электронный адрес отправлено письмо со ссылкой на подтверждение регистрации. Перейдите в почту для завершения регистрации.</v-label>
              </v-flex>
              <v-flex xs12 v-show="recoveryPassword.showBtn" >
                <a target="_blank" :href="recoveryPassword.url">
                  <v-btn @click="showDialog=false;snackbar.show=false"  color="primary">Перейти в почту</v-btn>
                </a>
              </v-flex>
            </template>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-layout align-center justify-center>
      <v-flex xs10 sm8 md5 lg4>
        <v-card v-show="showAuth"  class="elevation-12">
          <v-toolbar>
            <v-toolbar-title>Авторизация</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn v-show="passwordRecovery" @click="showDialog=true" small flat color="primary">Забыли пароль?</v-btn>
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
    <user-registration width="700px" :registration="registration"
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
      props: {
        showPassword: {
          type: Boolean
        },
        passwordRecovery: {
          type: String,
        },
        registration: {
          type: Object
        },
        controler: {
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
            recoveryPassword: {
              showText: false,
              showBtn: false,
              url: "",
            },
            showAuth:true,
            showDialog: false,
            userRegistration : false,
            errorAuth: {
              show: false,
              text: "имя пользователя или пароль введены не верно"
            } ,
            login: "",
            password: {
                text: "",
                show: ""
            },
            email: "",

            rules: {
              required: value => !!value || 'поле не должно быть пустым',
              email: value => {
                const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
                return (!value || pattern.test(value))  || 'не корректный Email.'
              }
            }
        }
      },
      methods: {
        showSnackBar(text, color) {
            this.snackbar.show = true
            this.snackbar.text = text
            this.snackbar.color = color || "error"
        },

        async onRecoveryPassword() {
          if (!this.validate) {
            this.showSnackBar("Заполните корректно все поля формы", "warning")
            return
          }
          const data = new FormData()
          data.append("email", this.email)
          try {
            const res = await axios.post(this.passwordRecovery, data)
            if (res && res.data && res.data.error) {
              this.showSnackBar(res.data.error, "warning")
              return
            }
            this.snackbar.show = false
            this.recoveryPassword.showText = true
            if (res && res.data && res.data.url) {
              this.recoveryPassword.showBtn = true
              this.recoveryPassword.url = res.data.url
            }
          } catch (e) {
            if (e.response && e.response.data && e.response.data.details) {
              this.showSnackBar(e.response.data.details, "warning")
              return
            }
            this.showSnackBar(`не удалось восстановить пароль. Ошибка: ${e}`)
          }

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
