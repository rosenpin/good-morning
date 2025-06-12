package config

// Config holds the application configuration
type Config struct {
	API            API
	Image          Image
	Search         Search
	MaxDailyReload int
}
