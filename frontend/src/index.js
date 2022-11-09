import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

// Default for production, use the same server as the frontend
let API_URL = "/query";

if (process.env.REACT_APP_API_URL) {
    API_URL = process.env.REACT_APP_API_URL + "/query";
} else if (process.env.NODE_ENV !== 'production') {
    // For dev, maybe add another condition for tests
    API_URL = "http://localhost:8080/query";
}

console.log('API_URL: ' + API_URL);

const client = new ApolloClient({
  uri: API_URL,
  cache: new InMemoryCache(),
});

// Supported in React 18+
const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
    <ApolloProvider client={client}>
      <App />
    </ApolloProvider>,
);
