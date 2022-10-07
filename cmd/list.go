package cmd

import (
	"fmt"

	"github.com/cli/go-gh"
	"github.com/tidwall/gjson"
)

func List(owner string, name string) {
	args := []string{
		"api", "graphql",
		"-F", fmt.Sprintf("owner=%s", owner),
		"-F", fmt.Sprintf("name=%s", name),
		"-f", fmt.Sprintf("query=%s", `query($name: String!, $owner: String!) {
			repository(name: $name, owner: $owner) {
				name
				pushedAt
				refs(first: 100, refPrefix: "refs/heads/") {
					nodes {
						name
						prefix
						target {
							... on Commit {
								history(first: 1) {
									nodes {
										message
										author {
											name
										}
										committedDate
									}
								}
							}
						}
					}
          totalCount
          pageInfo {
            hasNextPage
          }
				}
			}
		}`),
	}

	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		return
	}

	str := stdOut.String()
	commitMessage := gjson.Get(str, "data.repository.refs.nodes.0.target.history.nodes.0.message")
	authorName := gjson.Get(str, "data.repository.refs.nodes.0.target.history.nodes.0.author.name")
	committedDate := gjson.Get(str, "data.repository.refs.nodes.0.target.history.nodes.0.committedDate")
	fmt.Printf("%s %s %s\n", committedDate, authorName, commitMessage)
}