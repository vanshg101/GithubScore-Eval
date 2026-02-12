package github

import (
	"encoding/json"
	"fmt"
)

// OrgMember represents a GitHub organization member.
type OrgMember struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

// FetchOrgMembers returns all public members of a GitHub organization.
func (c *Client) FetchOrgMembers(org string) ([]OrgMember, error) {
	url := fmt.Sprintf("%s/orgs/%s/members?per_page=%d", baseURL, org, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching org members for %s: %w", org, err)
	}

	var members []OrgMember
	for _, raw := range rawItems {
		var m OrgMember
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("decoding org member: %w", err)
		}
		members = append(members, m)
	}

	return members, nil
}
