<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
import { CLIENT_ID, OAUTH_CALLBACK_URL } from "@/constants";
export default {
  name: "App",
  components: {},
  data() {
    return {};
  },
  async mounted() {
    if (this.$route.path !== "/wxpay") {
      this.GLOBAL.api.net.on(401, this.redirectToOAuth);
      this.GLOBAL.api.net.on(403, this.redirectToOAuth);
      this.GLOBAL.api.net.on(460, payload => {
        this.$router.push("/invitation/entry");
      });
    }
  },
  methods: {
    redirectToOAuth(payload) {
      let url = `https://mixin.one/oauth/authorize?client_id=${CLIENT_ID}&scope=PROFILE:READ+ASSETS:READ+MESSAGES:REPRESENT&response_type=code&return_to=${encodeURIComponent(
        OAUTH_CALLBACK_URL
      )}`;
      window.location.href = url;
    }
  }
};
</script>

<style>
html,
body {
  padding: 0;
  margin: 0;
  height: 100%;
  font-size: 14px;
  font-family: "Roboto", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#app {
  background: #f8f8f8;
  padding: 0;
  height: 100%;
}

.van-cell.van-field {
  padding: 0;
}
</style>
