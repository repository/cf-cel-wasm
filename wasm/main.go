//go:build js && wasm

package main

import (
	"syscall/js"
)

func main() {
	celAPI := js.Global().Get("Object").New()
	celAPI.Set("eval", js.FuncOf(evalCEL))
	celAPI.Set("analyzeType", js.FuncOf(analyzeType))
	celAPI.Set("analyzeTypeUnknown", js.FuncOf(analyzeTypeUnknown))
	js.Global().Set("$__celery", celAPI)
	select {}
}
