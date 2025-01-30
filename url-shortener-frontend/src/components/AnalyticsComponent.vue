<template>
    <div class="analytics-container">
      <h4>Analytics for Short URL: {{ shortUrl }}</h4>
  
      <!-- Loading Spinner -->
      <div v-if="loading" class="spinner"></div>
  
      <!-- Table to display URL entry -->
      <table v-if="urlEntry" class="url-entry-table">
        <thead>
          <tr>
            <th>Short URL</th>
            <th>Long URL</th>
            <th>Custom Alias</th>
            <th>Expiration Date</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>{{ urlEntry.short_url }}</td>
            <td>{{ urlEntry.long_url }}</td>
            <td>{{ urlEntry.custom_alias }}</td>
            <td>{{ urlEntry.expiration_date }}</td>
          </tr>
        </tbody>
      </table>
  
      <!-- Table Error message -->
      <div v-if="tableError" class="error-message">Loading table error: {{ tableError }}</div>
  
      <!-- Tabs for each date -->
      <div class="tabs">
        <button @click="prevDate" :disabled="currentIndex === 0">←</button>
        <span v-for="(day, index) in days" :key="index" @click="selectDate(index)" :class="{ active: currentIndex === index }">
          {{ day }}
        </span>
        <button @click="nextDate" :disabled="currentIndex === days.length - 1">→</button>
      </div>
  
      <!-- Display Line Chart -->
      <div v-if="selectedDayData">
        <h4 style="text-align:center;">{{ chartTitle }}</h4> <!-- Use dynamic chart title -->
        <apexchart type="line" :options="chartOptions" :series="chartSeries"></apexchart>
      </div>
  
      <!-- Chart Error message-->
      <div v-if="chartError" class="error-message">Loading chart error: {{ chartError }}</div>
    </div>
</template>

<script>
    import VueApexCharts from 'vue3-apexcharts';

    export default {
        components: {
            apexchart: VueApexCharts,
        },
  
    data() {
        return {
            shortUrl: this.$route.params.shortUrl,
            urlEntry: null,
            tableLoading: false,
            tableError: null,
            chartLoading: false,
            chartError: null,
            currentIndex: 0,
            days: [],
            selectedDayData: null,
            chartTitle: 'Total Click Counts', // Initialize dynamic chart title
            chartOptions: {
                chart: {
                id: 'click-counts-chart',
                toolbar: { show: false }, // Hide toolbar
                },
                xaxis: {
                title: {
                    text: 'Hourly Time',
                },
                categories: Array.from({ length: 24 }, (_, i) => `${i}:00`), // X-axis labels (0-23 hours)
                },
                yaxis: {
                title: {
                    text: 'Click Counts',
                },
                min: 0, // Set minimum value for Y-axis
                forceNiceScale: true, // Force a nice scale on Y-axis
                },
                tooltip: {
                shared: true,
                intersect: false,
                y: {
                    formatter(val) {
                    return `${val} requests`;
                    },
                    titleFormatter() {
                    return 'Time'; // Simplified tooltip title
                    }
                },
                },
                title: { // Ensure title is initialized here
                text: this.chartTitle, 
                align: 'center'
                }
            },
            chartSeries: [], // Initialize series data for the chart
        };
    },

    mounted() {
        this.getURLEntry(); // Fetch the URL entry when the component is mounted
        this.fetchAnalyticsData();
    },

    methods: {
        async getURLEntry() {
            this.tableLoading = true;
            this.tableError = null;

            try {
                const response = await fetch(`${process.env.VUE_APP_GO_SERVER_URL}/url?short_url=${this.shortUrl}`); // Adjust API endpoint as necessary
                if (!response.ok) throw new Error('Failed to fetch URL entry');
                
                this.urlEntry = await response.json();
                console.log('URL Entry:', this.urlEntry); // Log the URL entry
            } catch (err) {
                this.tableError = err.message;
                console.error('Error fetching URL entry:', err);
            } finally {
                this.tableLoading = false;
            }
        },

        async fetchAnalyticsData() {
            this.chartLoading = true;
            this.chartError = null;

            try {
                const response = await fetch(`${process.env.VUE_APP_GO_SERVER_URL}/analytics/${this.shortUrl}`); // Adjust API endpoint as necessary
                if (!response.ok) throw new Error('Failed to fetch analytics data');

                this.analyticsData = await response.json();
                console.log('Analytics Data:', this.analyticsData); // Log the analytics data
                this.days = Object.keys(this.analyticsData.day); // Extract days from the response
                this.selectDate(0); // Automatically select the first date to initialize chart
            } catch (err) {
                this.chartError = err.message;
                console.error('Error fetching analytics:', err);
            } finally {
                this.chartLoading = false;
            }
        },

        selectDate(index) {
            this.currentIndex = index;
            this.selectedDayData = this.analyticsData.day[this.days[index]]; // Update selected day data
            
            const clickCounts = this.selectedDayData.click_counts; // Get click counts for title
            const dateString = this.days[index]; // Get selected day string
      
            this.chartTitle = `Total Click Counts for ${dateString}: ${clickCounts}`;

            this.chartSeries = this.prepareChartData(this.selectedDayData); // Prepare new chart data
            console.log('Selected Day Data:', this.selectedDayData); 
            console.log('Chart Series:', this.chartSeries); 
        },

        nextDate() {
        if (this.currentIndex < this.days.length - 1) {
            this.currentIndex++;
            this.selectDate(this.currentIndex);
            }
        },

        prevDate() {
        if (this.currentIndex > 0) {
            this.currentIndex--;
            this.selectDate(this.currentIndex);
            }
        },

        prepareChartData(dayData) {
            if (!dayData || !dayData.timestamps) {
                console.error('Invalid dayData:', dayData); // Log invalid dayData
                return []; // Return empty array if dayData or its timestamps property is not defined
            }

            const hours = Array(24).fill(0); // Initialize an array for each hour of the day
            dayData.timestamps.forEach(timestamp => {
                const hour = new Date(timestamp).getHours();
                hours[hour]++; // Increment the count for the corresponding hour
            });

            return [{
                name: 'Click Counts',
                data: hours,
            }];
    },

  },
};
</script>