import React, {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
import {
    CreateForecastMutation,
    CreateForecastMutationVariables,
} from "./__generated__/graphql";
import {TextField, Button} from "@mui/material";
import {DateTimePicker, LocalizationProvider} from "@mui/x-date-pickers";
import dayjs, {Dayjs} from "dayjs";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";

const ADD_FORECAST = gql(`
    mutation createForecast($forecast: NewForecast!, $estimate: NewEstimate!) {
        createForecast(forecast: $forecast, estimate: $estimate) {
            id
            title
        }
    }
`);
//` as DocumentNode<CreateForecastMutation, CreateForecastMutationVariables>);
export const AddForecast: FC = () => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [closes, setCloses] = React.useState<Dayjs | null>(dayjs());
    const [resolves, setResolves] = React.useState<Dayjs | null>(dayjs());

//    useMutation<TableSizeMutation, TableSizeMutationVariables>
    const [addForecast, {error, data}] = useMutation<CreateForecastMutation, CreateForecastMutationVariables>(ADD_FORECAST, {
        refetchQueries: [
            {query: GET_FORECASTS}, // TODO needs refactor, should not be
            // referencing different component stuff
            'GetForecasts' // Query name
        ],
        variables: {
            forecast: {
                title,
                description,
                closes,
                resolves,
            },
            estimate: {
                reason: "TODO Just a gut feeling I have.",
                probabilities: [
                    {
                        value: 50,
                        outcome: {
                            text: "TODO Yes",
                        },
                    },
                    {
                        value: 50,
                        outcome: {
                            text: "TODO No",
                        },
                    },
                ],
            }
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
            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <form
                    onSubmit={e => {
                        e.preventDefault();
                        addForecast().then(() => {
                            setTitle('');
                            setDescription('');
                            setCloses(dayjs());
                            setResolves(dayjs());
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
                        <DateTimePicker
                            label="Closes"
                            value={closes}
                            onChange={(newValue: Dayjs | null) => {
                                setCloses(newValue);
                            }}
                            renderInput={(params) => <TextField {...params} />}
                        />
                    </p>
                    <p>
                        <DateTimePicker
                            label="Resolves"
                            value={resolves}
                            onChange={(newValue: Dayjs | null) => {
                                setResolves(newValue);
                            }}
                            renderInput={(params) => <TextField {...params} />}
                        />
                    </p>
                    <Button variant="outlined" type="submit">Add Forecast</Button>
                </form>
            </LocalizationProvider>
        </div>
    );
}
