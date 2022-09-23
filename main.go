package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// version string
var version string

// helper function to return version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("auth-proxy-server git=%s go=%s date=%s", version, goVersion, tstamp)
}

// main function
func main() {
	var jsonFile string
	flag.StringVar(&jsonFile, "json", "", "json file")
	var yamlFile string
	flag.StringVar(&yamlFile, "yaml", "", "yaml file")
	var version bool
	flag.BoolVar(&version, "version", false, "print version information")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)
	}
	err := convert(jsonFile, yamlFile)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
