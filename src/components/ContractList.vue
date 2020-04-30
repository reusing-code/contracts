<template>
  <div id="contract-list">
    <table>
      <thead>
        <tr>
          <th v-for="column in columns" :key="column">
            {{ column }}
          </th>
        </tr>
      </thead>
      <tbody>
        <TableRow
          v-for="contract in contracts"
          :key="contract.id"
          v-bind:data="contract"
          v-bind:columns="columns"
        ></TableRow>
      </tbody>
    </table>
  </div>
</template>

<script>
import TableRow from "@/components/TableRow.vue";

export default {
  components: { TableRow },
  data() {
    return {
      columns: ["id", "name", "value"],
      contracts: [
        { id: 1, name: "bank", value: "abc" },
        { id: 2, name: "insurenace", value: "def" },
        { id: 3, name: "mobile", value: "ghi" },
        { id: 4, name: "car", value: "jkl" }
      ]
    };
  }
};
</script>

<style scoped>
table {
  margin: auto;
}
</style>
