<template>
  <div class="locale-changer">
    <v-select v-model="currentLang" :items="langs" prepend-icon="mdi-translate">
    </v-select>
  </div>
</template>

<script>
export default {
  name: "locale-changer",
  data() {
    return {
      langs: ["en", "de"],
      currentLang: this.$i18n.locale,
    };
  },
  watch: {
    currentLang(lang) {
      this.$i18n.locale = lang;
      this.$vuetify.lang.current = lang;
    },
  },
};
</script>
