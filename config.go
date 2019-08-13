package dumpin

import "github.com/pkg/errors"

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

func (c Config) Check() error {
	if c.host == "" ||
		c.port == "" ||
		c.user == "" ||
		// c.Password == "" ||
		c.database == "" {
		return errors.New("Missing parameter")
	}
	return nil
}
