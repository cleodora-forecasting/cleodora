import {FC} from "react";
import {Forecast, Resolution} from "./__generated__/graphql";
import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Paper
} from "@mui/material";
import Draggable from 'react-draggable';
import {ResolutionChip} from "./ResolutionChip";

const PaperComponent: FC = (props) => {
    return (
        <Draggable handle="#draggable-dialog-title">
            <Paper {...props} />
        </Draggable>
    )
}

export const ForecastDetails: FC<{forecast: Forecast, open: boolean, handleClose: () => void}> = ({forecast, open, handleClose}) => {

    const getCorrectOutcome = () => {
        let outcome = <></>;
        if (forecast &&
            forecast.resolution === Resolution.Resolved &&
            forecast.estimates &&
            forecast.estimates.length > 0 &&
            forecast.estimates[0]
        ) {
            forecast.estimates[0].probabilities.forEach(function (p) {
                if (p && p.outcome.correct) {
                    outcome = <strong>({p.outcome.text})</strong>;
                    return
                }
            });
        }
        return outcome;
    };

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
            <ResolutionChip resolution={forecast.resolution} /> {getCorrectOutcome()}
            <ul>
                <li><strong>ID: </strong>{forecast.id}</li>
                <li><strong>Created: </strong>{new Date(forecast.created as string).toLocaleString()}</li>
                <li><strong>Resolves: </strong>{new Date(forecast.resolves as string).toLocaleString()}</li>
                <li><strong>Closes: </strong>{forecast.closes ? new Date(forecast.closes as string).toLocaleString() : ''}</li>
            </ul>
            <p>{forecast.description}</p>
            <h3>Estimates</h3>
            <>
                {forecast.estimates.map(estimate => {
                    if (estimate === null) {
                        return null
                    }
                    return (
                        <div key={estimate.id}>
                            <h4>{new Date(estimate.created as string).toLocaleString()}</h4>
                            <p><strong>Brier Score: </strong>{estimate.brierScore}</p>
                            <p><strong>Reason: </strong>{estimate.reason}</p>
                            <ul>
                                {estimate.probabilities.map(probability => {
                                    if (probability === null) {
                                        return null;
                                    }
                                    return (
                                        <li key={probability.id}>
                                            <strong>{probability.outcome.text}: </strong>{probability.value}%
                                        </li>
                                    )
                                })}
                            </ul>
                        </div>
                    )
                })}
            </>
        </DialogContent>
        <DialogActions>
            <Button onClick={handleClose} color="primary">
                Close
            </Button>
        </DialogActions>
    </Dialog>
    );
}
