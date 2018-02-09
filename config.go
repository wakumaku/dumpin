package dumpin

// Config ...
type Config struct {
	host     string
	port     string
	user     string
	password string
	database string
}

// NewConfig ...
func NewConfig(host, port, user, password, database string) Config {
	return Config{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
	}
}
