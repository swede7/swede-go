package main_test

import (
	"fmt"
	"testing"
	"time"

	"me.weldnor/swede/core/formatter"
	"me.weldnor/swede/core/parser"
)

func TestBenchmark(t *testing.T) {
	countOfTests := 2000

	fmt.Println(getExecutionTime(formatFile, countOfTests))
	fmt.Println(getExecutionTime(formatFileParallel, countOfTests))
}

func formatFileParallel() {
	parserResult := parser.ParseCode(CODE)
	rootNode := parserResult.RootNode
	formatter.NewFormatter(&rootNode).FormatParallel()
}

func formatFile() {
	parserResult := parser.ParseCode(CODE)
	rootNode := parserResult.RootNode
	formatter.NewFormatter(&rootNode).Format()
}

func getExecutionTime(f func(), countOfTests int) int64 {
	startTime := time.Now().UnixMilli()

	for range countOfTests {
		f()
	}

	endTime := time.Now().UnixMilli()
	return endTime - startTime
}

var CODE string = `
@example 
Feature: empty

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

@example 
Feature: empty

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

@example 
Feature: empty

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"

# comment

@positive 
Scenario: Test addition
- Add "2" and "2"
- Check that result is "4"

@negative @test 
Scenario: Test addition, but result is not correct
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
- Add "2" and "2"
- Check that result is "5"
`
