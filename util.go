package descriptor

import (
	"log"
	"strings"
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func readUrlOrFile(logger *log.Logger, location string) (base string, content []byte, err error) {
	if strings.Index(location, "http") == 0 {
		logger.Println("Loading URL", location)

		_, err = url.Parse(location)
		if err != nil {
			return
		}

		var response *http.Response
		response, err = http.Get(location)
		if err != nil {
			return
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			err = fmt.Errorf("HTTP Error getting the environment descriptor , error code %d", response.StatusCode)
			return
		}

		content, err = ioutil.ReadAll(response.Body)

		i := strings.LastIndex(location, "/")
		base = location[0 : i+1]
	} else {
		logger.Println("Loading file", location)
		var file *os.File
		file, err = os.Open(location)
		if err != nil {
			return
		}
		defer file.Close()
		content, err = ioutil.ReadAll(file)
		if err != nil {
			return
		}
		var absLocation string
		absLocation, err = filepath.Abs(location)
		if err != nil {
			return
		}
		base = filepath.Dir(absLocation) + string(filepath.Separator)
	}
	return
}
