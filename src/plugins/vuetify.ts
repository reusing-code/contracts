import Vue from "vue";
import Vuetify from "vuetify/lib";

Vue.use(Vuetify);

// translations
import en from "vuetify/es5/locale/en";
import de from "vuetify/es5/locale/de";

export default new Vuetify({
  lang: {
    locales: { en, de },
    current: "en",
  },
});
