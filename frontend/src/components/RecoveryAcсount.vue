<template>
  <v-form ref="form">
    <v-snackbar
        v-model="snackbar.show"
        :color="snackbar.color"
        :multi-line=true
        :timeout="10000"
    >{{snackbar.text}}</v-snackbar>
    <v-dialog v-model="show" persistent max-width="500" height="50px">
      <v-card>
        <v-toolbar>
          <v-toolbar-title class="headline">{{options.name}}</v-toolbar-title>
        </v-toolbar>
          <v-card-text>
              <v-layout wrap>
                <v-flex xs12 v-if ="isRadio">
                  <v-radio-group  v-model="radioGroup">
                    <v-radio label="по электронной почте" value="email"></v-radio>
                    <v-radio label="по телефону" value="phone"></v-radio>
                  </v-radio-group>
                </v-flex>
                <v-flex xs12 v-if = "isEmail">
                  <v-text-field v-model="email" prepend-icon="email" label="Электронная почта"
                                :rules="[rules.required, rules.email]"></v-text-field>
                </v-flex>
                <v-flex xs12 v-else>
                  <v-text-field v-model="phone" prepend-icon="phone" label="Мобильный телефон"
                                :rules="[rules.required, rules.phone]"></v-text-field>
                </v-flex>
                <template v-if = "code">
                  <v-flex xs8>
                    <v-text-field  label="Введите код" v-model="code"></v-text-field>
                  </v-flex>
                  <v-flex xs4>
                    <v-btn small color="success">Отправить код</v-btn>
                  </v-flex>
                </template>
              </v-layout>
          </v-card-text>
        <v-card-actions>
          <v-btn small color="primary" @click="onRecoveryPassword">Получить код</v-btn>
          <v-spacer></v-spacer>
          <v-btn small flat @click = "onClosed" color="primary">Закрыть</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>

<script>
    const axios = require('axios').default
    export default {
        name: "RecoveryAcсount",
        props: {
            show : {
                type: Boolean
            },
            options: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                radioGroup: "email",
                email: "",
                phone: "",
                code: "",
                recovery: {
                  showText: false,
                  showBtn: false,
                  url: ""
                },
                rules: {
                    required: value => !!value || 'поле не должно быть пустым',
                    email: value => {
                        const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
                        return (!value || pattern.test(value))  || 'не корректный Email.'
                    },
                    phone: value => {
                        const pattern =  /^(?!.{17,})(\s*)?(\+)?([- _():=+]?\d[- _():=+]?){11,14}(\s*)?/
                        return (!value || pattern.test(value) || 'телефон должен соответсвовать: 8 111 222 33-44')
                    }
                },
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
                ])
            }
        },
        computed: {
            isRadio () {
                if (this.options && this.options.email && this.options.email.urlReceiveCode && this.options.email.urlConfirmCode
                    && this.options.phone && this.options.phone.urlReceiveCode && this.options.phone.urlConfirmCode) {
                      return true
                }
              return false
            },
            isEmail() {
                if (this.isRadio) {
                    if (this.radioGroup === "email") {
                        return true
                    }
                    return  false
                }
                if (this.options && this.options.phone && this.options.phone.urlReceiveCode && this.options.phone.urlConfirmCode) {
                    return false
                }
                return true
            },
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
            validate () {
                return this.$refs.form.validate()
            }
        },
        methods: {
            showSnackBar(text, color) {
                this.snackbar.show = true
                this.snackbar.text = text
                this.snackbar.color = color || "error"
            },
            onClosed () {
                this.$emit('closed', false)
            },
            async onRecoveryPassword() {
                this.snackbar.show = false
                //общий метод, который должен определить восстановление идет по почте или телефону
                if (!this.validate) {
                    this.showSnackBar("Заполните корректно все поля формы", "warning")
                    return
                }
                if (!this.options) {
                    this.showSnackBar("Опции для восстановления акаунта не были переданы")
                    return
                }
                // const data = new FormData()
                // data.append("email", this.email)
                // try {
                //     const res = await axios.post(this.pas, data)
                //     if (res && res.data && res.data.error) {
                //         this.showSnackBar(res.data.error, "warning")
                //         return
                //     }
                //     this.snackbar.show = false
                //     this.recovery.showText = true
                //     let urlEmail = this.getEmail
                //     if (urlEmail) {
                //         this.recovery.showBtn = true
                //         this.recovery.url = urlEmail
                //     }
                // } catch (e) {
                //     if (e.response && e.response.data && e.response.data.error) {
                //         this.showSnackBar(e.response.data.error, "warning")
                //         return
                //     }
                //     this.showSnackBar(`не удалось восстановить пароль. Ошибка: ${e}`)
                // }
            }
        }
    }
</script>

<style scoped>

</style>