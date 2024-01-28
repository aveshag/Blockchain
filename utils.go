package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func getUUID() string {
	id := uuid.New()
	return id.String()
}

func getTimestamp() int64 {
	t := time.Now()
	return t.Unix()
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IntToHex(data int64) []byte {
	hex_data := fmt.Sprintf("%x", data)
	return []byte(hex_data)
}
