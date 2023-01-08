package field

import (
	"fmt"
	"regexp"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/data/schema"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

func ToColumn(f Field, targetSchema schema.Schema) (string, error) {
	snake := matchFirstCap.ReplaceAllString(f.Name(), "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)

	entity, err := targetSchema.ResolveAssociation(f.Entity())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`"%s"."%s"`, entity, snake), nil
}
