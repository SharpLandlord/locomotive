package railway

import (
	"context"

	"github.com/brody192/locomotive/internal/railway/gql/queries"
	"github.com/flexstack/uuid"
)

// VerifyAllServicesExistWithinEnvironment reports whether every wanted service is attached
// to the environment, matching against service instances rather than deployments so the
// check is independent of deployment state — a service that has never been deployed, or
// whose deployments are crashed, removed, or otherwise not running, still counts as
// existing.
func VerifyAllServicesExistWithinEnvironment(g *GraphQLClient, services []uuid.UUID, environmentID uuid.UUID) (bool, []uuid.UUID, []uuid.UUID, error) {
	environment := &queries.EnvironmentData{}

	variables := map[string]any{
		"id": environmentID,
	}

	if err := g.Client.Exec(context.Background(), queries.EnvironmentQuery, &environment, variables); err != nil {
		return false, nil, nil, err
	}

	foundServices := []uuid.UUID{}
	missingServices := []uuid.UUID{}

	for _, service := range services {
		found := false

		for _, edge := range environment.Environment.ServiceInstances.Edges {
			if edge.Node.ServiceID == service {
				found = true
				break
			}
		}

		if found {
			foundServices = append(foundServices, service)
		} else {
			missingServices = append(missingServices, service)
		}
	}

	return (len(missingServices) == 0), foundServices, missingServices, nil
}
