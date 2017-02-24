package helper

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadResource(filePath string) (*resource.Resource, string, error) {
	data, json, err := LoadData(filePath)
	if err != nil {
		return nil, "", err
	}
	return resource.NewResourceFromMap(data), json, nil
}

func LoadSchema(filePath string) (*resource.Schema, string, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	schema := &resource.Schema{}
	err = json.Unmarshal(fileBytes, &schema)
	if err != nil {
		return nil, "", err
	}

	return schema, string(fileBytes), nil
}

func LoadData(filePath string) (map[string]interface{}, string, error) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, "", err
	}

	return data, string(fileBytes), nil
}
