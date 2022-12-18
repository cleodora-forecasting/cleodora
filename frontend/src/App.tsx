import {ForecastList} from "./ForecastList";
import {AddForecast} from "./AddForecast";
import {Footer} from "./Footer";
import {AppBar, CssBaseline, Toolbar, Typography} from "@mui/material";
import { ReactComponent as Logo } from './logo.svg';

const App = () => {
  return (
      <>
        <CssBaseline />
        <AppBar position="relative">
          <Toolbar>
            <Logo />
            <Typography variant="h6">Cleodora</Typography>
          </Toolbar>
        </AppBar>
        <br/>
        <ForecastList />
        <br />
        <br />
        <AddForecast />
        <Footer />
      </>
  );
}

export default App;
