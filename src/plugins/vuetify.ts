import Vue from 'vue';
import Vuetify from 'vuetify/lib';

Vue.use(Vuetify);

// translations
import en from 'vuetify/src/locale/en';
import de from 'vuetify/src/locale/de';

export default new Vuetify({
  lang: {
    locales: {en, de},
    current: 'en',
  },
});
