import {FC} from "react";
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
            outcomes {
                id
                text
            }
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
    const {loading, error, data} = useQuery(GET_FORECASTS);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error :-(</p>;
    if (data === undefined) return <p>Error :-(</p>
    return (
        <div>
            <Typography variant="h5">Forecasts</Typography>

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
                        {
                            data.forecasts.map(
                                f => {
                                    let estimates = ""
                                    if (Array.isArray(f.estimates)) {
                                        if (f.estimates[0] != null) { // TODO use the latest
                                            estimates = f.estimates[0].probabilities.map(
                                                p => {
                                                    if (p != null) {
                                                        return p.outcome.text + ": " + p.value.toString() + "%"
                                                    }
                                                    return ""
                                                }
                                            ).join(" | ")
                                        }
                                    }
                                    return (
                                        <TableRow
                                            key={f.id}
                                            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                        >
                                            <TableCell component="th" scope="row">{f.title}</TableCell>
                                            <TableCell align="right">{new Date(f.created as string).toLocaleString()}</TableCell>
                                            <TableCell align="right">{new Date(f.closes as string).toLocaleString()}</TableCell>
                                            <TableCell align="right">{new Date(f.resolves as string).toLocaleString()}</TableCell>
                                            <TableCell align="right">{f.resolution}</TableCell>
                                            <TableCell align="right">{estimates}</TableCell>
                                        </TableRow>
                                    )
                                }
                            )
                        }
                    </TableBody>
                </Table>
            </TableContainer>
        </div>);
}
