package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

// TestYaml2Json provides test of GET method for our service
func TestYaml2Json(t *testing.T) {
	jtmpFile, err := os.CreateTemp(os.TempDir(), "json-data")
	jdata := make(map[string]interface{})
	jdata["foo"] = 1
	data, err := json.Marshal(jdata)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("JSON data", string(data))
	jtmpFile.Write(data)
	jtmpFile.Close()

	ytmpFile, err := os.CreateTemp(os.TempDir(), "yaml-data")
	defer ytmpFile.Close()

	jsonFile := jtmpFile.Name()
	yamlFile := ytmpFile.Name()
	err = convert(jsonFile, yamlFile)
	if err != nil {
		t.Error(err)
	}

	// read back yaml data
	data, err = ioutil.ReadAll(ytmpFile)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("YAML data", string(data))
	var ydata map[interface{}]interface{}
	err = yaml.Unmarshal(data, &ydata)
	if err != nil {
		t.Error("Unable to unmarshal yaml data")
	}

	// test new json map
	jmap := convertYaml(ydata)
	for k, v := range jmap {
		fmt.Println("json map, key", k, "value", v)
		val, ok := jdata[k]
		if ok {
			if val != v {
				t.Errorf("Mismatch of data")
			}
		} else {
			t.Errorf("Unable to find map key")
		}
	}

	// convert yaml to json file
	newjFile, err := os.CreateTemp(os.TempDir(), "new-json-data")
	defer newjFile.Close()
	err = convert(yamlFile, newjFile.Name())
	if err != nil {
		t.Error(err)
	}
	data, err = ioutil.ReadAll(newjFile)
	if err != nil {
		t.Error(err)
	}
	var rec map[string]interface{}
	err = json.Unmarshal(data, &rec)
	if err != nil {
		t.Error(err)
	}
	for k, v := range rec {
		fmt.Println("json map, key", k, "value", v)
		val, ok := jdata[k]
		if ok {
			if val != v {
				t.Errorf("Mismatch of data")
			}
		} else {
			t.Errorf("Unable to find map key")
		}
	}
}
