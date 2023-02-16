import React, {FC} from "react";
import {useQuery} from "@apollo/client";
import {gql} from "./__generated__"
import {Paper, Table,
    TableBody, TableCell, TableContainer, TableHead, TableRow, Typography} from '@mui/material'

export const GET_FORECASTS = gql(`
    query GetForecasts {
        forecasts {
            id
            title
            description
            created
            closes
            resolves
            resolution
            estimates {
                id
                probabilities {
                    id
                    value
                    outcome {
                        id
                        text
                    }
                }
            }
        }
    }
`);

export const ForecastList: FC = () => {
    const {error, data} = useQuery(GET_FORECASTS);

    let errorBox:JSX.Element;
    if (error) {
        const errMessages = new Array(error.message);
        if (error.message == 'NetworkError when attempting to fetch' +
            ' resource.') {
            errMessages.push('Make sure the API is running and reachable.');
        }
        errorBox =
            <p style={{color:"red"}}>
                <strong>{error.name}: </strong>
                {errMessages.join(' ')}
            </p>
    } else {
        errorBox = <></>
    }

    let tableBody:JSX.Element;
    if (data === undefined) {
        tableBody = <></>;
    } else {
        tableBody = <>
            {
                data.forecasts.map(
                    f => {
                        let estimates = ""
                        if (Array.isArray(f.estimates)) {
                            // TODO the next line depends much
                            //  on what the API returns. It
                            //  probably needs to be adjusted.
                            const lastEstimate = f.estimates[f.estimates.length - 1];
                            if (lastEstimate != null) {
                                estimates = lastEstimate.probabilities.map(
                                    p => {
                                        if (p != null) {
                                            return p.outcome.text + ": " + p.value.toString() + "%"
                                        }
                                        return ""
                                    }
                                ).join(" | ")
                            }
                        }
                        let closes = ""
                        if (f.closes) {
                            closes = new Date(f.closes as string).toLocaleString()
                        }
                        return (
                            <TableRow
                                key={f.id}
                                sx={{'&:last-child td, &:last-child th': {border: 0}}}
                            >
                                <TableCell component="th"
                                           scope="row">{f.title}</TableCell>
                                <TableCell
                                    align="right">{new Date(f.created as string).toLocaleString()}</TableCell>
                                <TableCell align="right">{closes}</TableCell>
                                <TableCell
                                    align="right">{new Date(f.resolves as string).toLocaleString()}</TableCell>
                                <TableCell
                                    align="right">{f.resolution}</TableCell>
                                <TableCell align="right">{estimates}</TableCell>
                            </TableRow>
                        )
                    }
                )
            }
        </>
    }
    return (
        <div>
            <Typography variant="h5">Forecasts</Typography>
            {errorBox}
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="forecasts">
                    <TableHead>
                        <TableRow>
                            <TableCell>Title</TableCell>
                            <TableCell align="right">Created</TableCell>
                            <TableCell align="right">Closes</TableCell>
                            <TableCell align="right">Resolves</TableCell>
                            <TableCell align="right">Resolution</TableCell>
                            <TableCell align="right">Estimate</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {tableBody}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>);
}
