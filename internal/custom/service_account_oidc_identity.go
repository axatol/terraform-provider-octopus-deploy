package custom

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrUnrecognisedResponse = errors.New("unrecognised response")
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

type UpdateServiceAccountOIDCIdentityResponse struct{}

type DeleteServiceAccountOIDCIdentityResponse struct{}

func (c *Client) ListServiceAccountOIDCIdentites(serviceAccountID string, skip, take int) (res ListServiceAccountOIDCIdentitesResponse, err error) {
	query := url.Values{"skip": {fmt.Sprint(skip)}, "take": {fmt.Sprint(take)}}
	endpoint := fmt.Sprintf("serviceaccounts/%s/oidcidentities/v1?%s", serviceAccountID, query.Encode())
	err = c.do(c.client.Sling().New().Get(endpoint), &res)
	return res, err
}

func (c *Client) GetServiceAccountOIDCIdentity(serviceAccountID, identityID string) (res GetServiceAccountOIDCIdentityResponse, err error) {
	endpoint := fmt.Sprintf("serviceaccounts/%s/oidcidentities/%s/v1", serviceAccountID, identityID)
	err = c.do(c.client.Sling().New().Get(endpoint), &res)
	return res, err
}

func (c *Client) CreateServiceAccountOIDCIdentity(identity OIDCIdentity) (res CreateServiceAccountOIDCIdentityResponse, err error) {
	endpoint := fmt.Sprintf("serviceaccounts/%s/oidcidentities/create/v1", identity.ServiceAccountID)
	err = c.do(c.client.Sling().New().Post(endpoint).BodyJSON(identity), &res)
	return res, err
}

func (c *Client) UpdateServiceAccountOIDCIdentity(identity OIDCIdentity) (res UpdateServiceAccountOIDCIdentityResponse, err error) {
	endpoint := fmt.Sprintf("serviceaccounts/%s/oidcidentities/%s/v1", identity.ServiceAccountID, *identity.ID)
	err = c.do(c.client.Sling().New().Put(endpoint).BodyJSON(identity), &res)
	return res, err
}

func (c *Client) DeleteServiceAccountOIDCIdentity(serviceAccountID, identityID string) (res DeleteServiceAccountOIDCIdentityResponse, err error) {
	endpoint := fmt.Sprintf("serviceaccounts/%s/oidcidentities/%s/v1", serviceAccountID, identityID)
	err = c.do(c.client.Sling().New().Delete(endpoint), &res)
	return res, err
}
