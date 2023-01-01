import {ApolloClient, InMemoryCache} from "@apollo/client";

export let API_URL = "/query";

if (process.env.REACT_APP_API_URL) {
    API_URL = process.env.REACT_APP_API_URL + "/query";
} else if (process.env.NODE_ENV !== 'production') {
    // For dev, maybe add another condition for tests
    API_URL = "http://localhost:8080/query";

}

console.log('API_URL: ' + API_URL);

export const client = new ApolloClient({
    uri: API_URL,
    cache: new InMemoryCache(),
});
