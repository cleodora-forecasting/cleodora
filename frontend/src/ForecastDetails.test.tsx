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
        estimates: [],
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
});
