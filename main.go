package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if ok := len(os.Args) > 1; !ok {
		fmt.Println("Please provide a project id")
		os.Exit(1)
	}

	gitlabHost := os.Getenv("GITLAB_HOST")
	if gitlabHost == "" {
		fmt.Println("error: GITLAB_HOST environment variable is not set")
		os.Exit(1)
	}

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		fmt.Println("error: GITLAB_TOKEN environment variable is not set")
		os.Exit(1)
	}

	projectId, _ := strconv.Atoi(os.Args[1])
	apiUrl := fmt.Sprintf("https://%s/api/v4/projects/%d/variables", gitlabHost, projectId)
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("PRIVATE-TOKEN", gitlabToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error: Gitlab API returned non-200 status code")
		os.Exit(1)
	}

	body, _ := io.ReadAll(resp.Body)
	var variables []map[string]interface{}
	err = json.Unmarshal(body, &variables)
	if err != nil {
		panic(err)
	}

	varGrouped := make(map[string][]map[string]interface{})
	for _, variable := range variables {
		if scope, ok := variable["environment_scope"]; ok {
			scopeStr, ok := scope.(string)
			if !ok {
				fmt.Println("Error: environment_scope is not a string")
				continue
			}
			varGrouped[scopeStr] = append(varGrouped[scopeStr], variable)
		}
	}

	for scope, vars := range varGrouped {
		fmt.Printf("# %s\n", scope)
		for _, variable := range vars {
			fmt.Printf("export %s=%s\n", variable["key"], strconv.Quote(variable["value"].(string)))
		}
	}
}
