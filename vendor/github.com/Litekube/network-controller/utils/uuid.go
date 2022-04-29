package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GetUniqueToken() string {
	uuid := uuid.NewV4().String()
	uuid = strings.ReplaceAll(uuid, "-", "")[:16]
	logger.Infof("gen uuid token: %+v", uuid)
	return uuid
}
