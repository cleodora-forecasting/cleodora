import React, {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
import {NewProbability} from "./__generated__/graphql";
import {
    CreateForecastMutation,
    CreateForecastMutationVariables,
} from "./__generated__/graphql";
import {
    TextField,
    Button,
    Grid,
    IconButton,
    InputAdornment
} from "@mui/material";
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
                    <Grid container direction="column" alignItems="flex-start" spacing={3} justifyItems="flex-start">
                        <Grid item>
                            <TextField
                                required
                                label="Title"
                                value={title}
                                onChange={e => setTitle(e.target.value)}
                                variant="outlined"
                            />
                        </Grid>
                        <Grid item>
                            <TextField
                                value={description}
                                onChange={e => setDescription(e.target.value)}
                                label="Description"
                                multiline
                                variant="outlined"
                            />
                        </Grid>
                        <Grid item>
                            <DateTimePicker
                                label="Resolves"
                                value={resolves}
                                onChange={(newValue) => {
                                    if (!newValue) {
                                        return
                                    }
                                    setResolves(newValue);
                                }}
                                renderInput={(params) => (
                                    <TextField
                                        {...params}
                                        required
                                        helperText="Date when you'll know the answer."
                                    />
                                )}
                            />
                        </Grid>
                        <Grid item>
                            <DateTimePicker
                                label="Closes"
                                value={closes}
                                onChange={(newValue: Dayjs | null) => {
                                    setCloses(newValue);
                                }}
                                renderInput={(params) => (
                                    <TextField
                                        {...params}
                                        helperText="Optional date when you'll stop updating."
                                    />
                                )}
                            />
                        </Grid>
                        <Grid item container direction="column" spacing={2}>
                            <Grid item>
                                <p style={{maxWidth: 400}}>Specify all possible outcomes/answers to the forecast, each with a probability of 0-100% . The total probability must add up to 100% . For example "Yes" with 30% and "No" with 70%.</p>
                            </Grid>
                        {probabilities.map((prob, index, {length}) => (
                            <Grid item container key={index} spacing={1} alignItems="center">
                                <Grid item>
                                    <TextField
                                        required
                                        value={prob.outcome.text}
                                        onChange={e => handleModifyProbability(index, e.target.value, prob.value)}
                                        label={`${index+1}. Outcome`}
                                        variant="outlined"
                                    />
                                </Grid>
                                <Grid item>
                                    <TextField
                                        required
                                        value={prob.value}
                                        onChange={e => {
                                            if (isNaN(Number(e.target.value))) {
                                                return
                                            }
                                            handleModifyProbability(index, prob.outcome.text, Number(e.target.value))
                                        }}
                                        inputProps={{inputMode: "numeric", pattern: '[0-9]+'}}
                                        InputLabelProps={{shrink: true}}
                                        label={`${index+1}. Probability`}
                                        variant="outlined"
                                        sx={{ m: 1, width: '25ch' }}
                                        InputProps={{
                                            startAdornment: <InputAdornment position="start">%</InputAdornment>,
                                        }}
                                    />
                                </Grid>
                                {index + 1 === length ?
                                    <Grid item>
                                        <IconButton
                                            aria-label="add probability"
                                            onClick={_ => setProbabilities(old => [...old, {outcome: {text: ''}, value: 0}])}
                                        >
                                            <AddCircleOutlineIcon />
                                        </IconButton>
                                    </Grid>
                                    : null
                                }
                            </Grid>
                        ))}
                        </Grid>
                        <Grid item>
                            <TextField
                                required
                                value={reason}
                                onChange={e => setReason(e.target.value)}
                                label="Reason"
                                multiline
                                variant="outlined"
                                helperText="Why these probabilities?"
                            />
                        </Grid>
                        <Grid item>
                            {error ?
                                <p style={{color: "red"}}>Oh no! {error.message}</p> : null}
                            {data && data.createForecast ?
                                <p style={{color: "green"}}>
                                    Saved "{data.createForecast.title}" with
                                    ID {data.createForecast.id}.
                                </p>
                                : null}
                            <Button variant="outlined" type="submit">Add Forecast</Button>
                        </Grid>
                    </Grid>
                </form>
            </LocalizationProvider>
        </div>
    );
}

export default AddForecast;
