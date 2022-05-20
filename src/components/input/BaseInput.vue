<template>
  <div v-if="type == 'string'">
    <v-text-field
      :value="value"
      :label="label"
      @input="updateValue"
      clearable
    ></v-text-field>
  </div>
  <div v-else-if="type == 'number'">
    <v-text-field
      :value="value"
      :label="label"
      type="number"
      @input="updateValue"
      clearable
    ></v-text-field>
  </div>
  <div v-else-if="type == 'text'">
    <v-textarea
      :value="value"
      :label="label"
      auto-grow
      rows="3"
      @input="updateValue"
      clearable
    ></v-textarea>
  </div>
  <div v-else-if="type == 'select'">
    <v-select
      :value="value"
      :label="label"
      :items="options.items"
      @input="updateValue"
    ></v-select>
  </div>
  <div v-else-if="type == 'currency'">
    <v-text-field
      :value="value"
      :label="label"
      type="number"
      suffix="€"
      @input="updateValue"
      clearable
    ></v-text-field>
  </div>
  <div v-else-if="type == 'date'">
    <v-menu
      v-model="menu"
      :close-on-content-click="false"
      transition="scale-transition"
      offset-y
      min-width="290px"
    >
      <template v-slot:activator="{ on }">
        <v-text-field
          :value="value"
          :label="label"
          prepend-icon="mdi-calendar"
          readonly
          v-on="on"
          @input="updateValue"
          clearable
        ></v-text-field>
      </template>
      <v-date-picker
        ref="picker"
        :value="value"
        @input="updateValue"
        no-title
        scrollable
      >
      </v-date-picker>
    </v-menu>
  </div>
</template>

<script>
export default {
  props: ["value", "type", "label", "options"],
  data() {
    return {
      menu: false,
    };
  },
  watch: {
    menu(val) {
      val && setTimeout(() => (this.$refs.picker.activePicker = "YEAR"));
    },
  },
  methods: {
    updateValue(val) {
      this.menu = false;
      this.$emit("input", val);
    },
  },
};
</script>

<style scoped></style>
