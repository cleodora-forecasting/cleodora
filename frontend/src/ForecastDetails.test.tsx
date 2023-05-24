import {render, screen} from "@testing-library/react";
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import {ForecastDetails} from "./ForecastDetails";
import {Forecast, Resolution} from "./__generated__/graphql";

test('the unresolved forecast is displayed', () => {
    const forecast: Forecast = {
        created: "2022-10-30T17:05:00+01:00",
        description: "If it rains continuously for more than 30 minutes, that" +
            " counts as rain.",
        estimates: [
            {
                id: "E01",
                created: "2022-10-30T17:05:00+01:00",
                reason: "Just a hunch.",
                brierScore: null,
                probabilities: [
                    {
                        id: "P01",
                        value: 20,
                        outcome: {
                            id: "O01",
                            correct: false,
                            text: "Yes",
                        },
                    },
                    {
                        id: "P02",
                        value: 80,
                        outcome: {
                            id: "O02",
                            correct: false,
                            text: "No",
                        },
                    },
                ],
            },
        ],
        id: "",
        resolution: Resolution.Unresolved,
        resolves: "2022-11-30T17:05:00+01:00",
        title: "Will it rain tomorrow?"
    };

    render(
        <ApolloProvider client={client}>
            <ForecastDetails forecast={forecast} open={true} handleClose={() => null}/>
        </ApolloProvider>
    );

    expect(screen.getByRole("heading", {name: "Will it rain tomorrow?"})).toBeInTheDocument();
    expect(screen.getByText(/more than 30 minutes/)).toBeInTheDocument();
    expect(screen.getByText(/UNRESOLVED/)).toBeInTheDocument();
    expect(screen.getByText(/Just a hunch./)).toBeInTheDocument();
    expect(screen.getByText(/Brier Score:/)).toBeInTheDocument();
    expect(screen.queryByText("(Yes)")).not.toBeInTheDocument();
    expect(screen.queryByText("(No)")).not.toBeInTheDocument();
});

test('the resolved forecast is displayed', () => {
    const forecast: Forecast = {
        created: "2022-10-30T17:05:00+01:00",
        description: "If it rains continuously for more than 30 minutes, that" +
            " counts as rain.",
        estimates: [
            {
                id: "E01",
                created: "2022-10-30T17:05:00+01:00",
                reason: "Just a hunch.",
                brierScore: 1.28,
                probabilities: [
                    {
                        id: "P01",
                        value: 20,
                        outcome: {
                            id: "O01",
                            correct: true,
                            text: "Yes",
                        },
                    },
                    {
                        id: "P02",
                        value: 80,
                        outcome: {
                            id: "O02",
                            correct: false,
                            text: "No",
                        },
                    },
                ],
            },
        ],
        id: "",
        resolution: Resolution.Resolved,
        resolves: "2022-11-30T17:05:00+01:00",
        title: "Will it rain tomorrow?"
    };

    render(
        <ApolloProvider client={client}>
            <ForecastDetails forecast={forecast} open={true} handleClose={() => null}/>
        </ApolloProvider>
    );

    expect(screen.getByRole("heading", {name: "Will it rain tomorrow?"})).toBeInTheDocument();
    expect(screen.getByText(/more than 30 minutes/)).toBeInTheDocument();
    expect(screen.getByText(/RESOLVED/)).toBeInTheDocument();
    expect(screen.getByText("(Yes)")).toBeInTheDocument(); // next to RESOLVED
    expect(screen.getByText(/Just a hunch./)).toBeInTheDocument();
    expect(screen.getByText((_, e) => e?.textContent === 'Brier Score: 1.28'))
        .toBeInTheDocument();
});
