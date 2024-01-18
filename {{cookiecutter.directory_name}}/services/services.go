package services

import (
	log "github.com/sirupsen/logrus"
)

type FormData = map[string][]string
type ValidationErrors = map[string]string

// enum for updated/no-update-but-no-error/erro when returning?

func parseForm(formData FormData, key string) (string, bool) {
	ret, ok := formData[key]
	if !ok || ret[0] == "" {
		log.Trace(formData)
		return "", false
	}
	return ret[0], true
}
