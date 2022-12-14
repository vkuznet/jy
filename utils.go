package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/exp/errors"
	yaml "gopkg.in/yaml.v2"
)

// fileExist checks if file exists
func fileExist(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if finfo.Size() == 0 {
		return false
	}
	return !errors.Is(err, os.ErrNotExist)
}

// helper function to convert json to yaml and vice versa
func convert(jsonFile, yamlFile string) error {
	if fileExist(jsonFile) && fileExist(yamlFile) {
		msg := fmt.Sprintf("Both input files %s and %s exist, please provide one non-existing input to perform conversion", jsonFile, yamlFile)
		return errors.New(msg)
	} else if fileExist(jsonFile) {
		fmt.Printf("Convert %s to %s\n", jsonFile, yamlFile)
		return convertJson2Yaml(jsonFile, yamlFile)
	} else if fileExist(yamlFile) {
		fmt.Printf("Convert %s to %s\n", yamlFile, jsonFile)
		return convertYaml2Json(yamlFile, jsonFile)
	}
	msg := fmt.Sprintf("Neither %s or %s exist\n", jsonFile, yamlFile)
	return errors.New(msg)
}

// helper function to read file content
func readFile(fname string) ([]byte, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	return data, err
}

// helper function to write file with given data
func writeFile(fname string, data []byte) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(data)
	return nil
}

// helper function to convert yaml map to json map interface
func convertMap(m map[any]any) map[string]any {
	res := map[string]any{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[any]any:
			res[fmt.Sprint(k)] = convertMap(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

// helper function to wrap error
func wrapError(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s, %s", err, msg))
}

// helper function to convert json to yaml file
func convertJson2Yaml(jsonFile, yamlFile string) error {
	data, err := readFile(jsonFile)
	if err != nil {
		return err
	}
	var record map[string]any
	err = json.Unmarshal(data, &record)
	if err != nil {
		// try to load the list of map records
		var records []map[string]any
		if err := json.Unmarshal(data, &records); err != nil {
			// try to load list of basic data-types, e.g. list of ints or strings
			var records []any
			if err := json.Unmarshal(data, &records); err != nil {
				msg := fmt.Sprintf("record is not []any or []map[string]any or map[string]any but %s of type %T", string(data), data)
				return wrapError(err, msg)
			}
			data, err = yaml.Marshal(records)
			if err != nil {
				return wrapError(err, "unable to marshal []any")
			}
		} else {
			data, err = yaml.Marshal(records)
			if err != nil {
				return wrapError(err, "unable to marshal []map[string]any")
			}
		}
	} else {
		data, err = yaml.Marshal(record)
		if err != nil {
			return err
		}
	}
	return writeFile(yamlFile, data)
}

// helper function to convert yaml to json file
func convertYaml2Json(yamlFile, jsonFile string) error {
	data, err := readFile(yamlFile)
	if err != nil {
		return err
	}
	var record map[any]any
	err = yaml.Unmarshal(data, &record)
	if err != nil {
		// try to load list of records
		var records []map[any]any
		err = yaml.Unmarshal(data, &records)
		if err != nil {
			var records []any
			err = yaml.Unmarshal(data, &records)
			if err != nil {
				return wrapError(err, "input data is not map[any]any or []map[any]any or []any")
			}
			data, err = json.Marshal(records)
			if err != nil {
				return err
			}
			return writeFile(jsonFile, data)
		}
		var out []map[string]any
		for _, r := range records {
			out = append(out, convertMap(r))
		}
		data, err = json.Marshal(out)
		if err != nil {
			return err
		}
		return writeFile(jsonFile, data)
	}
	jsonData := convertMap(record)
	data, err = json.Marshal(jsonData)
	if err != nil {
		return err
	}
	return writeFile(jsonFile, data)
}
