<template>
  <div class="locale-changer">
    <v-select
      v-model="$i18n.locale"
      :items="langs"
      prepend-icon="mdi-translate"
    >
    </v-select>
  </div>
</template>

<script>
export default {
  name: "locale-changer",
  data() {
    return { langs: ["en", "de"] };
  },
};
</script>
