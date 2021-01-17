package main

import (
	"envtool/utils"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
}

func main() {
	fmt.Println("Test..")
	utils.ListContainers()

}
