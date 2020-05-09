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
          <v-dialog v-model="dialog" max-width="1000px">
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
                    <v-col
                      cols="12"
                      md="6"
                      v-for="field in fields.filter(x => x.type !== '')"
                      :key="field.value"
                    >
                      <base-input
                        v-model="editedItem[field.value]"
                        :label="$t('contracts.' + field.value)"
                        :type="field.type"
                      ></base-input>
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
      fields: [
        { value: "id", type: "" },
        { value: "company", type: "string" },
        { value: "product", type: "string" },
        { value: "start", type: "date" },
        { value: "end", type: "date" },
        { value: "extensionMonths", type: "number" },
        { value: "noticePeriodMonths", type: "number" },
        { value: "pricePerMonth", type: "currency" },
        { value: "contractNumber", type: "string" },
        { value: "customerNumber", type: "string" },
        { value: "account", type: "string" },
        { value: "paymentOption", type: "string" },
        { value: "notes", type: "text" },
        { value: "category", type: "" },
        { value: "status", type: "" }
      ],
      contracts: [
        { id: 1, name: "bank", value: "abc" },
        { id: 2, name: "insurance", value: "def" },
        { id: 3, name: "mobile", value: "ghi" },
        { id: 4, name: "car", value: "jkl" }
      ],
      dialog: false,
      editedIndex: -1,
      editedItem: {
        id: -1,
        name: "",
        value: ""
      },
      defaultItem: {
        id: -1,
        name: "",
        value: ""
      }
    };
  },
  mounted() {
    this.loadPersisted();
  },
  computed: {
    formTitle() {
      return this.editedIndex === -1 ? this.$t("newitem") : this.$t("edititem");
    },
    headers() {
      var res = this.fields
        .filter(x => x.type !== "")
        .map(x => ({
          value: x.value,
          text: this.$t("contracts." + x.value)
        }));
      res.push({ value: "actions", text: this.$t("actions") });
      return res;
    }
  },
  watch: {
    dialog(val) {
      val || this.close();
    }
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
      this.persistLocally();
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
        const maxID = Math.max(...this.contracts.map(c => c.id), 0);
        this.editedItem.id = maxID + 1;
        this.contracts.push(this.editedItem);
      }
      this.persistLocally();
      this.close();
    },
    persistLocally() {
      const parsed = JSON.stringify(this.contracts);
      localStorage.setItem("contracts", parsed);
    },
    loadPersisted() {
      if (localStorage.getItem("contracts")) {
        try {
          this.contracts = JSON.parse(localStorage.getItem("contracts"));
        } catch (e) {
          localStorage.removeItem("contracts");
        }
      }
    }
  }
};
</script>

<style scoped></style>
