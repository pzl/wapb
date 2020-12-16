<template>
  <v-app dark>
    <v-app-bar app flat>
      <v-tabs centered class="ml-n9">
        <v-tab v-for="link in links" :key="link" nuxt :to="'/'+link">{{ link }}</v-tab>
      </v-tabs>
    </v-app-bar>
    <v-main>
      <v-container>
        <nuxt />
      </v-container>
    </v-main>
    <v-footer app>
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      links: ['text','file','link']
    }
  },
  methods: {
    setTheme() {
      if (window && window.matchMedia && window.matchMedia('(prefers-color-scheme:dark)').matches) {
        this.$vuetify.theme.dark = true;
      } else {
        this.$vuetify.theme.dark = false;
      } 
    }
  },
  mounted() {
    this.setTheme();
    if (window && window.matchMedia) {
      window.matchMedia('(prefers-color-scheme:dark)').addListener(e => {
        if (e.matches) {
          this.$vuetify.theme.dark = true;
        } else {
          this.$vuetify.theme.dark = false;
        }
      });
    }
  }
}
</script>
