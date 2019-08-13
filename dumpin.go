package dumpin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ErrMySqlClientNotFound is raised when the mysql client is not found in path to be executed
var ErrMySqlClientNotFound = errors.New("mysql client not found in PATH, please, install mysql client")

var ErrFileCanNotBeRead = errors.New("SQL File cannot be read")
var ErrCannotPipeInContent = errors.New("OS error, cannot pipe in the SQL file")

// Dumpin main struct
type Dumpin struct {
	config Config
}

// New Dumpin
func New(config Config) (*Dumpin, error) {

	if _, err := isMySqlCliAvailable(); err != nil {
		return nil, err
	}

	return &Dumpin{
		config: config,
	}, nil
}

func (m *Dumpin) hostSwitch() string {
	return fmt.Sprintf("-h%s", m.config.host)
}

func (m *Dumpin) portSwitch() string {
	return fmt.Sprintf("-p%s", m.config.host)
}

func (m *Dumpin) userSwitch() string {
	return fmt.Sprintf("-u%s", m.config.user)
}

func (m *Dumpin) passwordSwitch() string {
	return fmt.Sprintf("-p%s", m.config.password)
}

func (m *Dumpin) databaseSwitch() string {
	return fmt.Sprintf("%s", m.config.database)
}

func (m *Dumpin) buildArgs() []string {
	return []string{
		"--protocol",
		"tcp",
		m.hostSwitch(),
		m.userSwitch(),
		m.passwordSwitch(),
		m.databaseSwitch(),
	}
}

// ExecuteFile ...
func (m *Dumpin) ExecuteFile(sqlFilePath string, customArgs ...string) (string, error) {

	sql, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return "", ErrFileCanNotBeRead
	}

	return m.Execute(sql, customArgs...)
}

// Execute the sql file
func (m *Dumpin) Execute(sql []byte, customArgs ...string) (string, error) {

	var outbuf, errbuf bytes.Buffer

	customArgs = append(m.buildArgs(), customArgs...)

	cmd := exec.Command("/bin/sh", "-c", "mysql "+strings.Join(customArgs, " "))
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", ErrCannotPipeInContent
	}

	if _, err := stdin.Write(sql); err != nil {
		return "", errors.Wrap(err, "cannot pipe in the SQL")
	}

	if err := cmd.Start(); err != nil {
		return outbuf.String(), errors.New(errbuf.String())
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return outbuf.String(), errors.New(errbuf.String())
	}

	return outbuf.String(), nil
}

func isMySqlCliAvailable() (bool, error) {
	cmd := exec.Command("/bin/sh", "-c", "command -v mysql")
	if err := cmd.Run(); err != nil {
		return false, ErrMySqlClientNotFound
	}
	return true, nil
}
