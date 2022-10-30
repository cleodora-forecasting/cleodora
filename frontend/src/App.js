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

  return data.forecasts.map(({ id, summary, description, created, closes }) => (
      <div key={id}>
        <h3>{summary}</h3>
        <p>{description}</p>
      </div>
  ));
}

export default function App() {
  return (
      <div>
        <h2>My first Apollo app 🚀</h2>
        <br/>
        <DisplayForecasts />
      </div>
  );
}