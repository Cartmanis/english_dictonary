<template>
  <v-app>
    <v-content>
      <router-view @click-login="onLogin"/>
    </v-content>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data () {
    return {

    }
  },
  methods: {
    onLogin(isAuth) {
        if (isAuth) {
            this.$router.push('/')
        }
    }
  }
}
</script>
