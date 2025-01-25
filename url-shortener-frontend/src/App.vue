<template>
  <div id="app">
    <div class="header">
      <img src="@/assets/url_shortener_logo.jpeg" alt="URL Shortener Logo" class="logo" />
      <h1>URL Shortener</h1>
    </div>
    <form @submit.prevent="shortenURL">
      <input v-model="longURL" placeholder="Enter your long URL" />
      <button type="submit">Shorten</button>
    </form>
    <div v-if="shortenedURL">
      <p>Shortened URL: <a :href="shortenedURL" target="_blank">{{ shortenedURL }}</a></p>
    </div>
  </div>
</template>

<script>
import { shortenURL } from './api';
import Swal from 'sweetalert2';

export default {
  data() {
    return {
      longURL: '',
      shortenedURL: ''
    };
  },
  methods: {
    async shortenURL() {
      // Called when the form is submitted
      console.log('Submitted URL:', this.longURL); // Log the submitted URL

      if (!this.longURL) {
      // Check if URL is not provided
      Swal.fire({
        icon: 'error', // Change icon to error for better visibility
        title: 'Input Required',
        text: 'Please enter a URL to shorten.',
        confirmButtonText: 'OK',
        position: 'top',
        allowOutsideClick: false, // Prevent closing by clicking outside
        timer: null,  // Disable auto close timer
        showConfirmButton: true // Show confirm button
      });
      return;
    }

    try {
      this.shortenedURL = await shortenURL(this.longURL);
    } catch (error) {
      console.error('Error shortening URL:', error);
      Swal.fire({
        icon: 'error',
        title: 'Error',
        text: 'Failed to shorten URL. Please check the console for details.',
        confirmButtonText: 'OK',
        allowOutsideClick: false,
        timer: null,
        showConfirmButton: true
      });
    }
  }
}
};
</script>

<style>
/* Import the external CSS file */
@import './styles/styles.css';

</style>