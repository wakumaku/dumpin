package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wakumaku/dumpin/v3"
)

func main() {

	var sqlFile string

	var (
		host     string
		port     string
		user     string
		password string
		database string
	)

	flag.StringVar(&host, "h", "", "Hostname, eg.: my.db.com")
	flag.StringVar(&port, "P", "", "Port, eg.: 3306")
	flag.StringVar(&user, "u", "", "Username, eg.: root")
	flag.StringVar(&password, "p", "", "Password, eg.: root")
	flag.StringVar(&database, "d", "", "Database, eg.: my_database")

	flag.StringVar(&sqlFile, "file", "", "File, eg.: /full/path/to/file.sql")

	flag.Parse()

	config := dumpin.NewConfig(host, port, user, password, database)
	if err := config.Check(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dumpin, err := dumpin.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	output, err := dumpin.ExecuteFile(sqlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	fmt.Println(output)
}
