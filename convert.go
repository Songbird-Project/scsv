package scsv

import "fmt"

func AsCSV(parsedSCSV KeyValuePairs) string {
	conversion := ""

	for key, values := range parsedSCSV {
		for _, value := range values {
			csvValue := ""

			for idx, column := range value {
				csvValue += column

				if idx != len(value)-1 {
					csvValue += ","
				}
			}

			conversion += fmt.Sprintf("%s,%s\n", key, csvValue)
		}
	}

	return conversion
}

func AsJSON(parsedSCSV KeyValuePairs) string {
	return "JSON conversions are incorrect since SCSV v1.0.0, they will be fixed with v1.0.1"

	conversion := "{\n"
	keyIdx := 0

	for key, values := range parsedSCSV {
		keyIdx++

		arrayDelim := ",\n"
		if keyIdx == len(parsedSCSV) {
			arrayDelim = "\n"
		}
		conversion += fmt.Sprintf("  \"%s\": [\n", key)

		for i, value := range values {
			delim := ",\n"
			if i == len(values)-1 {
				delim = "\n"
			}
			conversion += fmt.Sprintf("    \"%s\"%s", value, delim)
		}

		conversion += fmt.Sprintf("  ]%s", arrayDelim)
	}

	conversion += "}\n"

	return conversion
}

func AsYAML(parsedSCSV KeyValuePairs) string {
	return "YAML conversions are incorrect since SCSV v1.0.0, they will be fixed with v1.0.1"

	conversion := ""

	for key, values := range parsedSCSV {
		conversion += fmt.Sprintf("%s:", key)

		if len(values) <= 1 {
			conversion += " []\n"
			continue
		} else {
			conversion += "\n"
		}

		for _, value := range values {
			conversion += fmt.Sprintf("  - %s\n", value)
		}
	}

	return conversion
}
