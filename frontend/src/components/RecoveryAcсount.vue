<template>
  <v-container v-show="show"  fluid fill-height>
    <v-snackbar
        v-model="snackbar.show"
        :color="snackbar.color"
        :multi-line=true
        :timeout="10000"
    >{{snackbar.text}}</v-snackbar>
    <v-layout align-center justify-center>
      <v-flex xs10 sm8 md5 lg4>
        <v-card  class="elevation-12">
          <v-toolbar>
            <v-toolbar-title>{{options.name}}</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form ref="form">
              <v-text-field v-model="email" prepend-icon="email" label="Электронная почта"
                            :rules="[rules.required, rules.email]"></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn @click="onRecoveryPassword">Отправить пароль</v-btn>
            <v-spacer></v-spacer>
            <v-btn small flat @click = "onClosed" color="primary">Закрыть</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
<!--  <v-dialog v-model="showDialog" persistent max-width="700">-->
<!--    <v-card>-->
<!--      <v-card-title>-->
<!--        <span class="headline">Восстановление акаунта</span>-->
<!--      </v-card-title>-->
<!--      <v-card-text>-->
<!--        <v-form ref="form">-->
<!--          <template>-->
<!--            <v-label small text-color="primary">Для восстановления пароля, необходимо заполнить электронную почту,-->
<!--              которая была указана при регистрации и нажать на кнопку "Отправить пароль". На эту почту будет выслан новый пароль.</v-label>-->
<!--            <v-flex xs12>-->
<!--              <v-text-field prepend-icon="email"-->
<!--                            v-model="email"-->
<!--                            label="Электронная почта"-->
<!--                            :rules="[rules.required, rules.email]">-->
<!--              </v-text-field>-->
<!--            </v-flex>-->
<!--            <v-btn @click="onRecoveryPassword">Отправить пароль</v-btn>-->
<!--            <v-btn @click="showDialog=false;snackbar.show=false">Закрыть</v-btn>-->
<!--            <template>-->
<!--              <v-flex xs12 v-show = "recoveryPassword.showText">-->
<!--                <v-label text-color ="success">На ваш электронный адрес отправлено письмо со ссылкой на подтверждение регистрации. Перейдите в почту для завершения регистрации.</v-label>-->
<!--              </v-flex>-->
<!--              <v-flex xs12 v-show="recoveryPassword.showBtn" >-->
<!--                <a target="_blank" :href="recoveryPassword.url">-->
<!--                  <v-btn @click="showDialog=false;snackbar.show=false"  color="primary">Перейти в почту</v-btn>-->
<!--                </a>-->
<!--              </v-flex>-->
<!--            </template>-->
<!--          </template>-->

<!--        </v-form>-->
<!--      </v-card-text>-->
<!--    </v-card>-->
<!--  </v-dialog>-->
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
                email: "",
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
                if (!this.validate) {
                    this.showSnackBar("Заполните корректно все поля формы", "warning")
                    return
                }
                if (!this.options) {
                    this.showSnackBar("Опции для восстановления акаунта не были переданы")
                    return
                }
                const data = new FormData()
                data.append("email", this.email)
                try {
                    const res = await axios.post(this.pas, data)
                    if (res && res.data && res.data.error) {
                        this.showSnackBar(res.data.error, "warning")
                        return
                    }
                    this.snackbar.show = false
                    this.recovery.showText = true
                    let urlEmail = this.getEmail
                    if (urlEmail) {
                        this.recovery.showBtn = true
                        this.recovery.url = urlEmail
                    }
                } catch (e) {
                    if (e.response && e.response.data && e.response.data.error) {
                        this.showSnackBar(e.response.data.error, "warning")
                        return
                    }
                    this.showSnackBar(`не удалось восстановить пароль. Ошибка: ${e}`)
                }
            }
        }
    }
</script>

<style scoped>

</style>