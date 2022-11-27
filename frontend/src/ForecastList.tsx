import {FC} from "react";
import {gql, useQuery} from "@apollo/client";

export const GET_FORECASTS = gql`
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
`;

export const ForecastList: FC = () => {
    const {loading, error, data} = useQuery(GET_FORECASTS);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error :(</p>;

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
                        ({
                             id,
                             summary,
                             description,
                             created, closes,
                             resolves,
                             resolution
                         }:
                             {
                                 id: string,
                                 summary: string,
                                 description: string,
                                 created: Date,
                                 closes: Date,
                                 resolves: Date,
                                 resolution: string
                             }) => (
                            <tr key={id}>
                                <td>{id}</td>
                                <td>{summary}</td>
                                <td>{description}</td>
                                <td>{new Date(created).toLocaleString()}</td>
                                <td>{new Date(closes).toLocaleString()}</td>
                                <td>{new Date(resolves).toLocaleString()}</td>
                                <td>{resolution}</td>
                            </tr>
                        ))
                }
                </tbody>
            </table>
        </div>);
}