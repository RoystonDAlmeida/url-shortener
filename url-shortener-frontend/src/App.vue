<template>
  <div id="app">
    <div class="header">
      <img src="@/assets/url_shortener_logo.jpeg" alt="URL Shortener Logo" class="logo" />
      <h1>URL Shortener</h1>
      <router-link to="/dashboard">
        <button class="dashboard-button">Dashboard</button>
      </router-link>
    </div>
    <router-view></router-view> <!-- Render components here depending on their route -->
  </div>
</template>

<script>
  export default {
    name: 'App',
  };
</script>

<style>
  @import './styles/styles.css';
</style>