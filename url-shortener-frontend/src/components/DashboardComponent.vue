<template>
    <div class="dashboard-container">
        <h1 class="dashboard-title">
            Link Management Dashboard
            <i class="fas fa-link" aria-hidden="true"></i>
            <i class="fas fa-chart-bar" aria-hidden="true" style="margin-left:5px;"></i>
        </h1>
        <hr/>
        <div class="scrollable-box">
            <!-- Loading Spinner -->
            <div v-if="loading" class="spinner"></div>

            <!-- Table for displaying URLs -->
            <table v-else class="urls-table">
                <thead>
                    <tr style="border-bottom: 2px solid gray;">
                        <th style="font-weight: bold; color: black; text-align: center;">Short URL</th>
                        <th style="font-weight: bold; color: black; text-align: center;">Original URL</th>
                        <th style="font-weight: bold; color: black; text-align: center;">Custom Alias</th>
                        <th style="font-weight: bold; color: black; text-align: center;">Expiration Date</th>
                        <th style="font-weight: bold; color: black; text-align: center;">Total Clicks</th>
                        <th style="font-weight: bold; color: black; text-align: center;">Edit</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="url in urls" :key="url.short_url">
                        <td style="font-weight: bold; color: blue;">
                            <a 
                                :href="`${getGoserverUrl}/${url.short_url}`" 
                                style="text-decoration: underline;"
                                target="_blank"
                            >
                                    {{ url.short_url }}
                            </a>
                        </td>
                        <td style="color: black;">{{ url.original_url }}</td>
                        <td style="color: black;">{{ url.custom_alias || 'N/A' }}</td>
                        <td style="color: black;">{{ formatDate(url.expiration_date) }}</td>
                        <td style="font-weight: bold; color: black;">
                            <a 
                                :href="`${getVueserverUrl}/analytics/${url.short_url}`" 
                                target="_blank" 
                                style="text-decoration: underline; color: black;" 
                            >
                                {{ url.total_clicks }}
                            </a>
                        </td>
                        <td @click="deleteUrl(url.short_url)" style="cursor: pointer;">
                            <i class="fas fa-trash" aria-hidden="true" title="Delete"></i> <!-- Delete icon -->
                        </td>
                    </tr>
                </tbody>
            </table>

        </div>
    </div>
</template>

<script>
    import Swal from 'sweetalert2';
    export default {
        name: 'DashboardComponent',
        data() {
            return {
                loading : true,
                urls: [],   // Array to hold data
            }
        },
        computed: {
            getVueserverUrl() {
                return process.env.VUE_APP_SERVER_URL;
            },
            getGoserverUrl() {
                // Returns Go server URL base path
                return process.env.VUE_APP_GO_SERVER_URL;
            }
        },
        methods: {
                async fetchUrls() {
                    try {
                        const response = await fetch(`${process.env.VUE_APP_GO_SERVER_URL}/urls`);
                        this.urls = await response.json();
                    } catch (error) {
                        console.error('Error fetching URLs:', error);
                    } finally {
                        this.loading = false;
                    }
                },
                formatDate(dateString) {
                    const options = { year: 'numeric', month: '2-digit', day: '2-digit' };
                    const date = new Date(dateString);
                    return date.toLocaleDateString('en-GB', options); // Format as DD/MM/YYYY
                },

                deleteUrl(shortUrl) {
                    if (confirm('Are you sure you want to delete this URL?')) {
                        // Call the API to delete the URL
                        fetch(`${process.env.VUE_APP_GO_SERVER_URL}/delete?short_url=${shortUrl}`, { method: 'DELETE' })
                            .then(async(response) => {
                                // Refresh the urls since shortUrl is deleteted
                                this.urls = this.urls.filter(url => url.short_url !== shortUrl);
                
                                // Handle error response(since deletion) and display message
                                return response.text().then(() => {
                                    Swal.fire({
                                        icon: 'success',
                                        title: 'Deleted!',
                                        text: 'URL is deleted successfully.',
                                        confirmButtonText: 'OK'
                                    });
                                });
                            })
                            .catch(error => {
                                console.error('Error deleting URL:', error);
                                Swal.fire({
                                    icon: 'error',
                                    title: 'Error!',
                                    text: 'Failed to delete URL due to network error.',
                                    confirmButtonText: 'OK'
                                });
                            });
                        }
                    }

        },
        mounted() {
            this.fetchUrls(); // Fetch URLs when the component is mounted
        }
    };
</script>