package google

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/arrikto/oidc-authservice/groups"
	"github.com/arrikto/oidc-authservice/settings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

var emailClaim = "email"

type googleMethod struct {
	adminSvc *admin.Service
}

func createDirectoryService(serviceAccountFilePath string, email string) (*admin.Service, error) {
	if serviceAccountFilePath == "" && email == "" {
		return nil, nil
	}
	if serviceAccountFilePath == "" || email == "" {
		return nil, fmt.Errorf("directory service requires both serviceAccountFilePath and adminEmail")
	}
	jsonCredentials, err := ioutil.ReadFile(serviceAccountFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading credentials from file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Impersonate an admin. This is mandatory for the admin APIs.
	config.Subject = email

	ctx := context.Background()
	client := config.Client(ctx)

	svc, err := admin.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to create directory service %v", err)
	}
	return svc, nil
}

func NewGoogleMethod(config *settings.Config) (groups.GroupsMethod, error) {
	log.Info("Intializing google directory service")
	svc, err := createDirectoryService(config.GoogleServiceAccountFilePath, config.GoogleAdminEmail)
	if err != nil {
		return nil, fmt.Errorf("failed intializing google directory service: %v", err)
	}
	return &googleMethod{adminSvc: svc}, nil
}

func (g *googleMethod) GetGroups(claims map[string]interface{}) ([]string, error) {
	rawEmail, ok := claims[emailClaim]
	if !ok {
		return nil, fmt.Errorf("did not find email claim")
	}
	email, ok := rawEmail.(string)
	if !ok {
		return nil, fmt.Errorf("error getting email claim value")
	}
	var userGroups []string
	var err error
	groupsList := &admin.Groups{}
	for {
		groupsList, err = g.adminSvc.Groups.List().
			UserKey(email).PageToken(groupsList.NextPageToken).Do()
		if err != nil {
			return nil, fmt.Errorf("could not list groups: %v", err)
		}

		log.Printf("email=%s, groups=%v", email, groupsList.Groups)
		for _, group := range groupsList.Groups {
			userGroups = append(userGroups, group.Email)
		}

		if groupsList.NextPageToken == "" {
			break
		}
	}

	return userGroups, nil
}
