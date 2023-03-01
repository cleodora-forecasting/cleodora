import {render, screen} from "@testing-library/react";
import {Forecast, Resolution} from "./__generated__/graphql";
import {client} from "./client";
import {ApolloProvider} from "@apollo/client";
import {ResolveForecastDialog} from "./ResolveForecastDialog";

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
