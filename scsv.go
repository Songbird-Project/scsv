package scsv

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CurrentFile struct {
	ValuePrecedence bool
	StrictMode      bool

	LineNumber     int
	FlowKey        string
	FlowKeyLines   int
	FlowValue      string
	FlowValueLines int
}

// eg. map[extra:["bat", "eza"] aur:["zen-browser-bin"]]
type KeyValuePairs map[string][]string

func ParseFile(path string) (KeyValuePairs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileOptions := &CurrentFile{
		ValuePrecedence: false,
		StrictMode:      false,

		LineNumber: 1,
	}
	keySets := make(KeyValuePairs)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		newFlowKey := false
		newFlowValue := false

		if fileOptions.FlowKeyLines > 0 {
			fileOptions.FlowKeyLines--
		} else if fileOptions.FlowKeyLines == 0 {
			fileOptions.FlowKey = ""
		}

		if fileOptions.FlowValueLines > 0 {
			fileOptions.FlowValueLines--
		} else if fileOptions.FlowValueLines == 0 {
			fileOptions.FlowValue = ""
		}

		if len(line) == 0 {
			fileOptions.LineNumber++
			continue

			// Check if line is a parsing option
		} else if fmt.Sprintf("%.*s", 2, line) == "#@" {
			keyValuePair := strings.Split(strings.TrimPrefix(line, "#@"), ",")
			if len(keyValuePair) > 2 || len(keyValuePair) < 2 {
				return nil, fmt.Errorf(
					"Parsing options should be in the form `key,value`: line %d",
					fileOptions.LineNumber,
				)
			}

			key := strings.ToLower(keyValuePair[0])
			value := strings.ToLower(keyValuePair[1])

			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}

			switch key {
			case "valueprecedence":
				fileOptions.ValuePrecedence = boolValue
			case "strictmode":
				fileOptions.StrictMode = boolValue
			}
		} else if line == "--" {
			fileOptions.FlowKeyLines = 0
			fileOptions.FlowValueLines = 0

			// Ignore comment lines
		} else if line[0] == '#' {
			fileOptions.LineNumber++
			continue

			// Handle other lines as key-value pairs
		} else {
			keyValuePair := strings.Split(line, ",")
			key := keyValuePair[0]
			value := keyValuePair[1]

			if strings.Contains(key, "|") {
				if fileOptions.FlowKey != "" && fileOptions.StrictMode {
					return nil, fmt.Errorf("Flow keys must be cleared before new ones are defined in strict mode: line %d", fileOptions.LineNumber)
				}

				flowKey := strings.Split(key, "|")
				flow := flowKey[0]
				key = flowKey[1]

				if flow != "" {
					flowLines, err := strconv.Atoi(flow)
					if err != nil {
						return nil, fmt.Errorf("Flow key must use a valid line controller if one is present: line %d", fileOptions.LineNumber)
					}

					fileOptions.FlowKeyLines = flowLines
				} else {
					fileOptions.FlowKeyLines = -1
				}

				fileOptions.FlowKey = key
			} else if key == "" {
				if fileOptions.FlowKey != "" {
					key = fileOptions.FlowKey
				} else {
					return nil, fmt.Errorf("No key found: line %d", fileOptions.LineNumber)
				}
			} else {
				fileOptions.FlowKey = ""
				fileOptions.FlowKeyLines = 0
			}

			if strings.Contains(value, "|") {
				if fileOptions.FlowValue != "" && fileOptions.StrictMode {
					return nil, fmt.Errorf("Flow values must be cleared before new ones are defined in strict mode: line %d", fileOptions.LineNumber)
				}

				flowValue := strings.Split(value, "|")
				flow := flowValue[0]
				value = flowValue[1]

				if flow != "" {
					flowLines, err := strconv.Atoi(flow)
					if err != nil {
						return nil, fmt.Errorf("Flow value must use a valid line controller if one is present: line %d", fileOptions.LineNumber)
					}

					fileOptions.FlowValueLines = flowLines
				} else {
					fileOptions.FlowValueLines = -1
				}

				fileOptions.FlowValue = value
			} else if value == "" {
				if fileOptions.FlowValue != "" {
					value = fileOptions.FlowValue
				} else {
					if fileOptions.StrictMode {
						return nil, fmt.Errorf("No value found: line %d", fileOptions.LineNumber)
					}
				}
			} else {
				fileOptions.FlowValue = ""
				fileOptions.FlowValueLines = 0
			}

			if newFlowKey && newFlowValue {
				if fileOptions.ValuePrecedence {
					fileOptions.FlowKeyLines = 0
				} else {
					fileOptions.FlowValueLines = 0
				}
			}

			keySets[key] = append(keySets[key], value)
		}

		fileOptions.LineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keySets, nil
}
