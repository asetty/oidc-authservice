package claim

import (
	"github.com/arrikto/oidc-authservice/groups"
	"github.com/arrikto/oidc-authservice/settings"
)

type claimMethod struct {
	key string
}

func NewClaimMethod(config *settings.Config) (groups.GroupsMethod, error) {
	return &claimMethod{key: config.GroupsClaim}, nil
}

func (c *claimMethod) GetGroups(claims map[string]interface{}) ([]string, error) {
	groups := []string{}
	groupsClaim := claims[c.key]
	if groupsClaim != nil {
		groups = interfaceSliceToStringSlice(groupsClaim.([]interface{}))
	}
	return groups, nil
}

func interfaceSliceToStringSlice(in []interface{}) []string {
	if in == nil {
		return nil
	}

	res := []string{}
	for _, elem := range in {
		res = append(res, elem.(string))
	}
	return res
}
