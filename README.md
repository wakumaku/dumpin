# Dumpin

## Dumps SQL files in!

WIP!

Load dumped SQL files, just for tests.

E.g.:

```
func loadDB() error {

	dumpCfg := dumpin.NewConfig("localhost", "3306", "root", "root", "mydatabase")
	dumpin, err := dumpin.New(dumpin.OsAUTODETECT, dumpin.EngMYSQL, dumpCfg)
	if err != nil {
		return err
	}

    output, err := dumpin.ExecuteFile("/path/to/dump.sql")
	if err != nil {
		return err
	}

    fmt.Println("Dump loaded!")
    fmt.Println(output)
}
```
