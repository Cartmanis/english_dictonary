<template>
  <v-dialog v-model="auth"  max-width="500px">
    <v-card>
      <v-toolbar height="60px">
        <v-toolbar-title>Авторизация</v-toolbar-title>
      </v-toolbar>
      <v-card-text>
        <v-layout wrap>
          <v-flex xs12>
            <v-text-field  prepend-icon="person" label="Имя пользователя"
                           :rules="[rules.required]"
            ></v-text-field>
            <v-text-field
                v-model='password.text' prepend-icon="lock" label="Пароль"
                 :type= "password.show? 'text' :'password'"
                :append-icon="password.show ? 'visibility' : 'visibility_off'"
                @click:append="password.show=!password.show"
                :rules="[rules.required]"></v-text-field>
          </v-flex>
          <v-flex xs8>
            <v-btn @click="onLogin" color = "primary">Вход</v-btn>
          </v-flex>
          <v-flex xs4>
            <v-btn @click="onRegistration" flat color = "primary"

             >Регистрация
            </v-btn>
          </v-flex>
        </v-layout>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
    export default {
        name: "Login",
      data () {
        return {
          auth: true,
          login: "",
          password: {
            text: "",
            show:""
          },
          rules: {
            required: value => !!value || 'поле не должно быть пустым',
          }
        }
      },
      methods: {
        onLogin() {
          console.log('попали')
        },
        onRegistration() {
          console.log('регистрация')
        }
      }
    }
</script>

<style scoped>

</style>