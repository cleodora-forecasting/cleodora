<img align="right" src="./logo_full.png">

# Cleodora Forecasting

Software to track personal forecasts/predictions and systematically improve at
making them.

Examples of such forecasts:

* Will "The Fabelmans" win "Best Picture" at the Oscars 2023?
* Will I get an A in my upcoming exam?
* Will there be nice weather on my birthday?
* Will the number of contributors for "Cleodora" be more than 3 at the end of 2022?

All information, including development guidelines, requirements and roadmap can
be found on the website https://cleodora.org .

# Run

Without embedded frontend:

```bash
go run .
```

With statically embedded frontend:

```bash
cd frontend
npm run build
cd ..
go run -tags frontend .
```

# Build

```bash
make build
```
