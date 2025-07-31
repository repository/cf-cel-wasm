//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

func evalCEL(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return map[string]any{
			"error": "Expected 2 arguments: expression and input data",
		}
	}

	expression := args[0].String()
	inputJS := args[1]

	inputData := convertJSObjectToMap(inputJS)

	env, err := createCELEnvironment()
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Failed to create CEL env: %v", err),
		}
	}

	var envOptions []cel.EnvOption
	for k := range inputData {
		envOptions = append(envOptions, cel.Variable(k, cel.DynType))
	}

	extendedEnv, err := env.Extend(envOptions...)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Failed to extend env with variables: %v", err),
		}
	}

	ast, issues := extendedEnv.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return map[string]any{
			"error": fmt.Sprintf("Compilation error: %v", issues.Err()),
		}
	}

	program, err := extendedEnv.Program(ast)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Program creation error: %v", err),
		}
	}

	result, _, err := program.Eval(inputData)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Evaluation error: %v", err),
		}
	}

	return map[string]any{
		"result": convertCELResultToJS(result),
	}
}

func convertCELResultToJS(result ref.Val) any {
	return result.Value()
}
