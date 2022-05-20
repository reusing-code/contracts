<template>
  <v-chip :color="color" text-color="white">
    <slot></slot>
  </v-chip>
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      default: "primary",
    },
  },
};
</script>

<style scoped></style>
