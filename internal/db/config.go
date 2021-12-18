package db

// Config holds database configuration.
type Config struct {
	URL  string `koanf:"url"`
	Name string `koanf:"name"`
}
