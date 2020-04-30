<template>
  <tr>
    <td v-for="column in columns" :key="column">{{ data[column] }}</td>
  </tr>
</template>

<script>
export default {
  props: {
    columns: {
      type: Array
    },
    data: {
      type: Object,
      default: function() {
        return {};
      }
    }
  }
};
</script>

<style scoped></style>
