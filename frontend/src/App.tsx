// Import everything needed to use the `useQuery` hook
import { useQuery, useMutation, gql } from '@apollo/client';
import {useState} from "react";

const GET_FORECASTS = gql`
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

function DisplayForecasts() {
  const { loading, error, data } = useQuery(GET_FORECASTS);

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
                      ({id,
                           summary,
                           description,
                           created, closes,
                           resolves,
                           resolution }:
                           {id:string,
                               summary:string,
                               description:string,
                               created:Date,
                               closes:Date,
                               resolves:Date,
                               resolution:string }) => (
                      <tr key={id}>
                          <td>{id}</td>
                          <td>{summary}</td>
                          <td>{description}</td>
                          <td>{ new Date(created).toLocaleString() }</td>
                          <td>{ new Date(closes).toLocaleString() }</td>
                          <td>{ new Date(resolves).toLocaleString() }</td>
                          <td>{resolution}</td>
                      </tr>
                  ))
              }
              </tbody>
          </table>
      </div>);
}

const ADD_FORECAST = gql`
    mutation createForecast($input: NewForecast!) {
        createForecast(input: $input) {
            id
            summary
        }
    }
`;

const AddForecast = () => {
    const [summary, setSummary] = useState('');
    const [description, setDescription] = useState('');
    const [closes, setCloses] = useState(''); // TODO date
    const [resolves, setResolves] = useState(''); // TODO date
    const [addForecast, { error, data }] = useMutation(ADD_FORECAST, {
        refetchQueries: [
            {query: GET_FORECASTS}, // DocumentNode object parsed with gql
            'GetForecasts' // Query name
        ],
        variables: {
            input: {
                summary,
                description,
                closes,
                resolves,
            },
        },
    });

    return (
        <div>
            <h3>Add Forecast</h3>
            {error ? <p style={{color: "red"}}>Oh no! {error.message}</p> : null}
            {data && data.createForecast ?
                <p style={{color: "green"}}>
                    Saved "{data.createForecast.summary}" with ID {data.createForecast.id}.
                </p>
                : null}
            <form
                onSubmit={e => {
                    e.preventDefault();
                    addForecast().then(() => {
                        setSummary('');
                        setDescription('');
                        setCloses('');
                        setResolves('');
                    });
                }}
            >
                <p>
                    <label>Summary</label>
                    <input
                        name="summary"
                        value={summary}
                        onChange={e => setSummary(e.target.value)}
                    />
                </p>
                <p>
                    <label>Description</label>
                    <textarea
                        name="description"
                        value={description}
                        onChange={e => setDescription(e.target.value)}
                    />
                </p>
                <p>
                    <small>Format for the dates: 2022-12-01T09:00:00+01:00</small>
                    <br />
                    <label>Closes</label>
                    <input
                        name="closes"
                        value={closes}
                        onChange={e => setCloses(e.target.value)}
                    />
                </p>
                <p>
                    <label>Resolves</label>
                    <input
                        name="resolves"
                        value={resolves}
                        onChange={e => setResolves(e.target.value)}
                    />
                </p>
                <button type="submit">Add Forecast</button>
            </form>
        </div>
    );
}

export default function App() {
  return (
      <div>
        <h2>Cleodora</h2>
        <br/>
        <DisplayForecasts />
        <br />
        <br />
        <AddForecast />
      </div>
  );
}
