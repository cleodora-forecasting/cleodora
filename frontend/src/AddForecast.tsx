import React, {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
import {NewProbability} from "./__generated__/graphql";
import {
    CreateForecastMutation,
    CreateForecastMutationVariables,
} from "./__generated__/graphql";
import {TextField, Button, Grid, IconButton} from "@mui/material";
import {DateTimePicker, LocalizationProvider} from "@mui/x-date-pickers";
import dayjs, {Dayjs} from "dayjs";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';

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
    const [closes, setCloses] = React.useState<Dayjs | null>(null);
    const [resolves, setResolves] = useState(dayjs());
    const [reason, setReason] = useState('');
    const [probabilities, setProbabilities] = useState([{outcome: {text: ''}, value: 0} as NewProbability]);

//    useMutation<TableSizeMutation, TableSizeMutationVariables>
    const [addForecast, {error, data}] = useMutation<CreateForecastMutation, CreateForecastMutationVariables>(ADD_FORECAST, {
        refetchQueries: [
            {query: GET_FORECASTS}, // TODO needs refactor, should not be
            // referencing different component stuff
            // to refetch outside of the mutation i.e. inside <App />?
            // https://www.apollographql.com/docs/react/data/refetching/
        ],
        variables: {
            forecast: {
                title,
                description,
                closes,
                resolves,
            },
            estimate: {
                reason,
                probabilities
            }
        },
    });

    // https://beta.reactjs.org/learn/updating-arrays-in-state#replacing-items-in-an-array
    function handleModifyProbability(index: number, outcome: string, value: number) {
        const nextProbabilities = probabilities.map((c, i) => {
            if (i === index) {
                // Increment the clicked counter
                return {outcome: {text: outcome}, value: value};
            } else {
                // The rest haven't changed
                return c;
            }
        });
        setProbabilities(nextProbabilities);
    }

    // https://rajputankit22.medium.com/add-dynamically-textfields-in-react-js-71320aee9a8d
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
                <form onSubmit={e => {
                    e.preventDefault();
                    addForecast().then(() => {
                        setTitle('');
                        setDescription('');
                        setCloses(null);
                        setResolves(dayjs());
                        setReason('');
                        setProbabilities([{value: 0, outcome: {text: ''}}]);
                    }).catch(reason => (console.log("error addForecast()", reason)));
                }}>
                    <Grid container direction="row" alignItems="flex-start" spacing={1} justifyItems="flex-start">
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <TextField
                                    label="Title"
                                    value={title}
                                    onChange={e => setTitle(e.target.value)}
                                    variant="filled"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    value={description}
                                    onChange={e => setDescription(e.target.value)}
                                    label="Description"
                                    multiline
                                    variant="filled"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <DateTimePicker
                                    label="Closes"
                                    value={closes}
                                    onChange={(newValue: Dayjs | null) => {
                                        setCloses(newValue);
                                    }}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            helperText="(optional) No probability updates after this date."
                                        />
                                    )}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <DateTimePicker
                                    label="Resolves"
                                    value={resolves}
                                    onChange={(newValue) => {
                                        if (!newValue) {
                                            return
                                        }
                                        setResolves(newValue);
                                    }}
                                    renderInput={(params) => <TextField {...params} />}
                                />
                            </Grid>
                            <Grid item xs={12}>
                            {probabilities.map((prob, index) => (
                                <Grid container key={index}>
                                    <Grid item xs={2}>
                                        <TextField
                                            value={prob.outcome.text}
                                            onChange={e => handleModifyProbability(index, e.target.value, prob.value)}
                                            label={`Outcome${index}`}
                                            variant="filled"/>
                                    </Grid>
                                    <Grid item xs={2}>
                                        <TextField
                                            value={prob.value}
                                            inputProps={{ inputMode: 'numeric', pattern: '[0-9]*' }}
                                            onChange={
                                                e => {
                                                    if (isNaN(Number(e.target.value))) {
                                                        return
                                                    }
                                                    handleModifyProbability(index, prob.outcome.text, Number(e.target.value))
                                                }
                                            }
                                            label={`Probability${index}`}
                                            variant="filled"/>
                                    </Grid>
                                    <Grid item xs={1}>
                                        <IconButton
                                            aria-label="add probability"
                                            onClick={_ => setProbabilities(old => [...old, {outcome: {text: ''}, value: 0}])}
                                        >
                                            <AddCircleOutlineIcon />
                                        </IconButton>
                                    </Grid>
                                </Grid>
                            ))}
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    value={reason}
                                    onChange={e => setReason(e.target.value)}
                                    label="Reason"
                                    multiline
                                    variant="filled"
                                    helperText="Why these probabilities?"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <Button variant="outlined" type="submit">Add Forecast</Button>
                            </Grid>
                        </Grid>
                    </Grid>
                </form>
            </LocalizationProvider>
        </div>
    );
}

export default AddForecast;
