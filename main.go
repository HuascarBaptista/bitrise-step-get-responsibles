package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
	"io/ioutil"
	"os"
	"strings"
)

type Responsible struct {
	Key              string   `json:"key"`
	Modules          []string `json:"modules"`
	SlackResponsible []string `json:"slack_responsible"`
}

type Config struct {
	AllowedKeys       string `env:"jira_keys"`
	Branch            string `env:"branch"`
	PathConfiguration string `env:"path_configuration"`
}

func main() {
	/*
		var cfg Config = Config{
			Folders:           "tools\nbasket\nbase\n",
			PathConfiguration: "tools/responsible.json",
			Branch:            "fix/SHP-22/huascar",
			AllowedKeys:       "BAS|SHP|OT",
		}
	*/
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	file, _ := ioutil.ReadFile(cfg.PathConfiguration)

	var jsonDataArray []Responsible

	errorJson := json.Unmarshal(file, &jsonDataArray)

	if errorJson != nil {
		failf("Error parsing JSON\n")
	}

	branchKey := extraBranchKey(cfg.Branch, cfg.AllowedKeys)

	var branchKeyIndex = getIndexOfKeyProject(jsonDataArray, branchKey)
	responsible := jsonDataArray[branchKeyIndex].SlackResponsible

	if err := tools.ExportEnvironmentWithEnvman("RESPONSIBLES", strings.Join(responsible, ", ")); err != nil {
		failf("error exporting variable", err)
	}

	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}

func getIndexOfKeyProject(jsonDataArray []Responsible, branchKey string) int {
	if branchKey != "" {
		for i := 0; i < len(jsonDataArray); i++ {
			if jsonDataArray[i].Key == branchKey {
				return i
			}
		}
	}
	return -1
}

func extraBranchKey(branch string, allowedKeys string) string {
	dividerBySlashPath := strings.Split(branch, "/")
	allowedKeysSeparated := strings.Split(allowedKeys, "|")
	if len(dividerBySlashPath) > 1 {
		var key = stringContainsInArray(dividerBySlashPath[1], allowedKeysSeparated)
		if key != "" {
			return key
		} else {
			fmt.Printf("Key %s in branch %s don't founded in allowed keys: %s\n", dividerBySlashPath[1], branch, allowedKeys)
		}
	}
	return ""
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(0)
}

func stringContainsInArray(branchPart string, allowedKeysSeparated []string) string {
	for _, allowedKey := range allowedKeysSeparated {
		fmt.Println(normalize(branchPart))
		fmt.Println(normalize(allowedKey))
		fmt.Println(strings.Contains(normalize(branchPart), normalize(allowedKey)))
		fmt.Println(strings.Contains(normalize(allowedKey), normalize(branchPart)))
		if strings.Contains(normalize(branchPart), normalize(allowedKey)) {
			return allowedKey
		}
	}
	return ""
}

func normalize(stringToCheck string) string {
	return strings.TrimSpace(strings.ToLower(stringToCheck))
}
