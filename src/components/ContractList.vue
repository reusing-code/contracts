<template>
  <v-card>
    <v-card-title>
      Contracts
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="contracts"
      :search="search"
    ></v-data-table>
  </v-card>
</template>

<script>
export default {
  data() {
    return {
      search: "",
      headers: [
        { text: "ID", value: "id" },
        { text: "Name", value: "name" },
        { text: "Value", value: "value" }
      ],
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

<style scoped></style>
