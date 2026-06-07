package data

import (
	"fmt"
	"reflect"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github.com/danielgtaylor/huma/v2"
)

// PaginationInput is a shared huma input struct embedded by paginated endpoints.
// Replaces the old $page/$size/$sort/$filter Gin query params.
type PaginationInput struct {
	Page   int    `query:"page"   default:"1"  minimum:"1" doc:"Page number (1-based)"`
	Size   int    `query:"size"   default:"40" minimum:"1" doc:"Page size"`
	Sort   string `query:"sort"   doc:"Sort expression, e.g. 'updatedAt desc,createdAt'"`
	Filter string `query:"filter" doc:"Filter expression"`
}

// OptionalParam wraps optional Huma params without using pointer fields, which
// Huma rejects for query/path/header/form parameters.
type OptionalParam[T any] struct {
	Value T
	IsSet bool
}

func (o *OptionalParam[T]) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeOf(o.Value))
}

func (o *OptionalParam[T]) Receiver() reflect.Value {
	return reflect.ValueOf(o).Elem().Field(0)
}

func (o *OptionalParam[T]) OnParamSet(isSet bool, _ any) {
	o.IsSet = isSet
}

func (o *OptionalParam[T]) Ptr() *T {
	if !o.IsSet {
		return nil
	}
	return &o.Value
}

// MapPageResponse converts a Page[DTO] to a PageResponse[VO] using a converter function.
func MapPageResponse[DTO any, VO any](page Page[DTO], converter func(DTO) VO) PageResponse[VO] {
	vos := make([]VO, len(page.Content()))
	for i, dto := range page.Content() {
		vos[i] = converter(dto)
	}
	return PageResponse[VO]{
		Content:          vos,
		PageNumber:       page.PageNumber(),
		PageSize:         page.PageSize(),
		TotalPage:        page.TotalPage(),
		NumberOfElements: page.NumberOfElements(),
	}
}

// PageRequestFromInput builds a PageRequest from a PaginationInput.
func PageRequestFromInput(p PaginationInput) (PageRequest, error) {
	pg, err := pageable.NewPageableFromQuery(p.Page, p.Size)
	if err != nil {
		return nil, fmt.Errorf("invalid pagination: %w", err)
	}

	sorts, err := sortable.NewSortableFromQuery(p.Sort)
	if err != nil {
		return nil, fmt.Errorf("invalid sort: %w", err)
	}

	filters, err := filterable.NewFilterableFromQuery(p.Filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	return NewPageRequest(pg, sorts, filters), nil
}

// WithTag returns a Huma API group that applies tag to operations without an explicit tag.
func WithTag(api huma.API, tag string) huma.API {
	group := huma.NewGroup(api)
	group.UseSimpleModifier(func(op *huma.Operation) {
		if len(op.Tags) == 0 {
			op.Tags = []string{tag}
		}
	})
	return group
}

// SchemaNamer prefixes project schemas with the internal module name so
// public DTOs like report.PeriodResponse and general_ledger.PeriodResponse do
// not collide in Huma's global OpenAPI component registry.
func SchemaNamer(t reflect.Type, hint string) string {
	name := huma.DefaultSchemaNamer(t, hint)
	module := internalModuleName(t)
	if module == "" {
		return name
	}
	return module + name
}

func internalModuleName(t reflect.Type) string {
	for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}

	pkgPath := t.PkgPath()
	const marker = "/internal/"
	idx := strings.Index(pkgPath, marker)
	if idx == -1 {
		return ""
	}

	rest := pkgPath[idx+len(marker):]
	module, _, ok := strings.Cut(rest, "/")
	if !ok || module == "" {
		return ""
	}
	return exportedName(module)
}

func exportedName(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '_' || r == '-'
	})

	var out strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		out.WriteString(strings.ToUpper(part[:1]))
		out.WriteString(part[1:])
	}
	return out.String()
}
