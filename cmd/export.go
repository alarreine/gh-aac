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

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/schollz/progressbar/v3"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

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
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
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

// PermissionInfo represents permissions for a repository.
type PermissionInfo struct {
	Repo   string `yaml:"repo"`
	Access string `yaml:"access"`
	Name   string `yaml:"name"`
	IsTeam bool   `yaml:"isTeam"`
}

// AccessConfig represents the overall structure of access-config.yaml.
type AccessConfig struct {
	Organization OrganizationInfo `yaml:"organization"`
	Repositories []RepositoryInfo `yaml:"repositories"`
	Teams        []TeamInfo       `yaml:"teams"`
	Members      []MemberInfo     `yaml:"members"`
	Permissions  []PermissionInfo `yaml:"permissions"`
}

var (
	allowedOrgs []string
	org         *string
	aacPath     *string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the current access configuration from GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		allowedOrgs = viper.GetStringSlice("organizations")

		if *org != "" {
			// Check if the specified organization is allowed
			allowed := false
			for _, allowedOrg := range allowedOrgs {
				if *org == allowedOrg {
					allowed = true
					break
				}
			}
			if !allowed {
				log.Panic("The specified organization is not allowed.")
			} else {
				allowedOrgs = nil
				allowedOrgs = append(allowedOrgs, *org)
			}

		} else {
			// Export all allowed organizations
			log.Printf("Exporting configurations for all allowed organizations:")

		}
		exportConfig()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	org = exportCmd.Flags().StringP("org", "o", "", "Name of the GitHub organization")
	aacPath = exportCmd.Flags().StringP("file", "f", "access-config.yaml", "Path to the output YAML file")

	exportCmd.MarkFlagFilename("file", "yaml")

}

func exportConfig() {

	var (
		accessConfig AccessConfig
	)

	orgProgress := progressbar.Default(int64(len(allowedOrgs)), "Starting..")
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("token")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	// client := githubv4.NewEnterpriseClient("viper.", httpClient)
	client := githubv4.NewClient(httpClient)

	ctx := context.Background()

	for _, allowedOrg := range allowedOrgs {
		orgInfo, err := getOrganizationInfo(ctx, client, allowedOrg)
		if err != nil {
			log.Printf("Failed to get organization info for %s: %v\n", allowedOrg, err)
			continue // Continúa con la siguiente organización si hay un error
		}

		repoInfo, err := getRepos(ctx, client, allowedOrg)
		if err != nil {
			log.Printf("Failed to get organization info for %s: %v\n", allowedOrg, err)
			continue // Continúa con la siguiente organización si hay un error
		}

		teamInfo, err := getTeams(ctx, client, allowedOrg)
		if err != nil {
			log.Printf("Failed to get organization info for %s: %v\n", allowedOrg, err)
			continue // Continúa con la siguiente organización si hay un error
		}
		memberInfo, err := getMembers(ctx, client, allowedOrg)
		if err != nil {
			log.Printf("Failed to get organization info for %s: %v\n", allowedOrg, err)
			continue // Continúa con la siguiente organización si hay un error
		}

		accessConfig.Organization = orgInfo
		accessConfig.Repositories = repoInfo
		accessConfig.Teams = teamInfo
		accessConfig.Members = memberInfo

		permissionInfo, err := getPermissions(ctx, client, allowedOrg, teamInfo, &repoInfo)
		if err != nil {
			log.Printf("Failed to get organization info for %s: %v\n", allowedOrg, err)
			continue
		}
		accessConfig.Permissions = permissionInfo
		orgProgress.Add(1)
	}

	err := SaveConfig("access-config.yaml", &accessConfig)
	if err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Println("Access configuration exported successfully to", *aacPath)
}

func getOrganizationInfo(ctx context.Context, client *githubv4.Client, organization string) (OrganizationInfo, error) {
	var query struct {
		Organization struct {
			ID          githubv4.String
			Name        githubv4.String
			Login       githubv4.String
			Description githubv4.String
			Url         githubv4.String
		} `graphql:"organization(login: $org)"`
	}

	// Definir las variables para la consulta
	variables := map[string]interface{}{
		"org": githubv4.String(organization),
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return OrganizationInfo{}, fmt.Errorf("error ejecutando la consulta: %v", err)
	}

	orgInfo := OrganizationInfo{
		ID:          string(query.Organization.ID),
		Name:        string(query.Organization.Name),
		Login:       string(query.Organization.Login),
		Description: string(query.Organization.Description),
		URL:         string(query.Organization.Url),
	}

	return orgInfo, nil
}

func getRepos(ctx context.Context, client *githubv4.Client, organization string) ([]RepositoryInfo, error) {
	var allRepos []RepositoryInfo
	var afterCursor *githubv4.String

	for {
		var query struct {
			Organization struct {
				Repositories struct {
					Edges []struct {
						Node struct {
							Name githubv4.String
							URL  githubv4.String
						}
					}
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"repositories(first: 100, after: $afterCursor)"`
			} `graphql:"organization(login: $org)"`
		}

		variables := map[string]interface{}{
			"org":         githubv4.String(organization),
			"afterCursor": afterCursor,
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
		}

		for _, edge := range query.Organization.Repositories.Edges {
			repoInfo := RepositoryInfo{
				Name: string(edge.Node.Name),
				URL:  string(edge.Node.URL),
			}
			allRepos = append(allRepos, repoInfo)
		}

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}

		afterCursor = &query.Organization.Repositories.PageInfo.EndCursor
	}

	return allRepos, nil
}

func getTeams(ctx context.Context, client *githubv4.Client, organization string) ([]TeamInfo, error) {
	var allTeams []TeamInfo
	var teamCursor *githubv4.String
	var memberCursor *githubv4.String
	var childTeamCursor *githubv4.String

TeamLoop:
	for {
		var query struct {
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
							} `graphql:"members(first: 100, after: $memberCursor)"`
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
							} `graphql:"childTeams(first: 100, after: $childTeamCursor)"`
						}
					}
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"teams(first: 100, after: $teamCursor)"`
			} `graphql:"organization(login: $org)"`
		}

		variables := map[string]interface{}{
			"org":             githubv4.String(organization),
			"teamCursor":      teamCursor,
			"childTeamCursor": childTeamCursor,
			"memberCursor":    memberCursor,
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
		}

		for _, teams := range query.Organization.Teams.Edges {

			teamInfo := TeamInfo{
				Name:        string(teams.Node.Name),
				Description: string(teams.Node.Description),
			}

		MemberLoop:
			for {

				for _, members := range teams.Node.Members.Edges {
					teamInfo.Members = append(teamInfo.Members, string(members.Node.Login))
				}
				if !teams.Node.Members.PageInfo.HasNextPage {
					break MemberLoop
				}
				memberCursor = &teams.Node.Members.PageInfo.EndCursor

				variables := map[string]interface{}{
					"org":             githubv4.String(organization),
					"teamCursor":      teamCursor,
					"childTeamCursor": childTeamCursor,
					"memberCursor":    memberCursor,
				}

				err := client.Query(ctx, &query, variables)
				if err != nil {
					return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
				}
			}
		ChildTeamLoop:
			for {

				for _, childTeam := range teams.Node.ChildTeams.Edges {
					teamInfo.ChildTeams = append(teamInfo.Members, string(childTeam.Node.Name))
				}

				if !teams.Node.ChildTeams.PageInfo.HasNextPage {
					break ChildTeamLoop
				}
				memberCursor = &teams.Node.Members.PageInfo.EndCursor

				variables := map[string]interface{}{
					"org":             githubv4.String(organization),
					"teamCursor":      teamCursor,
					"childTeamCursor": childTeamCursor,
					"memberCursor":    memberCursor,
				}

				err := client.Query(ctx, &query, variables)
				if err != nil {
					return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
				}
			}

			allTeams = append(allTeams, teamInfo)

		}

		if !query.Organization.Teams.PageInfo.HasNextPage {
			break TeamLoop
		}

		teamCursor = &query.Organization.Teams.PageInfo.EndCursor
	}

	return allTeams, nil
}

func getMembers(ctx context.Context, client *githubv4.Client, organization string) ([]MemberInfo, error) {
	var allMembers []MemberInfo
	var afterCursor *githubv4.String

	for {
		var query struct {
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

		variables := map[string]interface{}{
			"org":         githubv4.String(organization),
			"afterCursor": afterCursor,
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
		}

		for _, role := range query.Organization.MembersWithRole.Edges {

			memberInfo := MemberInfo{
				Role:  string(role.Role),
				Login: string(role.Node.Login),
			}
			allMembers = append(allMembers, memberInfo)
		}

		if !query.Organization.MembersWithRole.PageInfo.HasNextPage {
			break
		}

		afterCursor = &query.Organization.MembersWithRole.PageInfo.EndCursor
	}

	return allMembers, nil
}

func getPermissions(ctx context.Context, client *githubv4.Client, organization string, teams []TeamInfo, repos *[]RepositoryInfo) ([]PermissionInfo, error) {
	var allPermissions []PermissionInfo
	var repoCursor *githubv4.String
	var collabCursor *githubv4.String

TeamPermissions:
	for _, slug := range teams {

		var query struct {
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
					} `graphql:"repositories(first: 100, after: $repoCursor)"`
				} `graphql:"team(slug: $slug)"`
			} `graphql:"organization(login: $org)"`
		}
		variables := map[string]interface{}{
			"org":        githubv4.String(organization),
			"repoCursor": repoCursor,
			"slug":       githubv4.String(slug.Name),
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
		}

		for i, permission := range query.Organization.Team.Repositories.Nodes {
			teamPermissionInfo := PermissionInfo{
				Repo:   string(permission.Name),
				Access: string(query.Organization.Team.Repositories.Edges[i].Permission),
				IsTeam: true,
				Name:   slug.Name,
			}

			allPermissions = append(allPermissions, teamPermissionInfo)
		}

		if !query.Organization.Team.Repositories.PageInfo.HasNextPage {
			break TeamPermissions
		}

		repoCursor = &query.Organization.Team.Repositories.PageInfo.EndCursor
	}

RepoPermissions:
	for {

		var query struct {
			Organization struct {
				Repositories struct {
					Edges []struct {
						Node struct {
							Name          githubv4.String
							Collaborators []struct {
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
				} `graphql:"repositories(first: 100, after: $repoCursor)"`
			} `graphql:"organization(login: $org)"`
		}
		variables := map[string]interface{}{
			"org":          githubv4.String(organization),
			"repoCursor":   repoCursor,
			"collabCursor": collabCursor,
		}

		err := client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("error ejecutando la consulta: %v", err)
		}

		for _, repo := range query.Organization.Repositories.Edges {

			// teamPermissionInfo := PermissionInfo{
			// 	Repo:   string(permission.Name),
			// 	Access: string(query.Organization.Team.Repositories.Edges[i].Permission),
			// 	IsTeam: true,
			// 	Name:   slug.Name,
			// }

			// allPermissions = append(allPermissions, teamPermissionInfo)
		RepoMemberLoop:
			for {

			}
		}

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break RepoPermissions
		}

		repoCursor = &query.Organization.Repositories.PageInfo.EndCursor
	}

	return allPermissions, nil
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(filename string) (*AccessConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config AccessConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveConfig saves the configuration to a YAML file.
func SaveConfig(filename string, config *AccessConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
