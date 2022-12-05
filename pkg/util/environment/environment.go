package environment

import (
	"log"
	"os"
)

// Application Environment name
const (
	Development = "dev"
	Test        = "test"
	E2E         = "e2e"
)

// IsDev returns APP_ENV in development mode
func IsDev() bool {
	log.Println("APP_ENV: ", os.Getenv("APP_ENV"))
	return os.Getenv("APP_ENV") == Development
}

// IsDev returns APP_ENV in test mode
func IsTest() bool {
	log.Println("APP_ENV: ", os.Getenv("APP_ENV"))
	return os.Getenv("APP_ENV") == Test
}

// IsDev returns APP_ENV in e2e mode
func IsE2E() bool {
	log.Println("APP_ENV: ", os.Getenv("APP_ENV"))
	return os.Getenv("APP_ENV") == E2E
}
