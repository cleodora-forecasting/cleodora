import {FC} from "react";
import {Forecast}  from "./__generated__/graphql";
import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle, Paper
} from "@mui/material";
import Draggable from 'react-draggable';

const PaperComponent: FC = (props) => {
    return (
        <Draggable>
            <Paper {...props} />
        </Draggable>
    )
}

export const ForecastDetails: FC<{forecast: Forecast, open: boolean, handleClose: () => void}> = ({forecast, open, handleClose}) => {
    return (
    <Dialog
        open={open}
        onClose={handleClose}
        PaperComponent={PaperComponent}
        aria-labelledby="draggable-dialog-title"
    >
        <DialogTitle id="draggable-dialog-title">{forecast.title}</DialogTitle>
        <DialogContent>
            <ul>
                <li><strong>ID: </strong>{forecast.id}</li>
                <li><strong>Resolution: </strong>{forecast.resolution}</li>
                <li><strong>Created: </strong>{forecast.created}</li>
                <li><strong>Resolves: </strong>{forecast.resolves}</li>
                <li><strong>Closes: </strong>{forecast.closes}</li>
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
                            <h4>{estimate.created}</h4>
                            <p>{estimate.reason}</p>
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
