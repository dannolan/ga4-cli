package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func decodeJSONBody(path string, dest any) error {
	var reader io.Reader
	if path == "" || path == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	}
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("decode JSON body: %w", err)
	}
	return nil
}

func propertyResource(propertyID string) string {
	return "properties/" + strings.TrimPrefix(propertyID, "properties/")
}
