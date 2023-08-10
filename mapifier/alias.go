package mapifier

import "github.com/ivancorrales/knoa/mapifier/internal"

const (
	YAML = internal.YAML
	JSON = internal.JSON
)

var (
	AsJSON = internal.WithOutputFormat(JSON)
	AsYAML = internal.WithOutputFormat(YAML)
)

type SetterOpt = internal.SetterOpt

var (
	WithFuncPrefix   = internal.WithFuncPrefix
	WithStringPrefix = internal.WithStringPrefix
)

var (
	SetOp   = internal.SetOp
	UnsetOp = internal.UnsetOp
)
