package microservice

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)



func LogEnvironment(prefix string) {
 	log.Debug("environment variables: ")
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, prefix) {
			log.Debug(e)
		}
	}
}

