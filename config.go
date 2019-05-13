package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

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
	user, err := user.Current()
	if err != nil {
		return conf, fmt.Errorf("can't get user identity: %v", err)
	}
	// Default config location is XDG-CONFIG for Linux, BSD, ...
	prefsfile := filepath.Join(user.HomeDir, ".config", "hugoutil", "prefs.toml")
	switch runtime.GOOS {
	case "windows":
		// I think this is right for Windows, but I don't really use Windows
		prefsfile = filepath.Join(user.HomeDir, "AppData", "Local", "hugoutil", "prefs.toml")
	case "darwin":
		prefsfile = filepath.Join(user.HomeDir, "Library", "Preferences", "com.ath0.hugoutil", "prefs.toml")
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
