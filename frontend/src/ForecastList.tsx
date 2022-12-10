import {FC} from "react";
import {useQuery} from "@apollo/client";
import {gql} from "./__generated__"
import { Typography } from '@mui/material'

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
            <Typography variant="h3">Forecasts</Typography>
            <table>
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Title</th>
                    <th>Description</th>
                    <th>Created</th>
                    <th>Closes</th>
                    <th>Resolves</th>
                    <th>Resolution</th>
                    <th>Latest Estimate</th>
                </tr>
                </thead>
                <tbody>
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
                                <tr key={f.id}>
                                    <td>{f.id}</td>
                                    <td>{f.title}</td>
                                    <td>{f.description}</td>
                                    <td>{new Date(f.created as string).toLocaleString()}</td>
                                    <td>{new Date(f.closes as string).toLocaleString()}</td>
                                    <td>{new Date(f.resolves as string).toLocaleString()}</td>
                                    <td>{f.resolution}</td>
                                    <td>{estimates}</td>
                                </tr>
                            )
                        }
                    )
                }
                </tbody>
            </table>
        </div>);
}
