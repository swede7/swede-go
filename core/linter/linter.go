package linter

import (
	"strings"
	"sync"

	"me.weldnor/swede/core/common"
	"me.weldnor/swede/core/parser"
)

type Linter struct {
	rootNode *parser.Node
}

type LinterError struct {
	StartPosition common.Position
	EndPosition   common.Position
	Message       string
	Severety      LinterErrorSeverety
}

type LinterErrorSeverety string

const (
	INFO  LinterErrorSeverety = "INFO"
	WARN  LinterErrorSeverety = "WARN"
	ERROR LinterErrorSeverety = "ERROR"
)

func NewLinter(rootNode *parser.Node) *Linter {
	return &Linter{
		rootNode: rootNode,
	}
}

type LinterRule func(*Linter) []LinterError

var linterRules []LinterRule = []LinterRule{
	emptyFeatureTextRule,
	emptyScenarioTextRule,
}

func (l *Linter) Lint() []LinterError {
	errors := make([]LinterError, 0)

	wg := sync.WaitGroup{}

	for _, rule := range linterRules {
		wg.Add(1)

		go func(r LinterRule) {
			defer wg.Done()
			errors = append(errors, r(l)...)
		}(rule)
	}

	wg.Wait()

	return errors
}

func emptyFeatureTextRule(l *Linter) []LinterError {
	foundedErrors := make([]LinterError, 0)

	parser.VisitNode(l.rootNode, func(n *parser.Node) {

		if n.Type != parser.FEATURE {
			return
		}

		if strings.TrimSpace(n.Value) == "" {
			e := LinterError{StartPosition: n.StartPosition,
				EndPosition: n.EndPosition,
				Message:     "Feature name is empty",
				Severety:    WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}

	})

	return foundedErrors
}

func emptyScenarioTextRule(l *Linter) []LinterError {
	foundedErrors := make([]LinterError, 0)

	parser.VisitNode(l.rootNode, func(n *parser.Node) {

		if n.Type != parser.SCENARIO {
			return
		}

		if strings.TrimSpace(n.Value) == "" {
			e := LinterError{StartPosition: n.StartPosition,
				EndPosition: n.EndPosition,
				Message:     "Scenario name is empty",
				Severety:    WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}

	})

	return foundedErrors
}
