<template>
  <v-card>
    <v-data-table :headers="headers" :items="contracts" :search="search">
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>Contracts</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            :label="$t('search')"
            single-line
            hide-details
          ></v-text-field>
          <v-spacer></v-spacer>
          <v-dialog v-model="dialog" max-width="500px">
            <template v-slot:activator="{ on }">
              <v-btn color="primary" dark class="mb-2" v-on="on">{{
                $t("newitem")
              }}</v-btn>
            </template>
            <v-card>
              <v-card-title>
                <span class="headline">{{ formTitle }}</span>
              </v-card-title>

              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12" sm="6" md="4">
                      <v-text-field
                        v-model="editedItem.id"
                        label="ID"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-text-field
                        v-model="editedItem.name"
                        label="Name"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-text-field
                        v-model="editedItem.value"
                        label="Value"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" text @click="close">{{
                  $t("cancel")
                }}</v-btn>
                <v-btn color="primary" text @click="save">{{
                  $t("save")
                }}</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-toolbar>
      </template>
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="editItem(item)">
          mdi-pencil
        </v-icon>
        <v-icon small @click="deleteItem(item)">
          mdi-delete
        </v-icon>
      </template>
    </v-data-table>
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
        { text: "Value", value: "value" },
        { text: "Actions", value: "actions", sortable: false },
      ],
      contracts: [
        { id: 1, name: "bank", value: "abc" },
        { id: 2, name: "insurance", value: "def" },
        { id: 3, name: "mobile", value: "ghi" },
        { id: 4, name: "car", value: "jkl" },
      ],
      dialog: false,
      editedIndex: -1,
      editedItem: {
        id: -1,
        name: "",
        value: "",
      },
      defaultItem: {
        id: -1,
        name: "",
        value: "",
      },
    };
  },
  computed: {
    formTitle() {
      return this.editedIndex === -1 ? "New Item" : "Edit Item";
    },
  },
  watch: {
    dialog(val) {
      val || this.close();
    },
  },

  methods: {
    editItem(item) {
      this.editedIndex = this.contracts.indexOf(item);
      this.editedItem = Object.assign({}, item);
      this.dialog = true;
    },
    deleteItem(item) {
      const index = this.contracts.indexOf(item);
      confirm("Are you sure you want to delete this item?") &&
        this.contracts.splice(index, 1);
    },
    close() {
      this.dialog = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    },
    save() {
      if (this.editedIndex > -1) {
        Object.assign(this.contracts[this.editedIndex], this.editedItem);
      } else {
        this.contracts.push(this.editedItem);
      }
      this.close();
    },
  },
};
</script>

<style scoped></style>

<i18n>
{
  "en": {
    "search": "Search",
    "newitem": "New item",
    "cancel": "Cancel",
    "save": "Save"
  },
  "de": {
    "search": "Suche",
    "newitem": "Neu",
    "cancel": "Abbruch",
    "save": "Speichern"
  }
}
</i18n>
