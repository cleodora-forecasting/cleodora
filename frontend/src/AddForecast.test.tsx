import {render, screen} from '@testing-library/react';
import AddForecast from './AddForecast';
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import userEvent from "@testing-library/user-event";
import {server} from "./mocks/server";
import {graphql} from "msw";

type AddForecastRequest = {
    variables: {
        forecast:
            {
                title: string;
                resolves: string;
                description: string;
            };
        estimate: {
            reason: string;
            probabilities: {value: number; outcome: {text: string}; }[];
        };
    };
}

test('after adding a forecast a success msg is shown', async () => {
    const user = userEvent.setup()
    let requestBody: AddForecastRequest | undefined;

    server.use(
        graphql.mutation("createForecast", async (req, res, ctx) => {
            requestBody = await req.json();
            await new Promise(f => setTimeout(f, 200));
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
    const inputResolves = ["01", "13", "2023", "10", "00", "AM"];
    const expectedResolves = "2023-01-13T10:00:00.000Z";
    const expectedReason = "It was written carefully and is not complicated.";

    await user.type(screen.getByLabelText('Title *'), expectedTitle);

    // For some reason when selecting the 'resolves' field the last
    // element (AM/PM) is selected, so we need to move 5 times to the left.
    await user.click(screen.getByLabelText('Resolves *'));
    await user.keyboard('{ArrowLeft}'.repeat(5));
    for (const value of inputResolves) {
        await user.keyboard(value);
    }
    await user.type(screen.getByLabelText('Reason *'), expectedReason);
    await user.type(screen.getByLabelText('1. Outcome *'), 'Yes');
    await user.type(screen.getByLabelText('1. Probability *'), '95');
    await user.type(screen.getByLabelText('2. Outcome *'), 'No');
    await user.type(screen.getByLabelText('2. Probability *'), '99');
    await user.click(screen.getByLabelText('add outcome'));
    await user.type(await screen.findByLabelText('3. Outcome *'), 'Maybe');
    await user.type(await screen.findByLabelText('3. Probability *'), '1');
    await user.click(screen.getByLabelText('delete 2. outcome'));
    await user.click(screen.getByLabelText('add outcome'));
    await user.type(await screen.findByLabelText('3. Outcome *'), 'No');
    await user.type(await screen.findByLabelText('3. Probability *'), '4');
    await user.click(screen.getByRole("button", {name: "Add Forecast"}));

    expect(await screen.findByText('Saved "Mock title" with ID 999.')).toBeInTheDocument();
    expect(requestBody).toBeTruthy();
    if (!requestBody) {
        return
    }
    expect(requestBody.variables.forecast.title).toBe(expectedTitle);

    expect(requestBody.variables.forecast.resolves).toBe(expectedResolves);
    expect(requestBody.variables.forecast.description).toBe("");

    // Probability estimate
    expect(requestBody.variables.estimate.reason).toBe(expectedReason);
    expect(requestBody.variables.estimate.probabilities).toHaveLength(3);
    const expectedProbabilities = new Map<string, number>([
        ['Yes', 95],
        ['No', 4],
        ['Maybe', 1],
    ]);
    requestBody.variables.estimate.probabilities.forEach(p => {
        expect(expectedProbabilities.get(p.outcome.text)).toBe(p.value);
        // ensure every outcome only appears once
        expectedProbabilities.delete(p.outcome.text);
    });
}, 15000);


test('add a forecast by tabbing with the keyboard', async () => {
    const user = userEvent.setup();
    let requestBody: AddForecastRequest | undefined;

    server.use(
        graphql.mutation("createForecast", async (req, res, ctx) => {
            requestBody = await req.json();
            await new Promise(f => setTimeout(f, 200));
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
    const inputResolves = ["01", "13", "2023", "10", "00", "AM"];
    const expectedResolves = "2023-01-13T10:00:00.000Z";
    const expectedReason = "It was written carefully and is not complicated.";

    await user.type(screen.getByLabelText('Title *'), expectedTitle);
    await user.tab(); // to description
    await user.tab(); // to resolves
    await user.keyboard('{ArrowLeft}');
    for (const value of inputResolves) {
        await user.keyboard(value);
    }
    await user.tab(); // to calendar icon
    await user.tab(); // to 1. Outcome
    await user.keyboard("Yes");
    await user.tab();
    await user.keyboard("95");
    await user.tab();
    await user.tab();
    await user.keyboard("No");
    await user.tab();
    await user.keyboard("4");
    await user.tab();
    await user.tab();
    await user.keyboard('{Enter}')
    await user.keyboard("Maybe");
    await user.tab();
    await user.keyboard("1");
    await user.tab();
    await user.tab();
    await user.tab();
    await user.keyboard(expectedReason);
    await user.tab();
    await user.keyboard('{Enter}');

    expect(await screen.findByText('Saved "Mock title" with ID 999.')).toBeInTheDocument();

    // Verify the GraphQL mutation contains the data we expect
    expect(requestBody).toBeTruthy();
    if (!requestBody) {
        return
    }
    expect(requestBody.variables.forecast.title).toBe(expectedTitle);
    expect(requestBody.variables.forecast.resolves).toBe(expectedResolves);
    expect(requestBody.variables.forecast.description).toBe("");

    // Probability estimate
    expect(requestBody.variables.estimate.reason).toBe(expectedReason);
    expect(requestBody.variables.estimate.probabilities).toHaveLength(3);
    const expectedProbabilities = new Map<string, number>([
        ['Yes', 95],
        ['No', 4],
        ['Maybe', 1],
    ]);
    requestBody.variables.estimate.probabilities.forEach(p => {
        expect(expectedProbabilities.get(p.outcome.text)).toBe(p.value);
        // ensure every outcome only appears once
        expectedProbabilities.delete(p.outcome.text);
    });
}, 15000);
