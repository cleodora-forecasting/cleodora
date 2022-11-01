// Import everything needed to use the `useQuery` hook
import { useQuery, useMutation, gql } from '@apollo/client';

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
                  data.forecasts.map(({ id, summary, description, created, closes, resolves, resolution }) => (
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

function AddForecast() {
    let summary, description, closes, resolves;
    const [addForecast, { data, loading, error }] = useMutation(ADD_FORECAST);

    if (loading) return 'Submitting...';
    if (error) return `Submission error! ${error.message}`;

    return (
        <div>
            <h3>Add Forecast</h3>
            <form
                onSubmit={e => {
                    e.preventDefault();
                    addForecast({
                        variables: {
                            input: {
                                summary: summary.value,
                                description: description.value,
                                closes: closes.value,
                                resolves: resolves.value,
                            },
                        },
                    });
                    summary.value = '';
                }}
            >
                <p>
                    <label>Summary</label>
                    <input
                        ref={node => {
                            summary = node;
                        }}
                    />
                </p>
                <p>
                    <label>Description</label>
                    <textarea
                        ref={node => {
                            description = node;
                        }}
                    />
                </p>
                <p>
                    <small>Format for the dates: 2022-12-01T09:00:00+01:00</small>
                    <br />
                    <label>Closes</label>
                    <input
                        ref={node => {
                            closes = node;
                        }}
                    />
                </p>
                <p>
                    <label>Resolves</label>
                    <input
                        ref={node => {
                            resolves = node;
                        }}
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