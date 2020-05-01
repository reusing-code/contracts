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

    <v-content>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-content>

    <v-footer app>
      <!-- -->
    </v-footer>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      drawer: null,
      navItems: [
        { title: "Home", icon: "mdi-home", link: "/" },
        { title: "About", icon: "mdi-corn", link: "/about" },
        { title: "Contracts", icon: "mdi-file-cabinet", link: "/contracts" },
        { title: "Test", icon: "mdi-square-edit-outline", link: "/test" },
      ],
    };
  },
};
</script>

<style></style>
