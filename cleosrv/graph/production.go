//go:build production

package graph

// AddDummyData adds dummy data (e.g. example forecasts) to the database.
// It's not used in production.
func (r *Resolver) AddDummyData() error {
	return nil
}
