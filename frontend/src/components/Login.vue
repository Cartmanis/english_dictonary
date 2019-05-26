<template>
  <v-dialog v-model="dialog"  max-width="400px">
    <v-card>
      <v-toolbar height="60px">
        <v-toolbar-title>Авторизация</v-toolbar-title>
      </v-toolbar>
      <v-card-text>
        <v-layout wrap>
          <v-flex xs12>
            <v-text-field label="Имя пользователя" required></v-text-field>
            <v-text-field label="Пароль" required></v-text-field>
          </v-flex>
          <v-flex xs6>
            <v-spacer></v-spacer>
            <v-btn color = "primary">Вход</v-btn>
          </v-flex>
          <v-flex xs6>
            <v-spacer></v-spacer>
            <v-btn flat color = "primary">Регистрация</v-btn>
          </v-flex>
        </v-layout>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
    export default {
        name: "Login",
        data: () => ({
            dialog: true
        })
    }
</script>

<style scoped>

</style>