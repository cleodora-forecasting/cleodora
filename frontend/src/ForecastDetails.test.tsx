import {render, screen} from "@testing-library/react";
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import {ForecastDetails} from "./ForecastDetails";
import {Forecast, Resolution} from "./__generated__/graphql";

test('the forecast is displayed', () => {
    const forecast: Forecast = {
        created: undefined,
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
                ],
            },
        ],
        id: "",
        resolution: Resolution.Unresolved,
        resolves: undefined,
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
});
