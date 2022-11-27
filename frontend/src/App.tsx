// Import everything needed to use the `useQuery` hook
import {ForecastList} from "./ForecastList";
import {AddForecast} from "./AddForecast";


export default function App() {
  return (
      <div>
        <h2>Cleodora</h2>
        <br/>
        <ForecastList />
        <br />
        <br />
        <AddForecast />
      </div>
  );
}
