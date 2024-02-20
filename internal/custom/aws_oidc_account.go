package custom

import (
	"context"
	"fmt"
)

type AWSOIDCAccount struct {
	SpaceID                         string   `json:"SpaceId,omitempty"`
	ID                              string   `json:"Id,omitempty"`
	Slug                            string   `json:"Slug,omitempty"`
	Name                            string   `json:"Name"`
	Description                     string   `json:"Description"`
	TenantedDeploymentParticipation string   `json:"TenantedDeploymentParticipation"`
	AccountType                     string   `json:"AccountType"`
	RoleARN                         string   `json:"RoleArn"`
	SessionDuration                 string   `json:"SessionDuration"`
	EnvironmentIDs                  []string `json:"EnvironmentIds"`
	TenantIDs                       []string `json:"TenantIds"`
	TenantTags                      []string `json:"TenantTags"`
	DeploymentSubjectKeys           []string `json:"DeploymentSubjectKeys"`
	HealthCheckSubjectKeys          []string `json:"HealthCheckSubjectKeys"`
	AccountTestSubjectKeys          []string `json:"AccountTestSubjectKeys"`
}

func (c *Client) GetAWSOIDCAccount(ctx context.Context, spaceID, accountID string) (res *AWSOIDCAccount, err error) {
	endpoint := fmt.Sprintf("spaces/%s/accounts/%s", spaceID, accountID)
	err = c.do(ctx, c.client.Sling().New().Get(endpoint), &res)
	return res, err
}

func (c *Client) CreateAWSOIDCAccount(ctx context.Context, account AWSOIDCAccount) (res *AWSOIDCAccount, err error) {
	endpoint := fmt.Sprintf("spaces/%s/accounts", account.SpaceID)
	err = c.do(ctx, c.client.Sling().New().Post(endpoint).BodyJSON(account), &res)
	return res, err
}

func (c *Client) UpdateAWSOIDCAccount(ctx context.Context, account AWSOIDCAccount) (res *AWSOIDCAccount, err error) {
	endpoint := fmt.Sprintf("spaces/%s/accounts/%s", account.SpaceID, account.ID)
	err = c.do(ctx, c.client.Sling().New().Put(endpoint).BodyJSON(account), &res)
	return res, err
}

func (c *Client) DeleteAWSOIDCAccount(ctx context.Context, spaceID, accountID string) error {
	endpoint := fmt.Sprintf("spaces/%s/accounts/%s", spaceID, accountID)
	err := c.do(ctx, c.client.Sling().New().Delete(endpoint), nil)
	return err
}
