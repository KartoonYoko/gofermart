package pgsql

type Config struct {
	// Host             string
	// Port             uint16
	// ConnectTimeout   time.Duration
	// QueryTimeout     time.Duration
	// Username         string
	// Password         string
	// DBName           string
	ConnectionString string
	// MigrationVersion int64
}

func (c Config) connectionString() string {
	return c.ConnectionString
	// return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DBName)
}
