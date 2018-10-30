package initialize

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadYAML(filePath string) (map[string]string, error) {
	suppliers := map[string]string{}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &suppliers)
	if err != nil {
		return nil, err
	}

	return suppliers, nil
}
