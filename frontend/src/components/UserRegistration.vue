<template>
    <v-dialog v-model="show" persistent :max-width="width">
        <v-card>
            <v-card-title>
                <span class="headline">Регистрация пользователя</span>
            </v-card-title>
            <v-card-text>
                <v-container grid-list-md>
                    <v-layout wrap>
                        <v-flex xs12>
                            <v-text-field label="Имя пользователя*" prepend-icon="person"
                                          v-model="login"
                                          :rules="[rules.required, rules.maxLogin, rules.minLogin]"
                            ></v-text-field>
                        </v-flex>
                        <v-flex xs12>
                            <v-text-field label="Пароль*" prepend-icon="lock"
                                          :type = "password.show ? 'text' : 'password'"
                                          v-model="password.text"
                                          :append-icon="password.show ? 'visibility' : 'visibility_off'"
                                          @click:append="password.show=!password.show"
                                          :rules="[rules.required, rules.passwordValid, rules.minPassword]">
                            </v-text-field>
                        </v-flex>
                        <v-flex xs12>
                            <v-text-field prepend-icon="email"
                                          v-show="registration && registration.email"
                                          :label="checker.email.required ? 'Email*' : 'Email'"
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
                    </v-layout>
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
</template>

<script>
    export default {
        name: "UserRegistration",
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
          getParams () {
              let params = [] //map key, value
              if (this.registration && this.registration.login && this.registration.login.name) {
                  params.push(this.registration.login.name)
              }
              if (this.registration && this.registration.password && this.registration.password.name) {
                  params.push(this.registration.password.name)
              }
              if (this.registration && this.registration.email && this.registration.email.name) {
                  params.push(this.registration.email.name)
              }
              if (this.registration && this.registration.phone && this.registration.phone.name) {
                  params.push(this.registration.phone.name)
              }
              return params
          }
        },
        methods: {
          onClosed() {
              this.$emit('closed', false)
          },
          onRegistration() {
              const data = new FormData()
              this.getParams.foreach(item => {
                  data.append(item, this.login)
              })
          }
        },
        data () {
            return {
                login:"",
                password: {
                    text:"",
                    show: false
                },
                email:"",
                phone:"",
                checker : {
                    login: {
                        min: (this.registration && this.registration.login && this.registration.login.min) || 3,
                        max: (this.registration && this.registration.login && this.registration.login.max)  || 20,
                    },
                    password: {
                        min: (this.registration && this.registration.password && this.registration.password.min)  || 8
                    },
                    email: {
                        required: (this.registration && this.registration.email && this.registration.email.required) || false
                    },
                    phone: {
                        required: (this.registration && this.registration.phone && this.registration.phone.required) || false
                    }
                },
                rules: {
                    required: value => !!value || 'поле не должно быть пустым',
                    requiredEmail: value => (!!value || !this.checker.email.required) || 'Email не может быть пустым',
                    requiredPhone: value => (!!value || !this.checker.phone.required) || 'Телефон не может быть пустым',
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
