package goflow

import (
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xinzf/goflow.v2/definition"
	"io/ioutil"
	"os"
	"path/filepath"
)

var workflows map[string]definition.Workflow

func Start(xmlPath string) error {
	if pathExists(xmlPath) == false {
		return fmt.Errorf("Path: %s is not exists", xmlPath)
	}

	newMap := make(map[string]definition.Workflow)
	files, err := filepath.Glob(xmlPath + "/*.xml")
	if err != nil {
		return err
	}

	for _, f := range files {
		if !pathExists(f) {
			continue
		}

		d, err := loadFromFile(f)
		if err != nil {
			return err
		}

		newMap[d.ID] = d
	}

	workflows = newMap
	logrus.Debugln("装载 XML 完毕...")

	return nil
}

func loadFromFile(xmlFile string) (definition.Workflow, error) {
	file, err := os.Open(xmlFile)
	if err != nil {
		return definition.Workflow{}, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return definition.Workflow{}, err
	}

	return parse(data)
}

func Get(flowId string) (definition.Workflow, bool) {
	var (
		w     definition.Workflow
		found bool
	)

	if w, found = workflows[flowId]; !found {
		return definition.Workflow{}, false
	}
	return w, true
}

func LoadFromString(xmlData string) (definition.Workflow, error) {
	flow, err := parse([]byte(xmlData))
	if err != nil {
		return definition.Workflow{}, err
	}

	return flow, nil
}

func All() map[string]definition.Workflow {
	return workflows
}

func parse(data []byte) (definition.Workflow, error) {
	var w definition.Workflow
	err := xml.Unmarshal(data, &w)
	if err != nil {
		return definition.Workflow{}, err
	}

	return w, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
