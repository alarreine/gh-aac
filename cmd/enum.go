/*
Copyright © 2023 Agustin Larreinegabe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import "github.com/shurcooL/githubv4"

// OrganizationInfo represents basic information about an organization.
type OrganizationInfo struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Login       string `yaml:"login"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

// RepositoryInfo represents basic information about a repository.
type RepositoryInfo struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	Visibility string `yaml:"visibility"`
}

// TeamInfo represents basic information about a team.
type TeamInfo struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Members     []string `yaml:"members"`
	ChildTeams  []string `yaml:"childTeam"`
}

// MemberInfo represents basic information about a member.
type MemberInfo struct {
	Login string `yaml:"login"`
	Role  string `yaml:"role"`
}

// type TeamPermission struct {
// 	// Repo   string `yaml:"repo"`
// 	Slug   string `yaml:"slug"`
// 	Access string `yaml:"access"`
// }

type AccessPermission struct {
	// Repo   string `yaml:"repo"`
	Member string `yaml:"member"`
	Access string `yaml:"access"`
}

type AccessInfo struct {
	Teams         []AccessPermission `yaml:"teams"`
	Collaborators []AccessPermission `yaml:"collaborators"`
}

type PermissionsInfo struct {
	RepoName string     `yaml:"repoName"`
	Access   AccessInfo `yaml:"teams"`
}

type MembershipInfo struct {
	Teams        []TeamInfo   `yaml:"teams"`
	Organization []MemberInfo `yaml:"organization"`
}

// AccessConfig represents the overall structure of access-config.yaml.
type AccessConfig struct {
	Organization OrganizationInfo  `yaml:"organization"`
	Repositories []RepositoryInfo  `yaml:"repositories"`
	Memberships  MembershipInfo    `yaml:"memberships"` // List of members with admin access
	Permissions  []PermissionsInfo `yaml:"permissions"`
}

type OrganizationQuery struct {
	Organization struct {
		ID          githubv4.String
		Name        githubv4.String
		Login       githubv4.String
		Description githubv4.String
		Url         githubv4.String
	} `graphql:"organization(login: $org)"`
}

type RepoQuery struct {
	Organization struct {
		Repositories struct {
			Edges []struct {
				Node struct {
					Name       githubv4.String
					URL        githubv4.String
					Visibility githubv4.String
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: 100, after: $afterCursor, orderBy: {field:NAME,direction:ASC})"`
	} `graphql:"organization(login: $org)"`
}

type TeamQuery struct {
	Organization struct {
		Teams struct {
			Edges []struct {
				Node struct {
					Name        githubv4.String
					Description githubv4.String
					Members     struct {
						Edges []struct {
							Node struct {
								Login githubv4.String
							}
						}
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"members(first: 100, after: $memberCursor, orderBy:{field:LOGIN, direction:ASC})"`
					ChildTeams struct {
						Edges []struct {
							Node struct {
								Name githubv4.String
							}
						}
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"childTeams(first: 100, after: $childTeamCursor, orderBy:{field:NAME, direction:ASC})"`
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"teams(first: 100, after: $teamCursor, orderBy:{field:NAME, direction:ASC})"`
	} `graphql:"organization(login: $org)"`
}

type MemberQuery struct {
	Organization struct {
		MembersWithRole struct {
			Edges []struct {
				Role githubv4.String
				Node struct {
					Login githubv4.String
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"membersWithRole(first: 100, after: $afterCursor)"`
	} `graphql:"organization(login: $org)"`
}

type TeamPermissionQuery struct {
	Organization struct {
		Team struct {
			Repositories struct {
				Nodes []struct {
					Name githubv4.String
				}
				Edges []struct {
					Permission githubv4.String
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $repoCursor, orderBy: {field:CREATED_AT,direction:ASC})"`
		} `graphql:"team(slug: $slug)"`
	} `graphql:"organization(login: $org)"`
}

type RepoPermissionQuery struct {
	Organization struct {
		Repositories struct {
			Edges []struct {
				Node struct {
					Name          githubv4.String
					Collaborators struct {
						Edges []struct {
							Node struct {
								Login githubv4.String
							}
							Permission githubv4.String
						}
						PageInfo struct {
							EndCursor   githubv4.String
							HasNextPage bool
						}
					} `graphql:"collaborators(first: 100, after: $collabCursor)"`
				}
			}
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"repositories(first: 100, after: $repoCursor, orderBy: {field:CREATED_AT,direction:ASC})"`
	} `graphql:"organization(login: $org)"`
}
