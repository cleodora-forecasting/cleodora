import {ForecastList} from "./ForecastList";
import {AddForecast} from "./AddForecast";
import {AppBar, CssBaseline, Toolbar, Typography} from "@mui/material";
import logo from './logo.png'

const App = () => {
  return (
      <>
        <CssBaseline />
        <AppBar position="relative">
          <Toolbar>
            <img width="30px" src={logo} alt="Cleodora logo" />
            <Typography variant="h6">Cleodora</Typography>
          </Toolbar>
        </AppBar>
        <br/>
        <ForecastList />
        <br />
        <br />
        <AddForecast />
      </>
  );
}

export default App;