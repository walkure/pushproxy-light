package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type payloadJsonItem struct {
	Type    string `json:"type"`
	Help    string `json:"help"`
	Metrics []struct {
		Value     string            `json:"value"`
		Timestamp string            `json:"timestamp"`
		Labels    map[string]string `json:"labels"`
	} `json:"metrics"`
}

type PayloadJson map[string]payloadJsonItem

var quotedEscaper = strings.NewReplacer("\\", `\\`, "\n", `\n`, "\"", `\"`)

func ConvertMetrics(w io.Writer, body string) error {
	obj, err := decodeBody(body)
	if err != nil {
		return fmt.Errorf("Cannot Decode JSON: %w", err)
	}

	transformMetrics(w, obj)

	return nil
}

func decodeBody(body string) (*PayloadJson, error) {
	var decoded PayloadJson
	if err := json.Unmarshal([]byte(body), &decoded); err != nil {
		if err, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("%+v at %v", err, err.Offset)
		} else {
			return nil, err
		}
	}
	return &decoded, nil
}

func transformMetrics(w io.Writer, payload *PayloadJson) {

	nameSlice := make([]string, len(*payload))
	index := 0
	for key := range *payload {
		nameSlice[index] = key
		index++
	}
	sort.Strings(nameSlice)

	for _, name := range nameSlice {
		metrics := (*payload)[name]
		switch metrics.Type {
		case "counter", "gauge", "untyped":
			fmt.Fprintf(w, "# HELP %s %s\n", name, metrics.Help)
			fmt.Fprintf(w, "# TYPE %s %s\n", name, metrics.Type)
			for _, metric := range metrics.Metrics {
				fmt.Fprintf(w, "%s{", name)
				fmt.Fprint(w, convertLabels(metric.Labels))
				fmt.Fprintf(w, "} %s", metric.Value)
				if metric.Timestamp != "" {
					fmt.Fprintf(w, " %s", metric.Timestamp)
				}
				fmt.Fprint(w, "\n")
			}
		default:
			log.Printf("Unsupported metrics type:%s Skip!\n", metrics.Type)
		}
	}
}

func convertLabels(labels map[string]string) string {
	nameSlice := make([]string, len(labels))
	index := 0
	for key := range labels {
		nameSlice[index] = key
		index++
	}
	sort.Strings(nameSlice)

	index = 0
	elems := make([]string, len(labels))
	for _, label := range nameSlice {
		elems[index] = fmt.Sprintf("%s=\"%s\"", label, quotedEscaper.Replace(labels[label]))
		index++
	}

	return strings.Join(elems, ",")
}
