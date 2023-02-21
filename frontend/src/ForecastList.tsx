import React, {FC, useState} from "react";
import {useQuery} from "@apollo/client";
import {gql} from "./__generated__"
import {Paper, Table,
    TableBody, TableCell, TableContainer, TableHead, TableRow, Typography} from '@mui/material'
import {ForecastDetails} from "./ForecastDetails";
import {Forecast} from "./__generated__/graphql";
import {ResolutionChip} from "./ResolutionChip";

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
                created
                reason
                probabilities {
                    id
                    value
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

export const ForecastList: FC = () => {
    const {error, data} = useQuery(GET_FORECASTS);

    const [selectedForecast, setSelectedForecast] = useState<Forecast|null>(null);
    const [openDetails, setOpenDetails] = useState(false);

    let errorBox:JSX.Element;
    if (error) {
        const errMessages = new Array(error.message);
        if (error.message === 'NetworkError when attempting to fetch' +
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
                        let outcomes: JSX.Element | string = <></>;
                        if (Array.isArray(f.estimates)) {
                            // TODO the next line depends much
                            //  on what the API returns. It
                            //  probably needs to be adjusted.
                            const lastEstimate = f.estimates[f.estimates.length - 1];
                            if (lastEstimate != null && lastEstimate.probabilities.length > 0) {
                                outcomes = lastEstimate.probabilities.map(
                                    p => {
                                        if (p != null) {
                                            const content = p.outcome.text + ": " + p.value.toString() + "%";
                                            if (p.outcome.correct) {
                                                return <strong>{content}</strong>
                                            } else {
                                                return content
                                            }
                                        }
                                        return ""
                                    }
                                ).reduce((result, item) => <>{result} | {item}</>)
                            }
                        }
                        return (
                            <TableRow
                                key={f.id}
                                sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                hover
                                onClick={_ => {setSelectedForecast(f); setOpenDetails(true)}}
                                style={{cursor: 'pointer'}}
                            >
                                <TableCell
                                    component="th"
                                    scope="row"
                                >
                                    {f.title}
                                </TableCell>
                                <TableCell
                                    align="right">{new Date(f.resolves as string).toLocaleString()}</TableCell>
                                <TableCell
                                    align="right"><ResolutionChip resolution={f.resolution} /></TableCell>
                                <TableCell align="right">{outcomes}</TableCell>
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
            {selectedForecast ?
                <ForecastDetails
                    forecast={selectedForecast}
                    open={openDetails}
                    handleClose={() => setOpenDetails(false)}
                />
                : <></>
            }
        </div>);
}
