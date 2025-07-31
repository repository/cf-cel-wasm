//go:build js && wasm

package main

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

func createCELEnvironment() (*cel.Env, error) {
	env, err := cel.NewEnv(
		cel.HomogeneousAggregateLiterals(),
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		cel.CrossTypeNumericComparisons(true),
		cel.OptionalTypes(),
		ext.Strings(ext.StringsVersion(2)),
		ext.Math(),
		cel.Function("greet",
			cel.MemberOverload("string_greet_string", []*cel.Type{cel.StringType, cel.StringType}, cel.StringType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					return types.String(fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.", rhs.Value(), lhs.Value()))
				}),
			),
		),
	)
	return env, err
}
