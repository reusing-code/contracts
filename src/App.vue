<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-list-item
          link
          v-for="item in navItems"
          :key="item.title"
          :to="item.link"
        >
          <v-list-item-action>
            <v-icon v-if="item.icon">{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>Contracts</v-toolbar-title>
    </v-app-bar>

    <v-main>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-main>

    <v-footer app>
      <v-spacer></v-spacer>
      <v-col cols="1">
        <LocalChanger></LocalChanger>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script>
import LocalChanger from "@/components/LocaleChanger.vue";

export default {
  components: { LocalChanger },
  data() {
    return {
      drawer: null,
      navItems: [
        { title: "Home", icon: "mdi-home", link: "/" },
        { title: "About", icon: "mdi-corn", link: "/about" },
        { title: "Contracts", icon: "mdi-file-cabinet", link: "/contracts" }
      ]
    };
  }
};
</script>

<style></style>
