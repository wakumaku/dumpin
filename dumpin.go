package dumpin

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"runtime"
)

// Platform to use
type Platform string

// Platforms available to use
const (
	OsLINUX Platform = "linux"
	OsOSX   Platform = "osx"
)

// OsAUTODETECT tries to autodetect the OS
const OsAUTODETECT Platform = "auto"

// OsNOTSUPPORTED is returned if no valid OS found
const OsNOTSUPPORTED Platform = "notsupported"

// Engine types to support
type Engine string

// Engines list
const (
	EngMYSQL Engine = "mysql"
)

const binariesPath = "bin/"

var (
	binaries = map[Engine]map[Platform]string{
		EngMYSQL: map[Platform]string{
			OsLINUX: "mysql",
			OsOSX:   "mysql",
		},
	}
)

var (
	errPlatformNotSupported = errors.New("Platform not suported")
	errFileCanNotBeRead     = errors.New("SQL File cannot be read")
)

// Dumpin main struct
type Dumpin struct {
	platform Platform
	engine   Engine
	config   Config
}

// New Dumpin
func New(platform Platform, engine Engine, config Config) (*Dumpin, error) {

	if platform == OsAUTODETECT {
		platformDetected, err := determinePlatform()
		if err != nil {
			return nil, err
		}
		platform = platformDetected
	}

	return &Dumpin{
		platform: platform,
		engine:   engine,
		config:   config,
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

func (m *Dumpin) getExecutableFile() string {
	return fmt.Sprintf("%s%s%s/%s/%s", packagePath(), binariesPath, m.engine, m.platform, binaries[m.engine][m.platform])
}

// ExecuteFile ...
func (m *Dumpin) ExecuteFile(sqlFilePath string, customArgs ...string) (string, error) {

	sql, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return "", errFileCanNotBeRead
	}

	return m.Execute(sql, customArgs...)
}

// Execute ...
func (m *Dumpin) Execute(sql []byte, customArgs ...string) (string, error) {

	var outbuf, errbuf bytes.Buffer

	customArgs = append(m.buildArgs(), customArgs...)

	cmd := exec.Command(m.getExecutableFile(), customArgs...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", errors.New("Piping in")
	}

	stdin.Write(sql)

	if err := cmd.Start(); err != nil {
		return outbuf.String(), errors.New(errbuf.String())
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return outbuf.String(), errors.New(errbuf.String())
	}

	return outbuf.String(), nil
}

func determinePlatform() (Platform, error) {
	switch runtime.GOOS {
	case "linux":
		return OsLINUX, nil
	case "darwin":
		return OsOSX, nil
	}

	return OsNOTSUPPORTED, errPlatformNotSupported
}

func packagePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	return path.Dir(filename) + "/"
}
