package scsv

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FileInfo struct {
	ValuePrecedence bool
	StrictMode      bool

	LineNumber   int
	FlowKey      string
	FlowKeyLines int
	FlowValues   map[int][]any // map[COLUMN:[FLOW_VALUE_1,FLOW_LINES_1] COLUMN:[FLOW_VALUE_2,FLOW_LINES_2]]
}

// eg. map[extra:[["bat", "1.2"], ["eza", "3.4"]] aur:[["zen-browser-bin", "4.3"]]]
type KeyValuePairs map[string][][]string

func ParseFile(path string) (KeyValuePairs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return parse(func() (string, bool) {
		if scanner.Scan() {
			return scanner.Text(), true
		}
		return "", false
	}, scanner.Err)
}

func ParseString(scsv string) (KeyValuePairs, error) {
	lines := strings.Split(scsv, "\n")
	idx := 0
	return parse(func() (string, bool) {
		if idx >= len(lines) {
			return "", false
		}
		line := lines[idx]
		idx++
		return line, true
	}, func() error { return nil })
}

func parse(nextLine func() (string, bool), scanErr func() error) (KeyValuePairs, error) {
	fileOptions := &FileInfo{
		ValuePrecedence: false,
		StrictMode:      false,
		LineNumber:      1,
		FlowValues:      make(map[int][]any),
	}

	keySets := make(KeyValuePairs)

	for {
		line, ok := nextLine()
		if !ok {
			break
		}

		newFlowKey := false
		newFlowValue := false

		if fileOptions.FlowKeyLines != -1 {
			fileOptions.FlowKeyLines -= 1
		}

		if fileOptions.FlowKeyLines == 0 {
			fileOptions.FlowKey = ""
		}

		for key := range fileOptions.FlowValues {
			if fileOptions.FlowValues[key][1] != -1 {
				fileOptions.FlowValues[key][1] = fileOptions.FlowValues[key][1].(int) - 1
			}

			if fileOptions.FlowValues[key][1] == 0 {
				delete(fileOptions.FlowValues, key)
			}
		}

		if len(line) == 0 {
			fileOptions.LineNumber++
			continue
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
			for key := range fileOptions.FlowValues {
				delete(fileOptions.FlowValues, key)
			}
		} else if line[0] == '#' {
			fileOptions.LineNumber++
			continue
		} else {
			keyValuePair := strings.Split(line, ",")
			key := keyValuePair[0]
			values := keyValuePair[1:]

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
				newFlowKey = true
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

			parsedValues := []string{}

			for idx, value := range values {
				if strings.Contains(value, "|") {
					_, valueExists := fileOptions.FlowValues[idx]
					if valueExists && fileOptions.StrictMode {
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

						fileOptions.FlowValues[idx] = []any{value, flowLines}
					} else {
						fileOptions.FlowValues[idx] = []any{value, -1}
					}

					newFlowValue = true
				} else if value == "" {
					if len(fileOptions.FlowValues[idx]) <= 0 {
						if fileOptions.StrictMode {
							return nil, fmt.Errorf("No value found: line %d", fileOptions.LineNumber)
						}
					} else {
						value = fileOptions.FlowValues[idx][0].(string)
					}
				} else {
					if len(fileOptions.FlowValues[idx]) > 0 {
						fileOptions.FlowValues[idx][0] = ""
						fileOptions.FlowValues[idx][1] = 0
					}
				}

				parsedValues = append(parsedValues, value)
			}

			if newFlowKey && newFlowValue && len(values) > 1 {
				if fileOptions.ValuePrecedence {
					fileOptions.FlowKeyLines = 0
				} else {
					for key := range fileOptions.FlowValues {
						delete(fileOptions.FlowValues, key)
					}
				}
			}

			keySets[key] = append(keySets[key], parsedValues)
		}

		fileOptions.LineNumber++
	}

	if err := scanErr(); err != nil {
		return nil, err
	}

	return keySets, nil
}
