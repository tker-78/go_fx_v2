<template>
  <v-app>
    <v-app-bar scroll-behavior="collapse" color="blue">
      <v-app-bar-title>Go FX V2</v-app-bar-title>
      <template v-slot:append>
        <v-btn icon="mdi-heart"></v-btn>
    
        <v-btn icon="mdi-magnify"></v-btn>
    
        <v-btn icon="mdi-dots-vertical"></v-btn>
      </template>

      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link>
    </v-app-bar>
    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

nav {
  padding: 30px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #42b983;
}
</style>
