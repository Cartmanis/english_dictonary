<template>
  <div>
    <v-snackbar
        v-model="snackbar.show"
        :color="snackbar.color"
        :timeout="10000"
    >{{snackbar.text}}</v-snackbar>
    <v-card
        color="grey lighten-4"
        flat
        height="200px"
        tile
    >
      <v-toolbar dense>
        <v-toolbar-side-icon></v-toolbar-side-icon>
        <v-toolbar-title>Главное меню</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-label>Пользователь: {{userName}}</v-label>
        <v-btn flat @click="onLogOut">Выход</v-btn>
        <v-btn icon>
          <v-icon>more_vert</v-icon>
        </v-btn>
      </v-toolbar>
      <words></words>
    </v-card>
  </div>
</template>

<script>
  const axios = require('axios').default
  import Words from './Words'
  export default {
    components: {
      Words
    },
    data: () => ({
      snackbar: {
        show: false,
        text: "",
        color: ""
      },
      userName:"",
      items: [
        { title: 'Click Me' },
        { title: 'Click Me' },
        { title: 'Click Me' },
        { title: 'Click Me 2' }
      ]
    }),
    created: function () {
      this.getUserName()
    },
    methods: {
      async getUserName() {
        try {
          const url = "http://192.168.0.83:27333/api/v1/auth"
          const res = await axios.post(url, {}, {withCredentials:true})
          if (res && res.data && res.data.user) {
            this.userName = res.data.user
          }
        } catch(e) {
          this.showSnackBar(`не удалось получить имя пользователя. ${e}`)
        }
      },
      async onLogOut() {
        const url = "http://192.168.0.83:27333/api/v1/logout"
        try {
          const res = await  axios.post(url, {}, {withCredentials:true})
          if (res && res.data && res.data.result) {
            this.$router.push('/login')
          }
        } catch(e) {
          this.showSnackBar(`не удалось произвести выход пользователя ${this.userName}. ${e}`)
        }
      },
      showSnackBar(text, color) {
        this.snackbar.show = true
        this.snackbar.text = text
        this.snackbar.color = color || "error"
      },
    }
  }
</script>

<style>

</style>
