<template>
    <form @submit.prevent="shortenURL" style="text-align: center;">
        <input v-model="longURL" placeholder="Enter your long URL" ref="urlInput" />
        <input v-model="customAlias" placeholder="Custom Alias (optional)" ref="aliasInput" />
        <input type="date" v-model="expirationDate" ref="expirationInput" :min="today" required />

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

                <!-- New Analytics Button -->
                <button @click.prevent="viewAnalytics(shortenedURL)" class="analytics-button">
                    📊
                </button>
                <CopyTooltip v-if="tooltipVisible" :text="'Copied!'" :position="tooltipPosition" ref="tooltip"/>
            </div>
        </div>

    </form>
</template>

<script>
    import { shortenURL } from '../api';
    import Swal from 'sweetalert2';
    import CopyTooltip from './CopyTooltip.vue';

    export default {
    components: { CopyTooltip },

    data() {
        return {
            longURL: '',
            shortenedURL: '',
            customAlias: '',
            expirationDate: '',
            tooltipVisible: false, 
            tooltipPosition: { top: '0px', left: '0px' }, // Initial position
            loading: false, // Add loading state for spinner
        };
    },

    computed: {
        today() {
            const date = new Date();
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are zero-indexed
            const day = String(date.getDate()).padStart(2, '0');
            return `${year}-${month}-${day}`; // Format as YYYY-MM-DD
        }
    },

    methods: {
        async shortenURL() {
            this.loading = true;

            // Called when the form is submitted
            console.log('Submitted URL:', this.longURL); // Log the submitted URL

            if (!this.longURL || !this.expirationDate) {
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
                const data = {
                    longURL: this.longURL,
                    alias: this.customAlias,
                    expiration: this.expirationDate // Send expiration date
            };
                this.shortenedURL = await shortenURL(data);
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
        },

        async viewAnalytics(shortenedURL) {
            const shortUrlSegment = shortenedURL.split('/').pop(); // Extract short URL segment
            const analyticsUrl = `${process.env.VUE_APP_SERVER_URL}/analytics/${shortUrlSegment}`; // Construct the analytics URL
            window.open(analyticsUrl, '_blank'); // Open in a new tab/window
        },
    },
};
</script>