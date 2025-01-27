<template>
  <div id="app">
    <div class="header">
      <img src="@/assets/url_shortener_logo.jpeg" alt="URL Shortener Logo" class="logo" />
      <h1>URL Shortener</h1>
    </div>
    <form @submit.prevent="shortenURL" style="text-align: center;">
      <input v-model="longURL" placeholder="Enter your long URL" ref="urlInput" />
      <button type="submit">Shorten</button>

      <!-- Show spinner while loading -->
      <div v-if="loading" class="spinner"></div>

      <!-- Shortened URL Display -->
      <div v-if="shortenedURL" class="shortened-url-container">
        <div class="shortened-url">
          <p v-if="shortenedURL">Shortened URL: 
            <a :href="shortenedURL" target="_blank">{{ shortenedURL }}</a>
          </p>
          <button @click="copyToClipboard(shortenedURL)" class="copy-button" ref="copyButton">
            📋 
          </button>
          <CopyTooltip v-if="tooltipVisible" :text="'Copied!'" :position="tooltipPosition" ref="tooltip"/>
        </div>
      </div>

    </form>
  </div>
</template>

<script>
import { shortenURL } from './api';
import Swal from 'sweetalert2';
import CopyTooltip from './components/CopyTooltip.vue';

export default {
  components: { CopyTooltip },

  data() {
    return {
      longURL: '',
      shortenedURL: '',
      tooltipVisible: false, 
      tooltipPosition: { top: '0px', left: '0px' }, // Initial position
      loading: false, // Add loading state for spinner
    };
  },
  methods: {
    async shortenURL() {
      this.loading = true;

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
        text: error.message,
        confirmButtonText: 'OK',
        allowOutsideClick: false,
        timer: null,
        showConfirmButton: true
      });
      this.$refs.urlInput.value = ''; // Clear the input field
    }
    finally {
      this.loading = false;
    }
  },

  copyToClipboard() {
    const buttonElement = this.$refs.copyButton; // Reference to the button
    if (buttonElement) {
      navigator.clipboard.writeText(this.shortenedURL)
        .then(() => {
          const rect = buttonElement.getBoundingClientRect();
          this.tooltipPosition = { top: `${rect.top - 30}px`, left: `${rect.left}px` }; // Adjust position
          this.tooltipVisible = true; // Show tooltip
          
          // Use $nextTick to ensure tooltip is rendered before accessing it
          this.$nextTick(() => {
            if (this.$refs.tooltip) {
              this.$refs.tooltip.show(); // Trigger show method in CopyTooltip component
            } else {
              console.error('Tooltip reference is still undefined');
            }
          });
        })
        .catch(err => {
          console.error('Failed to copy:', err);
        });
    } else {
      console.error('Button element is undefined');
    }
  }

  }
};
</script>

<style>
  /* Import the external CSS file */
  @import './styles/styles.css';
</style>