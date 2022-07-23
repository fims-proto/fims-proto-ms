package query

import "context"

type ResolveIdentifierConfigurationReadModel interface {
	ResolveIdentifierConfiguration(ctx context.Context, targetBusinessObject string, objectsToMatch map[string]string) (IdentifierConfiguration, error)
}
