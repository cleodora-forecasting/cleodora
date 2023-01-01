import {render, waitFor, screen} from '@testing-library/react'
import '@testing-library/jest-dom'
import {ApolloProvider} from "@apollo/client";
import {Footer} from './Footer';
import {client} from "./client";

test('footer displays a version and a footer text', async () => {
    render(
        <ApolloProvider client={client}>
            <Footer/>
        </ApolloProvider>
    );

    expect(screen.getByText('cleodora.org')).toBeInTheDocument();

    await waitFor(() => screen.findByText("99.99.99+test"))
    expect(screen.getByText("99.99.99+test")).toBeInTheDocument();
    expect(screen.getByText("Footer text for a test")).toBeInTheDocument();
});
