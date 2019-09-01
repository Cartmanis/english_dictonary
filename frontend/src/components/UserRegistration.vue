<template>
  <v-form ref="form">
    <v-snackbar
        v-model="snackbar.show"
        :color="snackbar.color"
        :multi-line=true
        :timeout="7000"
    >{{snackbar.text}}</v-snackbar>
    <v-dialog v-model="show" persistent :max-width="maxWidth" height="50px">
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
                <v-flex xs12 v-show="options && options.surname">
                  <v-text-field :label="labelSurname" prepend-icon="account_circle"
                                v-model="surname"
                                :rules="[rules.requiredSurname, rules.LetterValid,
                                                  rules.minSurname, rules.maxSurname]">
                  </v-text-field>
                </v-flex>
                <v-flex xs12 v-show="options && options.name">
                  <v-text-field :label="labelName" prepend-icon="account_circle"
                                v-model="name"
                                :rules="[rules.requiredName, rules.LetterValid,
                                                  rules.minName, rules.maxName]">
                  </v-text-field>
                </v-flex>
                <v-flex xs12 v-show="options && options.birthDay">
                  <date-hidden :label="labelBirthDay" :required="checker.birthDay.required" @change-date="birthDay=$event"></date-hidden>
                </v-flex>
                <v-flex xs12 v-show="options && options.email">
                  <v-text-field :label="labelEmail" prepend-icon="email"
                                v-model="email"
                                :rules="[rules.requiredEmail, rules.email]">
                  </v-text-field>
                </v-flex>
                <v-flex xs12 v-show="options && options.phone">
                  <v-text-field :label="labelPhone" prepend-icon="phone"
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
                <template v-if = "confirmPhone.url">
                <v-flex xs12>
                  <v-label text-color ="success">На ваш телефон выслано sms c кодом подтверждения</v-label>
                </v-flex>
                <v-flex xs6>
                  <v-text-field v-model = "confirmPhone.code" label = 'Код подтверждения'></v-text-field>
                </v-flex>
                <v-flex xs6>
                    <v-btn @click="onConfirmPhone" color="primary">Подтвердить</v-btn>
                </v-flex>
                  <v-label v-show = "confirmPhone.text">{{confirmPhone.text}}</v-label>
                </template>
              </v-layout>
            </form>
          </v-container>
          <small>*поля обязательные для заполнения</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" flat @click="onClosed">Закрыть</v-btn>
          <v-btn color="primary" flat @click="onSave">Сохранить</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>

<script>
    const axios = require('axios').default
    import dateHidden from 'date-hidden'
    export default {
        name: "UserRegistration",
        components: {
            'date-hidden': dateHidden,
        },
        props: {
            show: {
                type: Boolean
            },
            options: {
                type: Object,
                required: true
            },
            maxWidth: {
                type: String
            },
            url: {
                type: String
            }
        },
        computed: {
            getEmail() {
              if (!this.email) {
                return ""
              }
              const arr = this.email.split('@')
              if (!arr || arr.length < 2) {
                return ""
              }
              return this.mapEmail.get(arr[1])
            },
            nameForm() {
                if (this.options && this.options.nameForm) {
                    return this.options.nameForm
                }
                return "User Registration"
            },
            labelLogin() {
                let required = "*"
                if (this.options && this.options.login && this.options.login.required === false) {
                    required = ""
                }
                if (this.options && this.options.login && this.options.login.label) {
                    return  this.options.login.label + required
                }
                return "User name" + required
            },
            labelPassword() {
                let required = "*"
                if (this.options && this.options.password && this.options.password.required === false) {
                    required = ""
                }
                if (this.options && this.options.password && this.options.password.label) {
                    return this.options.password.label + required
                }
                return "Password" + required
            },
            labelSurname() {
                let required = ""
                if (this.options && this.options.surname && this.options.surname.required) {
                    required = "*"
                }
                if (this.options && this.options.surname && this.options.surname.label) {
                    return this.options.surname.label + required
                }
                return "Surname" + required

            },
            labelName() {
                let required = ""
                if (this.options && this.options.name && this.options.name.required) {
                    required = "*"
                }
                if (this.options && this.options.name && this.options.name.label) {
                    return this.options.name.label + required
                }
                return "Name" + required

            },
            labelBirthDay() {
                if (this.options && this.options.birthDay && this.options.birthDay.label) {
                    return this.options.birthDay.label
                }
                return ''
            },
            labelEmail () {
                let required = ""
                if (this.options && this.options.email && this.options.email.required) {
                    required = "*"
                }
                if (this.options && this.options.email && this.options.email.label) {
                    return this.options.email.label + required
                }
                return "Email" + required
            },
            labelPhone () {
                let required = ""
                if (this.options && this.options.phone && this.options.phone.required) {
                    required = "*"
                }
                if (this.options && this.options.phone && this.options.phone.label) {
                    return this.options.phone.label + required
                }
                return "Phone" + required
            },
            getParams () {
                let params = new Map();

                if (this.options && this.options.login && this.options.login.name) {
                    params.set(this.options.login.name, this.login)
                }
                if (this.options && this.options.password && this.options.password.name) {
                    params.set(this.options.password.name, this.password.text)
                }
                if (this.options && this.options.email && this.options.email.name) {
                    params.set(this.options.email.name, this.email)
                }
                if (this.options && this.options.phone && this.options.phone.name) {
                    params.set(this.options.phone.name, this.phone)
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
            async onConfirmPhone() {
                if (!this.confirmPhone.url) {
                    this.showSnackBar("Сбой при отправке кода подтверждения. url для server пустой")
                    return
                }
                if (!this.confirmPhone.code) {
                    this.showSnackBar("Код подтверждения не отправлен. Заполните поле Код подтверждения", "warning")
                    return
                }
                try {
                    const data = new FormData()
                    data.append("code_phone", this.confirmPhone.code)
                    const res = await axios.post(`${this.confirmPhone.url}`, data)
                    if (res && res.data && res.data.error) {
                        this.showSnackBar(res.data.error, "warning")
                        return
                    }
                    this.confirmPhone.text = `Телефон успешно подтвержден`
                    return
                } catch (e) {
                    if (e && e.response && e.response.data && e.response.status) {
                        if (e.response.status === 401 || e.response.status === 403) {
                            this.confirmPhone.text = `Не удалось подтвердить телефон: неверный код подтверждения`
                        }
                        return
                    }
                    this.showSnackBar(`не удалось подтвердить телефон: ${e}`)
                }
              this.onClosed()
            },
            async onSave() {
                if (!this.options || !this.options.url) {
                    return
                }
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
                    const res = await axios.post(`${this.options.url}`, data)
                    if (res && res.data && res.data.error) {
                        this.showSnackBar(res.data.error, "warning")
                        return
                    }
                    this.snackbar.show = false

                    //Для электронной почты
                    let urlEmail = this.getEmail
                    if (this.options.email && this.options.email.confirm) {
                      this.confirmEmail.showText = true
                      if (urlEmail) {
                        this.confirmEmail.showBtn = true
                        this.confirmEmail.url = urlEmail
                      } else {
                        this.confirmEmail.showBtn = false
                        this.confirmEmail.url = ""
                      }
                    }
                    // Для телефона
                  if (this.options.phone && this.options.phone.confirmUrl) {
                    this.confirmPhone.url = this.options.phone.confirmUrl
                  }

                     // Проверка на автоматическое закрытие формы
                    if (!this.confirmEmail.showText && !this.confirmPhone.url) {
                        this.onClosed()
                    }

                } catch (e) {
                    if (e && e.response && e.response.data && e.response.data.error) {
                        this.showSnackBar(`не удалось сохранить пользователя. Ошибка: ${e.response.data.error}`)
                        return
                    }
                    this.showSnackBar(`не удалось сохранить пользователя: ${e}`)
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
                mapEmail : new Map([
                  ["mail.ru",        "https://e.mail.ru/"],
                  ["bk.ru",         "https://e.mail.ru/"],
                  ["list.ru",        "https://e.mail.ru/"],
                  ["inbox.ru",       "https://e.mail.ru/"],
                  ["yandex.ru",      "https://mail.yandex.ru/"],
                  ["ya.ru",          "https://mail.yandex.ru/"],
                  ["yandex.ua",      "https://mail.yandex.ua/"],
                  ["yandex.by",      "https://mail.yandex.by/"],
                  ["yandex.kz",      "https://mail.yandex.kz/"],
                  ["yandex.com",     "https://mail.yandex.com/"],
                  ["gmail.com",      "https://mail.google.com/"],
                  ["googlemail.com", "https://mail.google.com/"],
                  ["outlook.com",    "https://mail.live.com/"],
                  ["hotmail.com",    "https://mail.live.com/"],
                  ["live.ru",        "https://mail.live.com/"],
                  ["live.com",       "https://mail.live.com/"],
                  ["me.com",         "https://www.icloud.com/"],
                  ["icloud.com",     "https://www.icloud.com/"],
                  ["rambler.ru",     "https://mail.rambler.ru/"],
                  ["yahoo.com",      "https://mail.yahoo.com/"],
                  ["ukr.net",        "https://mail.ukr.net/"],
                  ["i.ua",           "http://mail.i.ua/"],
                  ["bigmir.net",     "http://mail.bigmir.net/"],
                  ["tut.by",         "https://mail.tut.by/"],
                  ["inbox.lv",       "https://www.inbox.lv/"],
                  ["mail.kz",        "http://mail.kz/"]
                ]),
                confirmEmail: {
                    showText: false,
                    showBtn: false,
                    url: "",
                },
                confirmPhone: {
                    url: "",
                    code: "",
                    text: ""
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
                        required : (this.options && this.options.login &&
                            this.options.login.required !== false) || this.options.login.required,
                        min: (this.options && this.options.login && this.options.login.min) || 3,
                        max: (this.options && this.options.login && this.options.login.max)  || 20,
                        character: this.options && this.options.login && this.options.login.character
                    },
                    password: {
                        required : (this.options && this.options.password &&
                            this.options.password.required !== false) || this.options.password.required,
                        min: (this.options && this.options.password && this.options.password.min)  || 8
                    },
                    surname: {
                        required: this.options && this.options.surname && this.options.surname.required,
                        min: (this.options && this.options.surname && this.options.surname.min) || 3,
                        max: (this.options && this.options.surname && this.options.surname.max)  || 30,
                    },
                    name: {
                        required: this.options && this.options.name && this.options.name.required,
                        min: (this.options && this.options.name && this.options.name.min) || 3,
                        max: (this.options && this.options.name && this.options.name.max)  || 30,
                    },
                    birthDay: {
                        required: this.options && this.options.birthDay && this.options.birthDay.required
                    },
                    email: {
                        required: this.options && this.options.email && this.options.email.required
                    },
                    phone: {
                        required: this.options && this.options.phone && this.options.phone.required
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
