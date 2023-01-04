import { render, screen } from '@testing-library/react';
import App from './App';
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";

test('complete overview is rendered', async () => {
    render(
        <ApolloProvider client={client}>
            <App />
        </ApolloProvider>
    );

    expect(screen.getByText("Cleodora")).toBeInTheDocument();
    expect(screen.getByText("cleodora.org")).toBeInTheDocument();
    expect(await screen.findByRole("heading", {name: "Forecasts"})).toBeInTheDocument();
    expect(await screen.findByRole("heading", {name: "Add Forecast"})).toBeInTheDocument();
});
