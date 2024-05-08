package linter

import (
	"strings"
	"sync"

	"github.com/swede7/swede-go/core/lang/common"
	"github.com/swede7/swede-go/core/lang/swede/parser"
)

type Linter struct {
	rootNode *common.Node
}

type LinterError struct {
	StartPosition common.Position
	EndPosition   common.Position
	Message       string
	Severity      LinterErrorSeverity
}

type LinterErrorSeverity string

const (
	INFO  LinterErrorSeverity = "INFO"
	WARN  LinterErrorSeverity = "WARN"
	ERROR LinterErrorSeverity = "ERROR"
)

func NewLinter(rootNode *common.Node) *Linter {
	return &Linter{
		rootNode: rootNode,
	}
}

type LinterRule func(*Linter) []LinterError

var linterRules []LinterRule = []LinterRule{
	emptyFeatureTextRule,
	emptyScenarioTextRule,
	emptyStepTextRule,
	featureNodeInAnotherPosition,
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

	common.VisitNode(l.rootNode, func(n *common.Node) {
		if n.Type != parser.FEATURE {
			return
		}

		if strings.TrimSpace(n.Value) == "" {
			e := LinterError{
				StartPosition: n.StartPosition,
				EndPosition:   n.EndPosition,
				Message:       "Feature name is empty",
				Severity:      WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}
	})

	return foundedErrors
}

func emptyScenarioTextRule(l *Linter) []LinterError {
	foundedErrors := make([]LinterError, 0)

	common.VisitNode(l.rootNode, func(n *common.Node) {
		if n.Type != parser.SCENARIO {
			return
		}

		if strings.TrimSpace(n.Value) == "" {
			e := LinterError{
				StartPosition: n.StartPosition,
				EndPosition:   n.EndPosition,
				Message:       "Scenario name is empty",
				Severity:      WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}
	})

	return foundedErrors
}

func emptyStepTextRule(l *Linter) []LinterError {
	foundedErrors := make([]LinterError, 0)

	common.VisitNode(l.rootNode, func(n *common.Node) {
		if n.Type != parser.STEP {
			return
		}

		if strings.TrimSpace(n.Value) == "" {
			e := LinterError{
				StartPosition: n.StartPosition,
				EndPosition:   n.EndPosition,
				Message:       "Step text is empty",
				Severity:      WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}
	})

	return foundedErrors
}

func featureNodeInAnotherPosition(l *Linter) []LinterError {
	foundedErrors := make([]LinterError, 0)

	for _, node := range l.rootNode.Children {
		if node.Type == parser.UNEXPECTED {
		}

		if node.Type == parser.FEATURE {
			break // ok
		} else {
			e := LinterError{
				StartPosition: node.StartPosition,
				EndPosition:   node.EndPosition,
				Message:       "Scenario should be declared on first position",
				Severity:      WARN,
			}

			foundedErrors = append(foundedErrors, e)
		}
	}
	return foundedErrors
}
