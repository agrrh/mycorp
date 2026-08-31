package template

// TODO: Replace with actual jinja package

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/agrrh/mycorp/internal/domain/modules"
	"github.com/agrrh/mycorp/internal/domain/scenario"
)

type TemplateEngine struct {
	tmpl *template.Template
}

type TemplateData struct {
	Steps  modules.PrevStepsResults
	Inputs scenario.SpecInputs
	Params scenario.SpecStepParams
}

func (td *TemplateData) GetStep(stepName string) map[string]any {
	if step, ok := td.Steps[stepName]; ok {
		return step
	}
	return nil
}

func (td *TemplateData) GetInput(inputName string) *scenario.SpecInputParameter {
	for _, input := range td.Inputs {
		if input.Name == inputName {
			return &input
		}
	}
	return nil
}

func NewTemplateEngine(tmplStr string) (*TemplateEngine, error) {
	converted := convertJinjaToGoTemplate(tmplStr)

	funcMap := template.FuncMap{
		"json": func(v any) string {
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("%v", v)
			}
			return string(b)
		},
		"toJSON": func(v any) string {
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("%v", v)
			}
			return string(b)
		},
		"int": func(v any) (int, error) {
			switch val := v.(type) {
			case int:
				return val, nil
			case float64:
				return int(val), nil
			case string:
				return strconv.Atoi(val)
			default:
				return 0, fmt.Errorf("cannot convert %v to int", v)
			}
		},
		"eq": func(a, b any) bool {
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		},
		"ne": func(a, b any) bool {
			return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
		},
		"gt": func(a, b any) bool {
			return fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b)
		},
		"lt": func(a, b any) bool {
			return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
		},
		"ge": func(a, b any) bool {
			return fmt.Sprintf("%v", a) >= fmt.Sprintf("%v", b)
		},
		"le": func(a, b any) bool {
			return fmt.Sprintf("%v", a) <= fmt.Sprintf("%v", b)
		},
	}

	tmpl, err := template.New("output").Funcs(funcMap).Parse(converted)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &TemplateEngine{tmpl: tmpl}, nil
}

func (te *TemplateEngine) Execute(data *TemplateData) (string, error) {
	var buf bytes.Buffer

	if err := te.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func RenderOutput(tmplStr string, stepResults modules.PrevStepsResults, inputs scenario.SpecInputs) (string, error) {
	return RenderOutputWithParams(tmplStr, stepResults, inputs, nil)
}

func RenderOutputWithParams(tmplStr string, stepResults modules.PrevStepsResults, inputs scenario.SpecInputs, params scenario.SpecStepParams) (string, error) {
	if tmplStr == "" {
		return "", nil
	}

	te, err := NewTemplateEngine(tmplStr)
	if err != nil {
		return "", err
	}

	data := &TemplateData{
		Steps:  stepResults,
		Inputs: inputs,
		Params: params,
	}

	return te.Execute(data)
}

func isStaticOutput(output scenario.SpecOutput) bool {
	s := string(output)
	return !strings.Contains(s, "{{") && !strings.Contains(s, "{%")
}

func ProcessOutput(output scenario.SpecOutput, stepResults modules.PrevStepsResults, inputs scenario.SpecInputs) (string, error) {
	return ProcessOutputWithParams(output, stepResults, inputs, nil)
}

func ProcessOutputWithParams(output scenario.SpecOutput, stepResults modules.PrevStepsResults, inputs scenario.SpecInputs, params scenario.SpecStepParams) (string, error) {
	outputStr := string(output)

	if outputStr == "" {
		return "", nil
	}

	if isStaticOutput(output) {
		return outputStr, nil
	}

	return RenderOutputWithParams(outputStr, stepResults, inputs, params)
}

func convertJinjaToGoTemplate(tmplStr string) string {
	result := tmplStr

	result = convertControlStructures(result)
	result = convertStepsAccess(result)
	result = convertInputsAccess(result)
	result = convertParamsAccess(result)
	result = convertComparisonOperators(result)

	return result
}

func convertComparisonOperators(s string) string {
	s = regexp.MustCompile(`(\([^)]+\))\s*==\s*(\d+)`).ReplaceAllString(s, "eq $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*==\s*"([^"]*)"`).ReplaceAllString(s, "eq $1 \"$2\"")
	s = regexp.MustCompile(`(\([^)]+\))\s*==\s*'([^']*)'`).ReplaceAllString(s, "eq $1 \"$2\"")
	s = regexp.MustCompile(`(\([^)]+\))\s*!=(\d+)`).ReplaceAllString(s, "ne $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*!=([^0-9])`).ReplaceAllString(s, "ne $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*>(\d+)`).ReplaceAllString(s, "gt $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*<(\d+)`).ReplaceAllString(s, "lt $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*>([^0-9])`).ReplaceAllString(s, "gt $1 $2")
	s = regexp.MustCompile(`(\([^)]+\))\s*<([^0-9])`).ReplaceAllString(s, "lt $1 $2")
	return s
}

func convertControlStructures(s string) string {
	s = regexp.MustCompile(`\{% if\s+(.+?)\s*%}`).ReplaceAllString(s, "{{ if $1 }}")
	s = regexp.MustCompile(`\{%\s*else\s*%\}`).ReplaceAllString(s, "{{ else }}")
	s = regexp.MustCompile(`\{%\s*end\s*%\}`).ReplaceAllString(s, "{{ end }}")
	return s
}

func convertStepsAccess(s string) string {
	re := regexp.MustCompile(`steps\[["']([^"']+)["']\]\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		extracted := re.FindStringSubmatch(match)
		if len(extracted) < 3 {
			return match
		}
		stepName := extracted[1]
		field := extracted[2]
		return fmt.Sprintf("(index .Steps %q %q)", stepName, field)
	})
}

func convertInputsAccess(s string) string {
	re := regexp.MustCompile(`inputs\[["']([^"']+)["']\]\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		extracted := re.FindStringSubmatch(match)
		if len(extracted) < 3 {
			return match
		}
		inputName := extracted[1]
		field := extracted[2]
		return fmt.Sprintf("(index .Inputs %q %q)", inputName, field)
	})
}

func convertParamsAccess(s string) string {
	re := regexp.MustCompile(`params\[["']([^"']+)["']\]\.([a-zA-Z_][a-zA-Z0-9_]*)`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		extracted := re.FindStringSubmatch(match)
		if len(extracted) < 3 {
			return match
		}
		paramName := extracted[1]
		field := extracted[2]
		return fmt.Sprintf("(index .Params %q %q)", paramName, field)
	})
}
