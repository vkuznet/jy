package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

// TestDict provides test of dict convertion
func TestDict(t *testing.T) {
	rec := make(map[string]any)
	rec["foo"] = 1
	testConvert(t, rec)
}

// TestList provides test of list convertion
func TestList(t *testing.T) {
	records := []int{1, 2, 3}
	testConvert(t, records)
}

// TestListDict provides test of list of dicts convertion
func TestListDict(t *testing.T) {
	var records []map[string]any
	rec := make(map[string]any)
	rec["foo"] = 1
	records = append(records, rec)
	testConvert(t, records)
}

// testConvert provides test of convert function
func testConvert(t *testing.T, jdata any) {
	jtmpFile, err := os.CreateTemp(os.TempDir(), "*.json")
	data, err := json.Marshal(jdata)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("JSON data", string(data))
	jtmpFile.Write(data)
	jtmpFile.Close()

	ytmpFile, err := os.CreateTemp(os.TempDir(), "*.yaml")

	// defer removal of temp files
	defer os.Remove(jtmpFile.Name())
	defer os.Remove(ytmpFile.Name())

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
	fmt.Printf("YAML data\n%s\n", string(data))
	var ydata map[any]any
	err = yaml.Unmarshal(data, &ydata)
	if err != nil {
		// try to unmarshal to []map[string]any
		var records []map[any]any
		err = yaml.Unmarshal(data, &records)
		if err != nil {
			var records []any
			err = yaml.Unmarshal(data, &records)
			if err != nil {
				msg := "unable to unmarshal yaml data to map[string]any or []map[string]any or []any"
				t.Error(wrapError(err, msg))
			} else {
				// test converted data
				srecords := fmt.Sprintf("%v", records)
				input := fmt.Sprintf("%v", jdata)
				if srecords != input {
					t.Errorf("input data '%v' differ from obtained data '%v", jdata, records)
				}
			}
		} else {
			// test converted data
			srecords := fmt.Sprintf("%v", records)
			input := fmt.Sprintf("%v", jdata)
			if srecords != input {
				t.Errorf("input data '%v' differ from obtained data '%v", jdata, records)
			}
		}
	}

	// test new json map
	jmap := convertMap(ydata)
	for k, v := range jmap {
		fmt.Println("json map, key", k, "value", v)
		switch m := jdata.(type) {
		case map[string]any:
			testMap(t, m, k, v)
		}
	}

	dumpFile(yamlFile)

	// convert yaml to json file
	newjFile, err := os.CreateTemp(os.TempDir(), "*.json")
	defer os.Remove(newjFile.Name())
	err = convert(newjFile.Name(), yamlFile)
	if err != nil {
		dumpFile(yamlFile)
		dumpFile(newjFile.Name())
		t.Error(err)
	}
	data, err = ioutil.ReadAll(newjFile)
	if err != nil {
		t.Error(err)
	}
	var rec map[string]any
	err = json.Unmarshal(data, &rec)
	if err != nil {
		// try to unmarshal to []map[string]any
		var records []map[any]any
		err = json.Unmarshal(data, &records)
		if err != nil {
			var records []any
			err = json.Unmarshal(data, &records)
			if err != nil {
				msg := "unable to unmarshal to map[string]any or []map[string]any or []any"
				t.Error(wrapError(err, msg))
			} else {
				// test converted data
				srecords := fmt.Sprintf("%v", records)
				input := fmt.Sprintf("%v", jdata)
				if srecords != input {
					t.Errorf("input data '%v' differ from obtained data '%v", jdata, records)
				}
			}
		} else {
			// test converted []map[string]any data
			srecords := fmt.Sprintf("%v", records)
			input := fmt.Sprintf("%v", jdata)
			if srecords != input {
				t.Errorf("input data '%v' differ from obtained data '%v", jdata, records)
			}
		}
	}
	for k, v := range rec {
		fmt.Println("json map, key", k, "value", v)
		switch m := jdata.(type) {
		case map[string]any:
			testMap(t, m, k, v)
		}
	}
}

// helper function to test k/v in a map
func testMap(t *testing.T, jdata map[string]any, k string, v any) {
	val, ok := jdata[k]
	if ok {
		sval := fmt.Sprintf("%v", val)
		sv := fmt.Sprintf("%v", v)
		if sval != sv {
			t.Errorf("Mismatch of data, for key=%s, value=%v of type %T expect=%v of type %T", k, val, val, v, v)
		}
	} else {
		t.Errorf("Unable to find map key")
	}
}

// helper function to dump file content
func dumpFile(fname string) {
	if fileExist(fname) {
		file, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fmt.Println("content of file", fname)
		if data, err := io.ReadAll(file); err == nil {
			fmt.Println(string(data))
		}
	}
}
