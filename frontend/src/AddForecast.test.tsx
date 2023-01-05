import { render, screen } from '@testing-library/react';
import AddForecast from './AddForecast';
import {ApolloProvider} from "@apollo/client";
import {client} from "./client";
import userEvent from "@testing-library/user-event";

test('after adding a forecast a succcess msg is shown', async () => {
    const user = userEvent.setup()

    render(
        <ApolloProvider client={client}>
            <AddForecast />
        </ApolloProvider>
    );

    // fill out form

    const title = screen.getByLabelText('Title');
    const resolves = screen.getByLabelText('Resolves');

    await user.type(title, "Will this test pass?");
    await user.type(resolves, "2022-01-31T12:00:00+01:00");

    // click submit
    const submit = await screen.findByRole("button", {name: "Add Forecast"});

    await user.click(submit);

    // TODO improve test to not use a standard mock but define it here to have
    // the same title and verify the request too.
    // https://www.stackbuilders.com/blog/testing-react-components-with-testing-library-and-mock-service-worker/
    // https://kentcdodds.com/blog/stop-mocking-fetch
    expect(await screen.findByText(/Saved "Mock title" with ID \d+./)).toBeInTheDocument();
});
