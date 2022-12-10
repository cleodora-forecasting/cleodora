import {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
import {
    CreateForecastMutation,
    CreateForecastMutationVariables,
} from "./__generated__/graphql";
import Button from "@mui/material/Button/Button";
import {TextField} from "@mui/material";

const ADD_FORECAST = gql(`
    mutation createForecast($input: NewForecast!) {
        createForecast(input: $input) {
            id
            title
        }
    }
`);
//` as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>);
export const AddForecast: FC = () => {
    const [title, setTitle] = useState('');
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
                title,
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
                    Saved "{data.createForecast.title}" with
                    ID {data.createForecast.id}.
                </p>
                : null}
            <form
                onSubmit={e => {
                    e.preventDefault();
                    addForecast().then(() => {
                        setTitle('');
                        setDescription('');
                        setCloses('');
                        setResolves('');
                    }).catch(reason => (console.log("error addForecast()", reason)));
                }}
            >
                <p>
                    <TextField
                        label="Title"
                        value={title}
                        onChange={e => setTitle(e.target.value)}
                        variant="filled"
                    />
                </p>
                <p>
                    <TextField
                        value={description}
                        onChange={e => setDescription(e.target.value)}
                        label="Description"
                        multiline
                        variant="filled"
                    />
                </p>
                <p>
                    <small>Format for the dates:
                        2022-12-01T09:00:00+01:00</small>
                    <br/>
                    <TextField
                        label="Closes"
                        value={closes}
                        onChange={e => setCloses(e.target.value)}
                        variant="filled"
                    />
                </p>
                <p>
                    <TextField
                        label="Resolves"
                        value={resolves}
                        onChange={e => setResolves(e.target.value)}
                        variant="filled"
                    />
                </p>
                <Button variant="outlined" type="submit">Add Forecast</Button>
            </form>
        </div>
    );
}
