import {
    render,
    screen, waitFor,
} from "@testing-library/react";
import {Forecast, Resolution} from "./__generated__/graphql";
import {client} from "./client";
import {ApolloProvider} from "@apollo/client";
import {ResolveForecastDialog} from "./ResolveForecastDialog";
import userEvent from "@testing-library/user-event";
import {server} from "./mocks/server";
import {graphql} from "msw";

function addDays(date: Date, days: number) {
    const result = new Date(date);
    result.setDate(result.getDate() + days);
    return result;
}

test('the resolve dialog is displayed', () => {
    const now = new Date();
    const forecast: Forecast = {
        id: "f01",
        resolution: Resolution.Unresolved,
        resolves: addDays(now, 30),
        title: "Will it rain tomorrow?",
        created: now,
        description: "If it rains continuously for more than 30 minutes, that" +
            " counts as rain.",
        estimates: [
            {
                id: "e01",
                reason: "The weather report says so",
                created: now,
                brierScore: null,
                probabilities: [
                    {
                        id: "p01",
                        value: 20,
                        outcome: {
                            id: "o01",
                            text: "Yes",
                            correct: false,
                        },
                    },
                    {
                        id: "p02",
                        value: 80,
                        outcome: {
                            id: "o02",
                            text: "No",
                            correct: false,
                        },
                    },
                ],
            },
        ],
    };

    render(
        <ApolloProvider client={client}>
            <ResolveForecastDialog forecast={forecast} open={true} handleClose={() => null} />
        </ApolloProvider>
    );

    expect(screen.getByRole("heading", {name: "Will it rain tomorrow?"})).toBeInTheDocument();
    expect(screen.getByText("Yes")).toBeInTheDocument();
    expect(screen.getByText("No")).toBeInTheDocument();
    expect(screen.getByText("No correct outcome (resolve forecast as N/A)")).toBeInTheDocument();
});

type ResolveForecastRequest = {
    variables: {
        forecastId: string,
        resolution: string,
        correctOutcomeId: string,
    };
}

test('forecast can be resolved with outcome', async () => {
    const user = userEvent.setup();
    let requestBody: ResolveForecastRequest | undefined;

    server.use(
        graphql.mutation("resolveForecast", async (req, res, ctx) => {
            requestBody = await req.json();
            // Simulate a delay like a real request
            await new Promise(f => setTimeout(f, 200));
            return res(
                ctx.data({
                    "resolveForecast": {
                        "id": "999",
                        "title": "Mock title",
                        "resolution": "RESOLVED",
                        "resolves": addDays(now, 30),
                        "closes": null,
                        "estimates": [
                            {
                                "id": "e01",
                                brierScore: 2.0,
                                probabilities: [
                                    {
                                        id: "p01",
                                        outcome: {
                                            id: "o01",
                                            text: "Yes",
                                            correct: true,
                                        },
                                    },
                                    {
                                        id: "p02",
                                        outcome: {
                                            id: "o02",
                                            text: "No",
                                            correct: false,
                                        },
                                    },
                                ],
                            },
                        ],
                        "__typename": "Forecast"
                    }
                }),
            )
        }),
    );

    const now = new Date();
    const forecast: Forecast = {
        id: "f01",
        resolution: Resolution.Unresolved,
        resolves: addDays(now, 30),
        title: "Will it rain tomorrow?",
        created: now,
        description: "If it rains continuously for more than 30 minutes, that" +
            " counts as rain.",
        estimates: [
            {
                id: "e01",
                reason: "The weather report says so",
                created: now,
                brierScore: null,
                probabilities: [
                    {
                        id: "p01",
                        value: 20,
                        outcome: {
                            id: "o01",
                            text: "Yes",
                            correct: false,
                        },
                    },
                    {
                        id: "p02",
                        value: 80,
                        outcome: {
                            id: "o02",
                            text: "No",
                            correct: false,
                        },
                    },
                ],
            },
        ],
    };

    const dialogCloseHandler = jest.fn();

    render(
        <ApolloProvider client={client}>
            <ResolveForecastDialog forecast={forecast} open={true} handleClose={dialogCloseHandler} />
        </ApolloProvider>
    );

    expect(screen.getByRole("heading", {name: "Will it rain tomorrow?"})).toBeInTheDocument();
    await user.click(screen.getByLabelText('Yes'));
    await user.click(screen.getByRole("button", {name: "Save"}));

    await waitFor(() => expect(dialogCloseHandler).toHaveBeenCalledTimes(1));

    expect(requestBody).toBeTruthy();
    if (!requestBody) {
        return
    }
    expect(requestBody.variables.forecastId).toBe("f01");
    expect(requestBody.variables.correctOutcomeId).toBe("o01");
    // The API assumes RESOLVED when no resolution is passed
    expect(requestBody.variables.resolution).toBeUndefined();
});

test('error is displayed', async () => {
    const user = userEvent.setup();
    let requestBody: ResolveForecastRequest | undefined;

    server.use(
        graphql.mutation("resolveForecast", async (req, res, ctx) => {
            requestBody = await req.json();
            await new Promise(f => setTimeout(f, 200));
            return res(
                ctx.errors([
                    {
                        name: "whatever",
                        message: "API test error",
                    }
                ])
            )
        }),
    );

    const now = new Date();
    const forecast: Forecast = {
        id: "f01",
        resolution: Resolution.Unresolved,
        resolves: addDays(now, 30),
        title: "Will it rain tomorrow?",
        created: now,
        description: "If it rains continuously for more than 30 minutes, that" +
            " counts as rain.",
        estimates: [
            {
                id: "e01",
                reason: "The weather report says so",
                created: now,
                brierScore: null,
                probabilities: [
                    {
                        id: "p01",
                        value: 20,
                        outcome: {
                            id: "o01",
                            text: "Yes",
                            correct: false,
                        },
                    },
                    {
                        id: "p02",
                        value: 80,
                        outcome: {
                            id: "o02",
                            text: "No",
                            correct: false,
                        },
                    },
                ],
            },
        ],
    };

    const dialogCloseHandler = jest.fn();

    render(
        <ApolloProvider client={client}>
            <ResolveForecastDialog forecast={forecast} open={true} handleClose={dialogCloseHandler} />
        </ApolloProvider>
    );

    expect(screen.getByRole("heading", {name: "Will it rain tomorrow?"})).toBeInTheDocument();
    await user.click(screen.getByLabelText('Yes'));
    await user.click(screen.getByRole("button", {name: "Save"}));

    const errFinder = (content: string, element: Element | null): boolean => {
        const message = "ApolloError: API test error";
        if (element &&
            element.tagName.toLowerCase() === 'p' &&
            element.textContent === message
        ) {
            // eslint-disable-next-line jest/no-conditional-expect
            expect(element).toHaveStyle("color: red");
            return true;
        }
        return false;
    }

    expect(await screen.findByText(errFinder)).toBeTruthy();

    // The function to close the dialog should not have been called. It's
    // safer to check for this after looking for the error messages because
    // closing the dialog is only attempted once the network requests
    // complete, so checking immediately clicking "Save" could lead to a
    // false negative.
    expect(dialogCloseHandler).not.toHaveBeenCalled();

    expect(requestBody).toBeTruthy();
    if (!requestBody) {
        return
    }
    expect(requestBody.variables.forecastId).toBe("f01");
    expect(requestBody.variables.correctOutcomeId).toBe("o01");
    // The API assumes RESOLVED when no resolution is passed
    expect(requestBody.variables.resolution).toBeUndefined();
});
