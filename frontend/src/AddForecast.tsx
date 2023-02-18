import React, {FC, useState} from "react";
import {useMutation} from "@apollo/client";
import {GET_FORECASTS} from "./ForecastList";
import {gql} from "./__generated__"
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
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import {v4 as uuid} from "uuid";

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
    const initialOutcomes = [
        {"id": uuid(), "outcome": "", value: 0},
        {"id": uuid(), "outcome": "", value: 0},
    ];
    const [outcomes, setOutcomes] = useState(initialOutcomes);

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
                probabilities: outcomes.map(v => {return {outcome: {text: v.outcome}, value: v.value}}),
            }
        },
    });

    let messageBox:JSX.Element = <></>;
    if (error) {
        const errMessages = new Array(error.message);
        if (error.message.includes('NetworkError')) {
            errMessages.push('Make sure the API is running and reachable.');
        }
        messageBox =
            <p style={{color:"red"}}>
                <strong>{error.name}: </strong>
                {errMessages.join(' ')}
            </p>
    } else if (data && data.createForecast) {
        messageBox =
            <p style={{color: "green"}}>
                Saved "{data.createForecast.title}" with
                ID {data.createForecast.id}.
            </p>
        }


    // https://beta.reactjs.org/learn/updating-arrays-in-state#replacing-items-in-an-array
    function updateOutcome(id: string, outcome: string, value: number) {
        setOutcomes(outcomes.map(p => {
            if (id === p.id) {
                return {"id": p.id, "outcome": outcome, "value": value};
            } else {
                return p;
            }
        }));
    }

    function deleteOutcome(id: string) {
        setOutcomes(outcomes.filter(value => value.id !== id));
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
                        setOutcomes(initialOutcomes);
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
                        <Grid item container direction="column" spacing={1}>
                        {outcomes.map((prob, index) =>  (
                            <Grid item container key={prob.id} spacing={1} alignItems="center">
                                <Grid item>
                                    <TextField
                                        required
                                        value={prob.outcome}
                                        onChange={e => updateOutcome(prob.id, e.target.value, prob.value)}
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
                                            updateOutcome(prob.id, prob.outcome, Number(e.target.value))
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
                                <Grid item>
                                    <IconButton
                                        style={{color: 'darkred'}}
                                        aria-label="delete outcome"
                                        onClick={_ => deleteOutcome(prob.id)}
                                    >
                                        <DeleteOutlineIcon />
                                    </IconButton>
                                </Grid>
                            </Grid>
                        ))}
                            <Grid item>
                                <Button
                                    size="small"
                                    variant="outlined"
                                    startIcon={<AddCircleOutlineIcon />}
                                    aria-label="add outcome"
                                    onClick={_ => setOutcomes(old => [...old, {"id": uuid(), "outcome": "", "value": 0}])}
                                >
                                    Outcome
                                </Button>
                            </Grid>
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
                            {messageBox}
                            <Button variant="outlined" type="submit">Add Forecast</Button>
                        </Grid>
                    </Grid>
                </form>
            </LocalizationProvider>
        </div>
    );
}

export default AddForecast;
