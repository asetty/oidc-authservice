package google

import (
	"github.com/arrikto/oidc-authservice/groups"
	"github.com/arrikto/oidc-authservice/settings"
)

//TODO decide how to authenticate to get groups
// * user's token
// * google service account file similar to dex
// * ...
// * maybe should do at same time as getting other claims (GetUserInfo)
//   (refactor from GroupsMethod to ProviderFlavor, customize getting
//    groups and userid, if desired.)
type googleMethod struct{}

func NewGoogleMethod(config *settings.Config) (groups.GroupsMethod, error) {
	// _, err := admin.NewService(context.Background())
	// if err != nil {
	// 	return nil, err
	// }
	return &googleMethod{}, nil
}

func (_ *googleMethod) GetGroups(claims map[string]interface{}) ([]string, error) {
	//TODO implement
	// var userGroups []string
	// var err error
	// groupsList := &admin.Groups{}
	// for {
	// 	groupsList, err = c.adminSrv.Groups.List().
	// 		UserKey(email).PageToken(groupsList.NextPageToken).Do()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not list groups: %v", err)
	// 	}

	// 	for _, group := range groupsList.Groups {
	// 		// TODO (joelspeed): Make desired group key configurable
	// 		userGroups = append(userGroups, group.Email)
	// 	}

	// 	if groupsList.NextPageToken == "" {
	// 		break
	// 	}
	// }

	// return userGroups, nil

	// TODO
	// Get groups either using user token or api credentials?
	// how do we get token source like in authenticator_session?
	// other refactoring?
	// option.WithTokenSource()
	// client, err := admin.NewService(context.Background())
	// if err != nil {
	// 	return nil, err
	// }
	return []string{}, nil
}
