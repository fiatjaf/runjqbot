package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/itchyny/gojq"
)

func runjq(actualCode string) (res string, err error) {
	code := strings.TrimSpace(actualCode)
	if code == "" {
		return "", nil
	}

	query, err := gojq.Parse(code)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	iter := query.RunWithContext(ctx, map[string]interface{}{})
	v, ok := iter.Next()
	if !ok {
		return "", nil
	}
	if err, ok := v.(error); ok {
		return "", err
	}

	bres, _ := json.MarshalIndent(v, "", "  ")
	result := string(bres)
	log.Debug().Str("code", actualCode).Str("ret", result).Msg("ran")

	return result, nil
}
