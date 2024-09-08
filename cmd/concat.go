package main

import (
	"fmt"
	"strings"
)

func concat(v ...any) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(fmt.Sprintf("%+v", value))
	}

	return builder.String()
}
