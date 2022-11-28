import {FC} from "react";
import {useQuery} from "@apollo/client";
import {gql} from "./__generated__"

export const GET_FORECASTS = gql(`
    query GetForecasts {
        forecasts {
            id
            summary
            description
            created
            closes
            resolves
            resolution
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
            <h3>Forecasts</h3>
            <table>
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Summary</th>
                    <th>Description</th>
                    <th>Created</th>
                    <th>Closes</th>
                    <th>Resolves</th>
                    <th>Resolution</th>
                </tr>
                </thead>
                <tbody>
                {
                    data.forecasts.map(
                        f => (
                            <tr key={f.id}>
                                <td>{f.id}</td>
                                <td>{f.summary}</td>
                                <td>{f.description}</td>
                                <td>{new Date(f.created as string).toLocaleString()}</td>
                                <td>{new Date(f.closes as string).toLocaleString()}</td>
                                <td>{new Date(f.resolves as string).toLocaleString()}</td>
                                <td>{f.resolution}</td>
                            </tr>
                        )
                    )
                }
                </tbody>
            </table>
        </div>);
}