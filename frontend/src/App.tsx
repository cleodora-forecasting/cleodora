import {ForecastList} from "./ForecastList";
import {AddForecast} from "./AddForecast";

const App = () => {
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

export default App;