// Import everything needed to use the `useQuery` hook
import { useQuery, gql } from '@apollo/client';

const GET_FORECASTS = gql`
  query GetForecasts {
    forecasts {
      id
      summary
      description
      created
      closes
    }
  }
`;

function DisplayForecasts() {
  const { loading, error, data } = useQuery(GET_FORECASTS);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;

  return (
      <table>
          <thead>
          <tr>
              <th>ID</th>
              <th>Summary</th>
              <th>Description</th>
              <th>Created</th>
              <th>Closes</th>
          </tr>
          </thead>
          <tbody>
          {
              data.forecasts.map(({ id, summary, description, created, closes }) => (
                  <tr key={id}>
                      <td>{id}</td>
                      <td>{summary}</td>
                      <td>{description}</td>
                      <td>{created}</td>
                      <td>{closes}</td>
                  </tr>
              ))
          }
          </tbody>
      </table>);
}

export default function App() {
  return (
      <div>
        <h2>Cleodora</h2>
        <br/>
        <DisplayForecasts />
      </div>
  );
}