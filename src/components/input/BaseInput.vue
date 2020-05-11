<template>
  <div v-if="type == 'string'">
    <v-text-field
      v-model="internalValue"
      :label="label"
      clearable
    ></v-text-field>
  </div>
  <div v-else-if="type == 'number'">
    <v-text-field
      v-model="internalValue"
      :label="label"
      type="number"
      clearable
    ></v-text-field>
  </div>
  <div v-else-if="type == 'text'">
    <v-textarea
      v-model="internalValue"
      :label="label"
      auto-grow
      rows="3"
      clearable
    ></v-textarea>
  </div>
  <div v-else-if="type == 'select'">
    <v-select
      v-model="internalValue"
      :label="label"
      :items="options.items"
    ></v-select>
  </div>
  <div v-else-if="type == 'currency'">
    <v-text-field
      v-model="internalValue"
      :label="label"
      type="number"
      suffix="€"
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
          v-model="internalValue"
          :label="label"
          prepend-icon="mdi-calendar"
          readonly
          v-on="on"
          clearable
        ></v-text-field>
      </template>
      <v-date-picker
        ref="picker"
        v-model="internalValue"
        @input="menu = false"
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
      internalValue: this.value,
      menu: false,
    };
  },
  watch: {
    internalValue(val) {
      this.$emit("input", val);
    },
    menu(val) {
      val && setTimeout(() => (this.$refs.picker.activePicker = "YEAR"));
    },
  },
};
</script>

<style scoped></style>
