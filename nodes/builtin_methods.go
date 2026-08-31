package nodes

import (
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func vals2anys[T any](vals []T) []any {
	res := make([]any, len(vals))
	for i, v := range vals {
		res[i] = v
	}
	return res
}

func blStringSplit(cx base.Context, inst any, args []any) (any, error) {
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.split: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("No args of in string.split")
	}
	var sep string
	switch sepv := args[0].(type) {
	case string:
		sep = sepv
	case objects.Glif:
		sep = string([]rune{sepv})
	default:
		return nil, fmt.Errorf("Bad separator in string.split: %T", args[0])
	}
	ee := strings.Split(s, sep)
	res := objects.NewListVal(vals2anys(ee))
	return res, nil
}
