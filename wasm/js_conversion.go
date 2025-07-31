//go:build js && wasm

package main

import (
	"syscall/js"
)

func convertJSObjectToMap(inputJS js.Value) map[string]any {
	inputData := make(map[string]any)

	if inputJS.Type() != js.TypeObject {
		return inputData
	}

	keys := js.Global().Get("Object").Call("keys", inputJS)
	keysLength := keys.Get("length").Int()

	for i := range keysLength {
		key := keys.Index(i).String()
		value := inputJS.Get(key)

		inputData[key] = convertJSValueToGoValue(value)
	}

	return inputData
}

func convertJSValueToGoValue(value js.Value) any {
	switch value.Type() {
	case js.TypeString:
		return value.String()
	case js.TypeNumber:
		return value.Float()
	case js.TypeBoolean:
		return value.Bool()
	case js.TypeObject:
		result := make(map[string]any)
		keys := js.Global().Get("Object").Call("keys", value)
		keysLength := keys.Get("length").Int()

		for i := range keysLength {
			key := keys.Index(i).String()
			nestedValue := value.Get(key)
			result[key] = convertJSValueToGoValue(nestedValue)
		}
		return result
	default:
		return value.String()
	}
}
