import {ApolloClient, HttpLink, InMemoryCache} from "@apollo/client";

export let API_URL = "/query";

if (process.env.REACT_APP_API_URL) {
    API_URL = process.env.REACT_APP_API_URL + "/query";
} else if (process.env.NODE_ENV !== 'production') {
    // For dev, maybe add another condition for tests
    API_URL = "http://localhost:8080/query";

}

console.log('API_URL: ' + API_URL);

// loggingFetch logs the requests and responses made by Apollo Client. It is
// only used when the ENV variable REACT_APP_DEBUG is set.
// The function was found here:
// https://github.com/apollographql/apollo-client/issues/4017#issuecomment-1260987147
async function loggingFetch(
    input: RequestInfo | URL,
    init?: RequestInit
): Promise<Response> {
    interface Request {
        operationName: string;
        query: string;
        // ESLint is complaining about the 'any' and I don't know how to
        // improve that right now.
        // eslint-disable-next-line
        variables: Map<string, any>;
    }
    const body = JSON.parse(init?.body as string ?? "{}") as Request;

    const start = Date.now();
    console.log(
        `${new Date().toISOString().slice(-13)} 📡 Sending ${
            body.operationName
        }\nrequest: ${body.query}\nWith variables: ${JSON.stringify(
            body.variables
        )}`
    );
    const response = await fetch(input, init);
    console.log(
        `${new Date().toISOString().slice(-13)} 📡 Received ${
            body.operationName
        } response in ${Date.now() - start}ms`
    );

    return {
        ...response,

        async text() {
            const start = Date.now();
            const result = await response.text();
            console.log(JSON.parse(result));
            console.log(
                `${new Date().toISOString().slice(-13)} ⚙️ in ${
                    Date.now() - start
                }ms (${result.length} bytes)`
            );
            return result;
        },
    };
}

let apolloLink = new HttpLink({ uri: API_URL });
if (process.env.REACT_APP_DEBUG === "true") {
    apolloLink = new HttpLink({fetch: loggingFetch, uri: API_URL});
}

export const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: apolloLink,
});
