package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config represents the configuration information used for access to an IBM Watson Cloud account.
type Config struct {
	APIKey string `toml:"api_key"`
	APIURL string `toml:"api_url"`
}

const sampleConfig = `# Sample configuration file for hugoutil, uncomment lines and set your API values.

[watson]
# api_key = "f755389ZmFrZSBjcmVkZW50aWFsc_wfd79855368006"
# api_url = "https://gateway.watsonplatform.net/natural-language-understanding/api"
`

func writeSampleConfig(filename string) {
	path := filepath.Dir(filename)
	err := os.MkdirAll(path, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create preferences directory %s: %v\n", path, err)
		return
	}
	err = ioutil.WriteFile(filename, []byte(sampleConfig), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Internal error creating preferences file %s: %v\n", filename, err)
		return
	}
	fmt.Printf("Created sample preferences file %s\n", filename)
}

func loadConfig() (Config, error) {
	conf := Config{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return conf, fmt.Errorf("can't get user home directory: %w", err)
	}
	confdir, err := os.UserConfigDir()
	if err != nil {
		return conf, fmt.Errorf("can't get user config directory: %w", err)
	}
	var prefsfile string
	switch runtime.GOOS {
	case "windows":
		// On Windows, %APPDATA%/Local is preferences local to the current machine.
		prefsfile = filepath.Join(confdir, "Local", "hugoutil", "prefs.toml")
	case "darwin":
		// On macOS, use Preferences rather than Application Support, i.e. don't use os.UserConfigDir.
		prefsfile = filepath.Join(homedir, "Library", "Preferences", "com.ath0.hugoutil", "prefs.toml")
	default:
		// Everything else is XDG spec so use a directory under os.UserConfigDir.
		prefsfile = filepath.Join(confdir, "hugoutil", "prefs.toml")
	}
	if *verbose {
		fmt.Printf("Checking for preferences in %s\n", prefsfile)
	}
	data, err := ioutil.ReadFile(prefsfile)
	if err != nil {
		writeSampleConfig(prefsfile)
		return conf, nil
	}
	_, err = toml.Decode(string(data), &conf)
	return conf, err
}
