<template>
    <v-form ref="form">
        <v-snackbar
          v-model="snackbar.show"
          :color="snackbar.color"
          :multi-line=true
          :timeout="10000"
        >{{snackbar.text}}</v-snackbar>
        <v-dialog v-model="show" persistent :max-width="width">
            <v-card>
                <v-card-title>
                    <span class="headline">{{nameForm}}</span>
                </v-card-title>
                <v-card-text>
                    <v-container grid-list-md>
                        <form>
                            <v-layout wrap>
                                <v-flex xs12>
                                    <v-text-field :label="labelLogin" prepend-icon="person"
                                                  v-model="login"
                                                  :rules="[rules.requiredLogin, rules.maxLogin, rules.minLogin, rules.characterLogin]"
                                    ></v-text-field>
                                </v-flex>
                                <v-flex xs12>
                                    <v-text-field :label="labelPassword"  prepend-icon="lock"
                                                  :type = "password.show ? 'text' : 'password'"
                                                  v-model="password.text"
                                                  :append-icon="password.show ? 'visibility' : 'visibility_off'"
                                                  @click:append="password.show=!password.show"
                                                  :rules="[rules.requiredPassword, rules.passwordValid, rules.minPassword]">
                                    </v-text-field>
                                </v-flex>
                                <v-flex xs12 v-show="registration && registration.surname">
                                    <v-text-field :label="labelSurname" prepend-icon="account_circle"
                                                  v-model="surname"
                                                  :rules="[rules.requiredSurname, rules.LetterValid,
                                                  rules.minSurname, rules.maxSurname]">
                                    </v-text-field>
                                </v-flex>
                                <v-flex xs12 v-show="registration && registration.name">
                                    <v-text-field :label="labelName" prepend-icon="account_circle"
                                                  v-model="name"
                                                  :rules="[rules.requiredName, rules.LetterValid,
                                                  rules.minName, rules.maxName]">
                                    </v-text-field>
                                </v-flex>
                                <v-flex xs12 v-show="registration && registration.birthDay">
                                    <date-hidden :label="labelBirthDay" :required="checker.birthDay.required" @change-date="birthDay=$event"></date-hidden>
                                </v-flex>
                                <v-flex xs12 v-show="registration && registration.email">
                                    <v-text-field prepend-icon="email"
                                                  :label="checker.email.required ? 'Электронная почта*' : 'Электронная почта'"
                                                  v-model="email"
                                                  :rules="[rules.requiredEmail, rules.email]">
                                    </v-text-field>
                                </v-flex>
                                <v-flex xs12>
                                    <v-text-field prepend-icon="phone"
                                                  v-show="registration && registration.phone"
                                                  :label="checker.phone.required ? 'Телефон*' : 'Телефон'"
                                                  v-model="phone"
                                                  :rules="[rules.requiredPhone, rules.phone]">

                                    </v-text-field>
                                </v-flex>
                                <template>
                                    <v-flex xs12 v-show = "confirmEmail.showText">
                                        <v-label text-color ="success">На ваш электронный адрес отправлено письмо со ссылкой на подтверждение регистрации. Перейдите в почту для завершения регистрации.</v-label>
                                    </v-flex>
                                    <v-flex xs12 v-show="confirmEmail.showBtn" >
                                        <a target="_blank" :href="confirmEmail.url">
                                            <v-btn @click="onClosed"  color="primary">Перейти в почту</v-btn>
                                        </a>
                                    </v-flex>
                                </template>
                            </v-layout>
                        </form>
                    </v-container>
                    <small>*поля обязательные для заполнения</small>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="primary" flat @click="onClosed">Закрыть</v-btn>
                    <v-btn color="primary" flat @click="onRegistration">Сохранить</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-form>
</template>

<script>
    const axios = require('axios').default
    import dateHidden from './DateHidden'
    export default {
        name: "UserRegistration",
        components: {
            'date-hidden': dateHidden,
        },
        props: {
            show: {
              type: Boolean
            },
            registration: {
                type: Object
            },
            width: {
                type: String
            },
            controler: {
                type: Object
            }
        },
        computed: {
          nameForm() {
            if (this.registration && this.registration.nameForm) {
                return this.registration.nameForm
            }
            return "User Registration"
          },
           labelLogin() {
               let required = "*"
               if (this.registration && this.registration.login && this.registration.login.required === false) {
                   required = ""
               }
               if (this.registration && this.registration.login && this.registration.login.label) {
                   return  this.registration.login.label + required
               }
               return "User name" + required
           },
            labelPassword() {
                let required = "*"
                if (this.registration && this.registration.password && this.registration.password.required === false) {
                    required = ""
                }
                if (this.registration && this.registration.password && this.registration.password.label) {
                    return this.registration.password.label + required
                }
                return "Password" + required
            },
            labelSurname() {
              let required = ""
                if (this.registration && this.registration.surname && this.registration.surname.required) {
                    required = "*"
                }
                if (this.registration && this.registration.surname && this.registration.surname.label) {
                    return this.registration.surname.label + required
                }
                return "Surname" + required

            },
            labelName() {
                let required = ""
                if (this.registration && this.registration.name && this.registration.name.required) {
                    required = "*"
                }
                if (this.registration && this.registration.name && this.registration.name.label) {
                    return this.registration.name.label + required
                }
                return "Name" + required

            },
            labelBirthDay() {
              if (this.registration && this.registration.birthDay && this.registration.birthDay.label) {
                  return this.registration.birthDay.label
              }
              return ''
            },
          getParams () {
              let params = new Map();

              if (this.registration && this.registration.login && this.registration.login.name) {
                  params.set(this.registration.login.name, this.login)
              }
              if (this.registration && this.registration.password && this.registration.password.name) {
                  params.set(this.registration.password.name, this.password.text)
              }
              if (this.registration && this.registration.email && this.registration.email.name) {
                  params.set(this.registration.email.name, this.email)
              }
              if (this.registration && this.registration.phone && this.registration.phone.name) {
                  params.set(this.registration.phone.name, this.phone)
              }
              return params
          },
           validate () {
              return this.$refs.form.validate()
           }
        },
        methods: {
          onClosed() {
              this.$emit('closed', false)
          },
          async onRegistration() {
              if (!this.registration || !this.registration.controler) {
                  return
              }
              const url = `http://${this.registration.controler.ip}:${this.registration.controler.port}/${this.registration.controler.url}`
              if (!this.validate) {
                  this.showSnackBar("Сохранение не выполнено. Заполните корректно все поля формы", "warning")
                  return
              }
              const data = new FormData()
              this.getParams.forEach( (value, key) => {
                  if (value === "") {
                      return
                  }
                  data.append(key, value)
              });
              try {
                  const res = await axios.post(`${url}`, data)
                  if (res && res.data && res.data.error) {
                      this.showSnackBar(res.data.error, "warning")
                      return
                  }
                  this.snackbar.show = false
                  if (this.registration.email && this.registration.email.confirm) {
                      this.confirmEmail.showText = true
                      if (res && res.data && res.data.url) {
                          this.confirmEmail.showBtn = true
                          this.confirmEmail.url = res.data.url
                      }
                      return
                  }
                  this.onClosed()

              } catch (e) {
                  this.showSnackBar(`не удалось сохранить пользователя. Ошибка: ${e}`)
              }

          },
          showSnackBar(text, color) {
              this.snackbar.show = true
              this.snackbar.text = text
              this.snackbar.color = color || "error"
          }
        },
        data () {
            return {
                snackbar: {
                    show: false,
                    text: "",
                    color: ""
                },
                confirmEmail: {
                    showText: false,
                    showBtn: false,
                    url: "",
                },
                login:"",
                password: {
                    text:"",
                    show: false
                },
                surname:"",
                name:"",
                email:"",
                birthDay:"",
                phone:"",
                checker : {
                    login: {
                        required : (this.registration && this.registration.login &&
                            this.registration.login.required !== false) || this.registration.login.required,
                        min: (this.registration && this.registration.login && this.registration.login.min) || 3,
                        max: (this.registration && this.registration.login && this.registration.login.max)  || 20,
                        character: this.registration && this.registration.login && this.registration.login.character
                    },
                    password: {
                        required : (this.registration && this.registration.password &&
                            this.registration.password.required !== false) || this.registration.password.required,
                        min: (this.registration && this.registration.password && this.registration.password.min)  || 8
                    },
                    surname: {
                        required: this.registration && this.registration.surname && this.registration.surname.required,
                        min: (this.registration && this.registration.surname && this.registration.surname.min) || 3,
                        max: (this.registration && this.registration.surname && this.registration.surname.max)  || 30,
                    },
                    name: {
                        required: this.registration && this.registration.name && this.registration.name.required,
                        min: (this.registration && this.registration.name && this.registration.name.min) || 3,
                        max: (this.registration && this.registration.name && this.registration.name.max)  || 30,
                    },
                    birthDay: {
                        required: this.registration && this.registration.birthDay && this.registration.birthDay.required
                    },
                    email: {
                        required: this.registration && this.registration.email && this.registration.email.required
                    },
                    phone: {
                        required: this.registration && this.registration.phone && this.registration.phone.required
                    }
                },
                rules: {
                    required: value => !!value || 'поле не должно быть пустым',
                    requiredLogin: value => (!!value || !this.checker.login.required) || 'Имя пользователя не может быть пустым',
                    requiredPassword: value => (!!value || !this.checker.password.required) || 'Пароль не может быть пустым',
                    requiredEmail: value => (!!value || !this.checker.email.required) || 'Email не может быть пустым',
                    requiredPhone: value => (!!value || !this.checker.phone.required) || 'Телефон не может быть пустым',
                    requiredSurname: value => (!!value || !this.checker.surname.required) || 'Фамилия должна быть заполнена',
                    requiredName: value => (!!value || !this.checker.name.required) || 'Имя должно быть заполнено',
                    characterLogin: value => {
                        if (this.checker.login.character&& value) {
                            const pattern = /^[a-zA-Zа-яА-Я][0-9a-zA-Zа-яА-Я]+$/
                            // return (!this.checker.login.character || pattern.test((value))) || 'необходма латинская буква'
                            return pattern.test((value)) || 'имя пользователя может содержать цирфры и буквы, первый символ - буква'
                        }
                        return true
                    },
                    maxLogin: value => !value || value.length <= this.checker.login.max ||
                        `максимальное количество символов: ${this.checker.login.max}`,
                    minLogin: value => !value || value.length >= this.checker.login.min ||
                        `минимальное количество символов: ${this.checker.login.min}`,
                    LetterValid: value => {
                        if (value) {
                            const pattern = /^[a-zA-Zа-яА-Я]+$/
                            return pattern.test(value) || 'в данном поле могут присутствовать только буквы'
                        }
                        return true
                    },
                    maxSurname: value => !value || value.length <= this.checker.surname.max ||
                        `максимальное количество символов: ${this.checker.surname.max}`,
                    minSurname: value => !value || value.length >= this.checker.surname.min ||
                        `минимальное количество символов: ${this.checker.surname.min}`,
                    maxName: value => !value || value.length <= this.checker.name.max ||
                        `максимальное количество символов: ${this.checker.name.max}`,
                    minName: value => !value || value.length >= this.checker.name.min ||
                        `минимальное количество символов: ${this.checker.name.min}`,
                    passwordValid: value => {
                        if (value) {
                            const pattern = /(?=.*[0-9])(?=.*[a-zA-Zа-яА-Я])/
                            return pattern.test(value) || 'пароль должен содержать символы и цифры'
                        }
                        return true
                    },
                    minPassword: value => !value || value.length >= this.checker.password.min ||
                        `пароль не должен быть менее ${this.checker.password.min}  символов`,
                    email: value => {
                        const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
                        return (!value || pattern.test(value))  || 'не корректный Email.'
                    },
                    phone: value => {
                        const pattern =  /^(?!.{17,})(\s*)?(\+)?([- _():=+]?\d[- _():=+]?){11,14}(\s*)?/
                        return (!value || pattern.test(value) || 'телефон должен соответсвовать: 8 111 222 33-44')
                    }
                }
            }
        }
    }
</script>

<style scoped>

</style>
