package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"golang.org/x/exp/errors"
	yaml "gopkg.in/yaml.v2"
)

// version string
var version string

// helper function to return version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("auth-proxy-server git=%s go=%s date=%s", version, goVersion, tstamp)
}

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

// fileExist checks if file exists
func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func convert(jsonFile, yamlFile string) error {
	if fileExist(jsonFile) {
		fmt.Printf("Convert %s to %s\n", jsonFile, yamlFile)
		return convertJson2Yaml(jsonFile, yamlFile)
	} else if fileExist(yamlFile) {
		fmt.Printf("Convert %s to %s\n", yamlFile, jsonFile)
		return convertYaml2Json(yamlFile, jsonFile)
	}
	msg := fmt.Sprintf("Neither %s or %s exist\n", jsonFile, yamlFile)
	return errors.New(msg)
}

func readFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	return data, err
}

func convertJson2Yaml(jsonFile, yamlFile string) error {
	data, err := readFile(jsonFile)
	if err != nil {
		return err
	}
	var record map[string]interface{}
	err = json.Unmarshal(data, &record)
	data, err = yaml.Marshal(record)
	file, err := os.Create(yamlFile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}

func convertYaml2Json(yamlFile, jsonFile string) error {
	data, err := readFile(yamlFile)
	if err != nil {
		return err
	}
	var record map[string]interface{}
	err = yaml.Unmarshal(data, &record)
	data, err = json.Marshal(record)
	file, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}
