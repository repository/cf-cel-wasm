//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/google/cel-go/cel"
)

func analyzeType(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return map[string]any{
			"error": "Expected 2 arguments: expression and variable types",
		}
	}

	expression := args[0].String()
	variableTypesJS := args[1]

	variableTypes := convertJSObjectToMap(variableTypesJS)

	env, err := createCELEnvironment()
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Failed to create CEL env: %v", err),
		}
	}

	var envOptions []cel.EnvOption
	for k, v := range variableTypes {
		celType := inferCELType(v)
		fmt.Printf("Variable %s has type %s (Go type: %T)\n", k, celType.String(), v)
		envOptions = append(envOptions, cel.Variable(k, celType))
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

	resultType := ast.OutputType()
	
	resultTypeStr := resultType.String()

	return map[string]any{
		"resultType": resultTypeStr,
		"isValid":    true,
	}
}

func inferCELType(value any) *cel.Type {
	switch value.(type) {
	case string:
		return cel.StringType
	case int, int32, int64:
		return cel.IntType
	case float32, float64:
		return cel.DoubleType
	case bool:
		return cel.BoolType
	case map[string]any:
		return cel.DynType
	default:
		return cel.DynType
	}
}

func analyzeTypeUnknown(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return map[string]any{
			"error": "Expected 2 arguments: expression and variable names",
		}
	}

	expression := args[0].String()
	variableNamesJS := args[1]

	var variableNames []string
	if variableNamesJS.Type() == js.TypeObject && variableNamesJS.Get("length").Type() != js.TypeUndefined {
		length := variableNamesJS.Get("length").Int()
		for i := range length {
			variableNames = append(variableNames, variableNamesJS.Index(i).String())
		}
	} else {
		keys := js.Global().Get("Object").Call("keys", variableNamesJS)
		keysLength := keys.Get("length").Int()
		for i := range keysLength {
			variableNames = append(variableNames, keys.Index(i).String())
		}
	}

	env, err := createCELEnvironment()
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Failed to create CEL env: %v", err),
		}
	}

	var envOptions []cel.EnvOption
	for _, varName := range variableNames {
		envOptions = append(envOptions, cel.Variable(varName, cel.DynType))
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
			"error":   fmt.Sprintf("Compilation error: %v", issues.Err()),
			"isValid": false,
		}
	}

	resultType := ast.OutputType()
	
	resultTypeStr := resultType.String()

	return map[string]any{
		"resultType": resultTypeStr,
		"isValid":    true,
	}
}

