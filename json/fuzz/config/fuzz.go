package fuzzconfig

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/json"
)

func Fuzz(data []byte) int {
	_, diags := json.Parse(data, "<fuzz-conf>", hcl.Pos{Byte: 0, Line: 1, Column: 1})

	if diags.HasErrors() {
		return 0
	}

	return 1
}
