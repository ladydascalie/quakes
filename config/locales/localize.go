package locales

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const base = "config/locales/"

var (
	// Locale - the currently selected locale. Provided as part of the request URLs
	Locale = "en"
	// Strings unwraps the JSON translations.
	// It uses interface{} as values in case a value is included in the JSON which isn't a string. This way a type assertion can be made to fail, allowing us to bail and return an error.
	Strings   = make(map[string]interface{})
	available = []string{"en", "jp"}
)

// Load the locale translation, given the language key
func Load(locale string) error {
	if !Valid(locale) {
		locale = "en"
	}
	fpath := filepath.Join(base, locale+".json")
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	return json.Unmarshal(contents, &Strings)
}

// Valid checks if the provided locale key has an available translation
func Valid(locale string) bool {
	for _, loc := range available {
		if locale == loc {
			return true
		}
	}
	return false
}

// T is a tersely named function which fetches the translated key from the loaded locale file.
// It is made available in the templates FuncMap under the same name `T`
func T(msg string) string {
	str, ok := Strings[msg].(string)
	if !ok {
		return ""
	}
	return str
}
