import React, {ChangeEvent, FC, useState} from "react";
import {
    Forecast,
    Resolution,
    ResolveForecastMutation,
    ResolveForecastMutationVariables
} from "./__generated__/graphql";
import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    FormControlLabel,
    FormLabel,
    Paper,
    Radio,
    RadioGroup,
} from "@mui/material";
import Draggable from "react-draggable";
import {gql} from "./__generated__";
import {useMutation} from "@apollo/client";

const PaperComponent: FC = (props) => {
    return (
        <Draggable handle="#draggable-dialog-title">
            <Paper {...props} />
        </Draggable>
    )
}

const RESOLVE_FORECAST = gql(`
    mutation resolveForecast($forecastId: ID!, $resolution: Resolution, $correctOutcomeId: ID) {
        resolveForecast(forecastId: $forecastId, correctOutcomeId: $correctOutcomeId, resolution: $resolution) {
            id
            title
            resolution
            resolves
            closes
            estimates {
                id
                brierScore
                probabilities {
                    id
                    outcome {
                        id
                        text
                        correct
                    }
                }
            }
        }
    }
`);

export const ResolveForecastDialog: FC<{
    forecast: Forecast,
    open: boolean,
    handleClose: () => void
}> = ({
    forecast,
    open,
    handleClose
}) => {
    const [outcome, setOutcome] = useState("");

    const handleChange = (event: ChangeEvent, value: string) => {
        setOutcome(value);
    }

    const OUTCOME_NONE = "__none";

    const resolveForecast = () => {
        let queryVars: ResolveForecastMutationVariables;
        if (outcome === OUTCOME_NONE) {
            queryVars = {
                forecastId: forecast.id,
                resolution: Resolution.NotApplicable,
            };
        } else {
            queryVars = {
                forecastId: forecast.id,
                correctOutcomeId: outcome,
            }
        }
        void resolveForecastMutation({variables: queryVars});
    }


    const [resolveForecastMutation, {error}] = useMutation<
        ResolveForecastMutation,
        ResolveForecastMutationVariables>(
            RESOLVE_FORECAST,
        {
            onCompleted: _ => handleClose(),
            // it seems to be necessary to call onError so that the
            // 'error' variable works correctly.
            // eslint-disable-next-line @typescript-eslint/no-empty-function
            onError: _ => {},
            },
        );

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
    }

    return (
        <Dialog
            open={open}
            onClose={(event, reason) => handleClose()}
            PaperComponent={PaperComponent}
            aria-labelledby="draggable-dialog-title"
        >
            <DialogTitle id="draggable-dialog-title" style={{cursor: "grab"}}>
                {forecast.title}
            </DialogTitle>
            <DialogContent>
                {messageBox}
                <FormControl>
                    <FormLabel id="demo-controlled-radio-buttons-group">What is the correct outcome?</FormLabel>
                    <RadioGroup
                        aria-labelledby="demo-controlled-radio-buttons-group"
                        name="controlled-radio-buttons-group"
                        value={outcome}
                        onChange={handleChange}
                    >
                        {
                            forecast.estimates[0]?.probabilities.map(p => {
                                return p ? <FormControlLabel key={p.outcome.id} value={p.outcome.id} control={<Radio />} label={p.outcome.text} />:<></>
                            })
                        }
                        <FormControlLabel value={OUTCOME_NONE} control={<Radio />} label="No correct outcome (resolve forecast as N/A)" style={{marginTop: "10px"}} />
                    </RadioGroup>
                </FormControl>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose} color="secondary">
                    Cancel
                </Button>
                <Button onClick={resolveForecast} color="primary">
                    Save
                </Button>
            </DialogActions>
        </Dialog>
    );
}
