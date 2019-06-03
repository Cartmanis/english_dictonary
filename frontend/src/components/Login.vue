<template>
  <v-container fluid fill-height>
      <v-dialog v-model="dialog" persistent max-width="500px">
        <v-card>
          <v-card-title>
            <span class="headline">Регистрация пользователя</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout wrap>
                <v-flex xs12>
                  <v-text-field label="Имя пользователя*"
                                :rules="[rules.required, rules.maxLogin, rules.minLogin]"
                  ></v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field label="Пароль*" type = "password"
                                :rules="[rules.required, rules.passwordValid, rules.minPassword]">
                  </v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field v-show="registration.email" label="Email*"
                            :rules="[rules.requiredEmail, rules.email]">
                  </v-text-field>
                </v-flex>
                <v-flex xs12>
                  <v-text-field v-show="registration.phone" label="Phone"></v-text-field>
                </v-flex>
              </v-layout>
            </v-container>
            <small>*поля обязательные для заполнения</small>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat @click="dialog = false">Закрыть</v-btn>
            <v-btn color="blue darken-1" flat @click="dialog = false">Сохранить</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
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
            <v-btn @click = "onClosed" color="primary">Закрыть</v-btn>
            <v-btn v-show = "registration" @click="onRegistration" flat color="info">Регистрация</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>
<script>
    export default {
        name: "Login",
      data () {
        return {
          showAuth: true,
          error: false,
          dialog: false,
          login: "",
          password: {
            text: "",
            show: ""
          },
          checker : {
            login: {
              min: (this.registration.login && this.registration.login.min) || 3,
              max: (this.registration.login && this.registration.login.max)  || 20,
            },
            password: {
              min: (this.registration.password && this.registration.password.min)  || 8
            },
            email: {
              required: (this.registration.email && this.registration.email.required) || false
            }
          },
          rules: {
            required: value => !!value || 'поле не должно быть пустым',
            requiredEmail: value => this.checker.email.required || 'Email не может быть пустым',
            maxLogin: value => !value || value.length <= this.checker.login.max ||
              `максимальное количество символов: ${this.checker.login.max}`,
            minLogin: value => !value || value.length >= this.checker.login.min ||
              `минимальное количество символов: ${this.checker.login.min}`,

            passwordValid: value => {
              const pattern = /(?=.*[0-9])(?=.*[a-zA-Zа-яА-Я])/
              return pattern.test(value) || 'пароль должен содержать символы и цифры'
            },
            minPassword: value => !value || value.length >= this.checker.password.min ||
              `пароль не должен быть менее ${this.checker.password.min}  символов`,
            email: value => {
              const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
              return pattern.test(value) || 'не вернный Email.'
            }
          }
        }
      },
      props: {
          showPassword: {
            type: Boolean
          },
          registration: {
            type: Object
          }
      },
      methods: {
        onClosed () {
          this.showAuth = false
        },
        onRegistration () {
          console.log(this.registration)
          this.dialog = true
        },
        onLogin () {
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
