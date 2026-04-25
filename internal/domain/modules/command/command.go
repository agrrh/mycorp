package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/agrrh/mycorp/internal/domain/modules"
	"github.com/agrrh/mycorp/internal/domain/scenario"
)

// CommandModule handles direct command execution.
type CommandModule struct {
	Name string
}

// New creates a new CommandModule instance.
func New() *CommandModule {
	return &CommandModule{
		Name: "command",
	}
}

// Run executes the specified command with the given parameters.
//
// Expected stepParams:
//   - cmd (required): string - the command to execute
//   - args (optional): []string - command arguments
//   - env (optional): []string - environment variables in KEY=VALUE format
//   - pwd (optional): string - working directory for command execution
//   - timeout (optional): string - command timeout (e.g., "5s", "1m")
func (m *CommandModule) Run(ctx context.Context, stepName string, scInputs scenario.SpecInputs, stepParams scenario.SpecStepParams, prevStepResults modules.PrevStepsResults) (modules.StepResults, error) {
	results := make(modules.StepResults)

	// Parse required command
	cmdStr, ok := stepParams["cmd"].(string)
	if !ok || cmdStr == "" {
		return results, fmt.Errorf("missing or invalid 'cmd' parameter in step params")
	}

	// Parse optional arguments
	var args []string
	if argsRaw, exists := stepParams["args"]; exists {
		args = toStringSlice(argsRaw)
	}

	// Parse optional environment variables
	var env []string
	if envRaw, exists := stepParams["env"]; exists {
		env = toStringSlice(envRaw)
	}

	// Parse optional working directory
	var dir string
	if dirRaw, exists := stepParams["pwd"]; exists {
		dir = dirRaw.(string)
	}

	// TODO: Use relatively to some Scenario dir

	// Create the command
	cmd := exec.CommandContext(ctx, cmdStr, args...)

	// Set working directory if provided
	if dir != "" {
		cmd.Dir = dir
	}

	// Set environment variables if provided (merge with current env)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	} else {
		cmd.Env = os.Environ()
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// TODO: Run in a thread
	// TODO: Respect timeout param

	// Execute the command
	err := cmd.Run()

	// Record results
	results["stdout"] = strings.TrimSuffix(stdout.String(), "\n")
	results["stderr"] = strings.TrimSuffix(stderr.String(), "\n")
	results["exit_code"] = cmd.ProcessState.ExitCode()

	if err != nil {
		// Return results along with error for debugging
		results["error"] = err.Error()
		return results, fmt.Errorf("command execution failed: %w", err)
	}

	return results, nil
}

// toStringSlice converts various input types to a string slice.
func toStringSlice(v any) []string {
	switch val := v.(type) {
	case []string:
		return val
	case []any:
		result := make([]string, len(val))
		for i, item := range val {
			if str, ok := item.(string); ok {
				result[i] = str
			}
		}
		return result
	case string:
		return strings.Fields(val)
	default:
		return nil
	}
}
