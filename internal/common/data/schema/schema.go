package schema

type Schema interface {
	ResolveAssociation(entity string) (string, error)
}
