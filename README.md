# URL Shortener Application

This project is a URL shortener application consisting of a backend service built with Go and a frontend built with Vue.js. The backend handles URL shortening, storage, and analytics, while the frontend provides a user-friendly interface for users to shorten URLs and view analytics.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Requirements](#requirements)
- [Backend Setup](#backend-setup)
  - [Environment Variables](#backend-environment-variables)
  - [Installing Air for Hot Reloading](#installing-air-for-hot-reloading)
  - [Running the Backend](#running-the-backend)
- [Frontend Setup](#frontend-setup)
  - [Environment Variables](#frontend-environment-variables)
  - [Running the Frontend](#running-the-frontend)
- [Usage](#usage)
- [Analytics Dashboard](#analytics-dashboard)
- [Contributing](#contributing)
- [License](#license)

## Features

- Shorten long URLs into shorter, more manageable links.
- Retrieve original URLs from shortened links.
- Custom alias option for personalized shortened URLs.
- Track click counts for each shortened URL over time.
- Dashboard displaying detailed information about each URL, including:
  - Short URL
  - Original URL
  - Custom alias (if provided)
  - Click analytics
- User-friendly web interface for easy interaction.
- Hot reloading during development for a better workflow.

## Architecture

The application consists of two main components:

1. **Backend**: Built using Go, it manages the URL shortening logic, interacts with the SQLite database, and tracks click analytics.
2. **Frontend**: Built using Vue.js, it provides a graphical interface for users to input URLs, view shortened links, and access analytics.

## Requirements

Before running the application, ensure you have the following installed:

- Go (version 1.16 or higher) for the backend
- Node.js (version 14 or higher) and npm for the frontend
- SQLite3 database for storing URLs and analytics data

## Backend Setup

### Environment Variables

Create a `.env` file in the `url-shortener-backend` directory with the following variables:

```bash
DATABASE_URL = <your_relative_database_url>
SERVER_ADDRESS = "localhost:8080"
CHARSET = <your_sequence_of characters>
```

- **DATABASE_URL**: Connection string for your SQLite database.
- **SERVER_ADDRESS**: Address where the backend server will run.
- **CHARSET**: Characters used for generating shortened URLs.

### Installing Air for Hot Reloading

To enable hot reloading in your Go application, you need to install Air. Follow these steps:

1. Install Air using Go:
```bash
go install github.com/cosmtrek/air@latest
```

2. Navigate to your project’s root directory and initialize Air:
```bash
cd url-shortener-backend
air init
```

This will create a `.air.toml` configuration file with default settings.

### Running the Backend

1. Navigate to the backend directory:
```bash
cd url-shortener-backend
```

2. Start the server with hot reloading using Air:
```bash
air
```

3. The backend should now be running on the address specified in `SERVER_ADDRESS`.

## Frontend Setup

### Environment Variables

Create a `.env` file in the `url-shortener-frontend` directory with the following variables:
```bash
VUE_APP_SERVER_URL = "http://localhost:8081"
VUE_APP_GO_SERVER_URL = "http://localhost:8080"
```

- **VUE_APP_SERVER_URL**: The URL of your backend service.
- **VUE_APP_GO_SERVER_URL**: Another variable pointing to your backend service (can be used for different purposes in your app).

### Running the Frontend

1. Navigate to the frontend directory:
```bash
cd url-shortener-frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run serve
```

4. The frontend should now be running at `http://localhost:8081` or another port specified by Vue CLI.

## Usage

1. Open your web browser and navigate to `http://localhost:8081` (or wherever your frontend is running).
2. Enter a long URL in the input field and click "Shorten" to generate a shortened link.
3. Optionally, provide a custom alias for your shortened link.
4. Copy and use the shortened link as needed.

## Analytics Dashboard

The application includes an analytics dashboard that allows users to view detailed information about their shortened URLs:

1. **Short URL**: Displays the generated short link.
2. **Original URL**: Shows the original long link that was shortened.
3. **Custom Alias**: If provided, displays the custom alias used for shortening.
4. **Click Analytics**: Tracks and displays click counts at different times, allowing users to monitor how many times their links have been accessed.

To access the dashboard, navigate to `/dashboard` after shortening URLs.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](https://mit-license.org/) file for details.