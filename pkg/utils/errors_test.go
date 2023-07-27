package utils

import (
	"log"
	"testing"
)

func TestErrors(t *testing.T) {
	log.Println(BadRequest())
	log.Println(Unauthorized())
	log.Println(Forbidden())
}
