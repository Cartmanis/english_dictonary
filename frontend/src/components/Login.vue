<template>
  <v-container fluid fill-height>
    <v-layout align-center justify-center>
      <v-flex xs12 sm4 md3>
        <v-card class="elevation-12">
          <v-toolbar>
            <v-toolbar-title>Авторизация</v-toolbar-title>
            <v-spacer></v-spacer>
          </v-toolbar>
          <v-card-text>
            <v-form>
              <v-text-field  prepend-icon="person" label="Имя пользователя"></v-text-field>
              <v-text-field
                      v-model='password.text' prepend-icon="lock" label="Пароль"
                      :type= "password.show? 'text' :'password'"
                      :append-icon="password.show ? 'visibility' : 'visibility_off'"
                      @click:append="password.show=!password.show">
              </v-text-field>
              <h4 v-if="error" align="center" style="color: red;">имя пользователя или пароль введены не верно</h4>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary">Вход</v-btn>
            <v-btn color="primary">Закрыть</v-btn>
            <v-spacer></v-spacer>
            <v-btn flat color="info">Регистрация</v-btn>
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
          error: false,
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
      props: {
          hidePassword: {
            type: Boolean
          }
      },
      methods: {
        onChange () {
          this.$emit('change-data', this.settings)
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
