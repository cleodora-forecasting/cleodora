import {render, screen, within} from '@testing-library/react'
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

    const footer = screen.getByRole('contentinfo');

    expect(within(footer).getByText("cleodora.org")).toBeInTheDocument();

    await within(footer).findByText("99.99.99+test");
    await within(footer).findByText("Footer text for a test");
});
