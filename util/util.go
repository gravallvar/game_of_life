package util

import (
	"log"
	"math/rand"
	"os"
	"time"
)

// GetRandomCellState returns a random bool.
// The chance of the bool being true is a value convenient for generating a
// random matrix of cells that doesn't over or underpopulate the matrix.
func GetRandomCellState() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)%3 == 0
}

// LogToFile for debugging purposes.
func LogToFile(fileName string, msg string) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}
