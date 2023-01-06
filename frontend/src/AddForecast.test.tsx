import {render, screen} from '@testing-library/react';
import AddForecast from './AddForecast';
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import userEvent from "@testing-library/user-event";
import {server} from "./mocks/server";
import {graphql} from "msw";

test('after adding a forecast a success msg is shown', async () => {
    const user = userEvent.setup()
    let requestBody;

    server.use(
        graphql.mutation("createForecast", async (req, res, ctx) => {
            requestBody = await req.json();
            return res(
                ctx.data({
                    "createForecast": {
                        "id": "999",
                        "title": "Mock title",
                        "__typename": "Forecast"
                    }
                }),
            )
        }),
    )

    render(
        <ApolloProvider client={client}>
            <AddForecast />
        </ApolloProvider>
    );

    const expectedTitle = "Will this test pass?";
    const inputResolves = "01/13/2023 10:00 AM";
    const expectedResolves = "2023-01-13T09:00:00.000Z";

    await user.type(screen.getByLabelText('Title'), expectedTitle);
    await user.clear(screen.getByLabelText('Closes'));
    await user.clear(screen.getByLabelText('Resolves'));
    await user.type(screen.getByLabelText('Resolves'), inputResolves);
    await user.click(await screen.findByRole("button", {name: "Add Forecast"}));

    expect(await screen.findByText('Saved "Mock title" with ID 999.')).toBeInTheDocument();
    expect(requestBody).toBeTruthy();
    if (requestBody) {
        expect(requestBody!.variables.forecast.title).toBe(expectedTitle);
        expect(requestBody!.variables.forecast.resolves).toBe(expectedResolves);
        expect(requestBody!.variables.forecast.closes).toBeNull();
        expect(requestBody!.variables.forecast.description).toBe("");
    }
});
