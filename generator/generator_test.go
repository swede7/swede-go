package generator_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/swede7/swede-go/generator"
	"os"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	testCases := getAllTestCases()

	for _, testCase := range testCases {
		t.Run("testing "+testCase.name, func(t *testing.T) {
			g := generator.NewGenerator(generator.Options{
				Source: testCase.source,
				Debug:  true,
			})

			result := g.Generate()
			fmt.Println(result)
			assert.Equal(t, testCase.expected, result)
		})
	}

}

type testCase struct {
	name     string
	source   string
	expected string
}

func getAllTestCases() []testCase {
	dir, err := os.ReadDir("./testdata/")
	if err != nil {
		return nil
	}

	testCaseNames := make(map[string]bool)

	for _, file := range dir {
		fileName := file.Name()
		testCaseName := strings.ReplaceAll(fileName, "_expected", "")
		testCaseNames[testCaseName] = true
	}

	testCases := make([]testCase, 0)

	for testCaseName := range testCaseNames {
		const prefix = "./testdata/"
		input := readFileAsString(prefix + testCaseName)
		expected := readFileAsString(prefix + testCaseName + "_expected")
		testCase := testCase{
			name:     testCaseName,
			source:   input,
			expected: expected,
		}
		testCases = append(testCases, testCase)
	}

	return testCases
}

func readFileAsString(path string) string {
	_bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(_bytes)
}
