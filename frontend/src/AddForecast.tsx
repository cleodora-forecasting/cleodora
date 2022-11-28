import {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
import {
    CreateForecastMutation,
    CreateForecastMutationVariables,
} from "./__generated__/graphql";

const ADD_FORECAST = gql(`
    mutation createForecast($input: NewForecast!) {
        createForecast(input: $input) {
            id
            summary
        }
    }
`);
//` as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>);
export const AddForecast: FC = () => {
    const [summary, setSummary] = useState('');
    const [description, setDescription] = useState('');
    const [closes, setCloses] = useState(''); // TODO date
    const [resolves, setResolves] = useState(''); // TODO date

//    useMutation<TableSizeMutation, TableSizeMutationVariables>
    const [addForecast, {error, data}] = useMutation<CreateForecastMutation, CreateForecastMutationVariables>(ADD_FORECAST, {
        refetchQueries: [
            {query: GET_FORECASTS}, // TODO needs refactor, should not be
            // referencing different component stuff
            'GetForecasts' // Query name
        ],
        variables: {
            input: {
                summary,
                description,
                closes,
                resolves,
            },
        },
    });

    return (
        <div>
            <h3>Add Forecast</h3>
            {error ?
                <p style={{color: "red"}}>Oh no! {error.message}</p> : null}
            {data && data.createForecast ?
                <p style={{color: "green"}}>
                    Saved "{data.createForecast.summary}" with
                    ID {data.createForecast.id}.
                </p>
                : null}
            <form
                onSubmit={e => {
                    e.preventDefault();
                    addForecast().then(() => {
                        setSummary('');
                        setDescription('');
                        setCloses('');
                        setResolves('');
                    }).catch(reason => (console.log("error addForecast()", reason)));
                }}
            >
                <p>
                    <label>Summary</label>
                    <input
                        name="summary"
                        value={summary}
                        onChange={e => setSummary(e.target.value)}
                    />
                </p>
                <p>
                    <label>Description</label>
                    <textarea
                        name="description"
                        value={description}
                        onChange={e => setDescription(e.target.value)}
                    />
                </p>
                <p>
                    <small>Format for the dates:
                        2022-12-01T09:00:00+01:00</small>
                    <br/>
                    <label>Closes</label>
                    <input
                        name="closes"
                        value={closes}
                        onChange={e => setCloses(e.target.value)}
                    />
                </p>
                <p>
                    <label>Resolves</label>
                    <input
                        name="resolves"
                        value={resolves}
                        onChange={e => setResolves(e.target.value)}
                    />
                </p>
                <button type="submit">Add Forecast</button>
            </form>
        </div>
    );
}
