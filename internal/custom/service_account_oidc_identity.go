package custom

import (
	"fmt"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/services"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/services/api"
)

type OIDCIdentity struct {
	ID               *string `json:"Id,omitempty"`
	ServiceAccountID string  `json:"ServiceAccountId"`
	Name             string  `json:"Name"`
	Issuer           string  `json:"Issuer"`
	Subject          string  `json:"Subject"`
}

type GetServiceAccountOIDCIdentityResponse OIDCIdentity

type ListServiceAccountOIDCIdentitesResponse struct {
	ServerURL      string         `json:"ServerUrl"`
	ExternalID     string         `json:"ExternalId"`
	Count          int64          `json:"Count"`
	OIDCIdentities []OIDCIdentity `json:"OidcIdentities"`
}

type CreateServiceAccountOIDCIdentityResponse struct {
	ID string `json:"Id"`
}

type UpdateServiceAccountOIDCIdentityResponse CreateServiceAccountOIDCIdentityResponse

type DeleteServiceAccountOIDCIdentityResponse struct{}

func (c *Client) ListServiceAccountOIDCIdentites(serviceAccountID string, skip, take int) (*ListServiceAccountOIDCIdentitesResponse, error) {
	query := url.Values{"skip": {fmt.Sprint(skip)}, "take": {fmt.Sprint(take)}}.Encode()
	path := fmt.Sprintf("serviceaccounts/%s/oidcidentities/v1?%s", serviceAccountID, query)
	res, err := api.ApiGet(c.client.Sling(), new(ListServiceAccountOIDCIdentitesResponse), path)
	if err != nil {
		return nil, err
	}

	return res.(*ListServiceAccountOIDCIdentitesResponse), nil
}

func (c *Client) GetServiceAccountOIDCIdentity(serviceAccountID, identityID string) (*GetServiceAccountOIDCIdentityResponse, error) {
	path := fmt.Sprintf("serviceaccounts/%s/oidcidentities/%s/v1", serviceAccountID, identityID)
	res, err := api.ApiGet(c.client.Sling(), new(GetServiceAccountOIDCIdentityResponse), path)
	if err != nil {
		return nil, err
	}

	return res.(*GetServiceAccountOIDCIdentityResponse), nil
}

func (c *Client) CreateServiceAccountOIDCIdentity(identity OIDCIdentity) (*CreateServiceAccountOIDCIdentityResponse, error) {
	path := fmt.Sprintf("serviceaccounts/%s/oidcidentities/v1", identity.ServiceAccountID)
	res, err := services.ApiAdd(c.client.Sling(), identity, new(CreateServiceAccountOIDCIdentityResponse), path)
	if err != nil {
		return nil, err
	}

	return res.(*CreateServiceAccountOIDCIdentityResponse), nil
}

func (c *Client) UpdateServiceAccountOIDCIdentity(identity OIDCIdentity) (*UpdateServiceAccountOIDCIdentityResponse, error) {
	path := fmt.Sprintf("serviceaccounts/%s/oidcidentities/v1", identity.ServiceAccountID)
	res, err := services.ApiPost(c.client.Sling(), identity, new(UpdateServiceAccountOIDCIdentityResponse), path)
	if err != nil {
		return nil, err
	}

	return res.(*UpdateServiceAccountOIDCIdentityResponse), nil
}

func (c *Client) DeleteServiceAccountOIDCIdentity(serviceAccountID, identityID string) error {
	path := fmt.Sprintf("serviceaccounts/%s/oidcidentities/%s/v1", serviceAccountID, identityID)
	_, err := c.client.Sling().Post(path).Receive(new(DeleteServiceAccountOIDCIdentityResponse), new(core.APIError))
	return err
}
