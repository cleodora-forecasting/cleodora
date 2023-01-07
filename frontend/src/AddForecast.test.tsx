import {render, screen} from '@testing-library/react';
import AddForecast from './AddForecast';
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import userEvent from "@testing-library/user-event";
import {server} from "./mocks/server";
import {graphql} from "msw";

test('after adding a forecast a success msg is shown', async () => {
    const user = userEvent.setup()
    let requestBody: {
        variables: {
            forecast:
                {
                    title: string;
                    resolves: string;
                    closes: string;
                    description: string;
                };
            estimate: {
                reason: string;
                probabilities: {value: number; outcome: {text: string}; }[];
            };
        };
    } | undefined;

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
    const expectedResolves = "2023-01-13T10:00:00.000Z";
    const expectedReason = "It was written carefully and is not complicated.";

    await user.type(screen.getByLabelText('Title'), expectedTitle);
    await user.clear(screen.getByLabelText('Closes'));
    await user.clear(screen.getByLabelText('Resolves'));
    await user.type(screen.getByLabelText('Resolves'), inputResolves);
    await user.type(screen.getByLabelText('Reason'), expectedReason);
    await user.type(await screen.findByLabelText('Outcome0'), 'Yes');
    await user.type(await screen.findByLabelText('Probability0'), '95');
    await user.click(await screen.findByLabelText('add probability'));
    await user.type(await screen.findByLabelText('Outcome1'), 'No');
    await user.type(await screen.findByLabelText('Probability1'), '5');
    await user.click(await screen.findByRole("button", {name: "Add Forecast"}));

    expect(await screen.findByText('Saved "Mock title" with ID 999.')).toBeInTheDocument();
    expect(requestBody).toBeTruthy();
    if (!requestBody) {
        return
    }
    expect(requestBody.variables.forecast.title).toBe(expectedTitle);

    expect(requestBody.variables.forecast.resolves).toBe(expectedResolves);
    expect(requestBody.variables.forecast.closes).toBeNull();
    expect(requestBody.variables.forecast.description).toBe("");

    // Probability estimate
    expect(requestBody.variables.estimate.reason).toBe(expectedReason);
    expect(requestBody.variables.estimate.probabilities).toHaveLength(2);
    let probYes = requestBody.variables.estimate.probabilities[0];
    let probNo = requestBody.variables.estimate.probabilities[1];
    if (probYes.value !== 95) {
        // switch around, since we shouldn't assume order
        probYes = requestBody.variables.estimate.probabilities[1];
        probNo = requestBody.variables.estimate.probabilities[0];
    }
    expect(probYes.value).toBe(95);
    expect(probYes.outcome.text).toBe('Yes');
    expect(probNo.value).toBe(5);
    expect(probNo.outcome.text).toBe('No');
}, 15000);
