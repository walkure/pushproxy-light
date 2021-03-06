package converter

import (
	"strings"
	"testing"
)

const testJson = `{"temperature":{"type":"gauge","help":"message",
"metrics":[{"value":"11.1","labels":{"place":"inside","location":"saitama"}},
{"value":"810.1919","labels":{"place":"outside"}}]},
"humidity":{"type":"gauge","help":"message","metrics":[{"value":"33.4"},
{"value":"21.4","labels":{"place":"outside"}}]},"pressure":{"type":"hoge",
"help":"maaa","metrics":[{"value":"1031.2","labels":{"location":"saitama"}},
{"value":"931.4","labels":{"place":"outside"}}]}}`

const promTest = `# HELP humidity message
# TYPE humidity gauge
humidity 33.4
humidity{place="outside"} 21.4
# HELP temperature message
# TYPE temperature gauge
temperature{location="saitama",place="inside"} 11.1
temperature{place="outside"} 810.1919
`

func TestTransformMetrics(t *testing.T) {

	decoded, err := decodeBody(testJson)
	if err != nil {
		t.Fatalf("Test Broken!! :%v", err)
	}
	var result strings.Builder
	transformMetrics(&result, decoded)
	if result.String() != promTest {
		t.Errorf("Result[%s] does not match [%s]\n", result.String(), promTest)
	}
}

func TestConvertLabels(t *testing.T) {
	testValue := map[string]string{
		"place":    "inside",
		"location": "kyoto",
		"org":      "hoge",
	}

	testResult := "location=\"kyoto\",org=\"hoge\",place=\"inside\""

	result := convertLabels(testValue)
	if result != testResult {
		t.Errorf("Result[%s] does not match [%s]\n", result, testResult)
	}
}
func TestConvertLabelsQuoted(t *testing.T) {
	testValue := map[string]string{
		"location": "kyoto",
		"org":      "h\"o\"ge",
	}

	testResult := "location=\"kyoto\",org=\"h\\\"o\\\"ge\""

	result := convertLabels(testValue)
	if result != testResult {
		t.Errorf("Result[%s] does not match [%s]\n", result, testResult)
	}

}
