// src/api.js
import axios from 'axios';

// Vue environment variables should be prefixed with VUE_APP
const BASE_URL = process.env.VUE_APP_GO_SERVER_URL;

// Function to validate the URL
const isValidURL = (urlString) => {
  try {
    new URL(urlString); // This will throw an error if the URL is invalid
    return true; // If no error is thrown, the URL is valid
  } catch (error) {
    return false; // If an error is thrown, the URL is invalid
  }
};

// Function to validate the URL by making a request to the backend
const validateURL = async (longURL) => {
    try {
        const response = await axios.post(`${BASE_URL}/validate`, new URLSearchParams({ url: longURL }));
        return response.status === 200; // URL is valid if status is 200
    } catch (error) {
        return false; // URL is invalid if any error occurs
    }
};

export const shortenURL = async (data) => {
    const { longURL, alias, expiration } = data; // Destructure data object

    // First, validate the URL before making the request
    if (!isValidURL(longURL)) {
        throw new Error('Invalid URL format. Please provide a valid URL.');
    }

    // Then, validate the URL by making a request
    const isValid = await validateURL(longURL);

    if (!isValid) {
        throw new Error('Invalid URL. Please provide a reachable URL.');
    }

    try {
        const response = await axios.post(`${BASE_URL}/shorten`, 
            new URLSearchParams({ 
                url: longURL,
                alias: alias || '', // Optional custom alias
                expiration: expiration 
        }));

        return response.data;
    } catch (error) {
        throw new Error('Failed to shorten URL. Please check the console for details.');
    }
};
