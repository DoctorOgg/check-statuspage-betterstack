package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

type StatusPage struct {
	Page struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		URL       string    `json:"url"`
		TimeZone  string    `json:"time_zone"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"page"`
	Incidents []struct {
		CreatedAt       string `json:"created_at"`
		ID              string `json:"id"`
		Impact          string `json:"impact"`
		IncidentUpdates []struct {
			Body       string `json:"body"`
			CreatedAt  string `json:"created_at"`
			DisplayAt  string `json:"display_at"`
			ID         string `json:"id"`
			IncidentID string `json:"incident_id"`
			Status     string `json:"status"`
			UpdatedAt  string `json:"updated_at"`
		} `json:"incident_updates"`
		MonitoringAt interface{} `json:"monitoring_at"`
		Name         string      `json:"name"`
		PageID       string      `json:"page_id"`
		ResolvedAt   interface{} `json:"resolved_at"`
		Shortlink    string      `json:"shortlink"`
		Status       string      `json:"status"`
		UpdatedAt    string      `json:"updated_at"`
	} `json:"incidents"`
}

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Url string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-check-statuspage",
			Short:    "Check for statuspage.com",
			Keyspace: "sensu.io/plugins/sensu-check-statuspage/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "url",
			Env:       "STATUS_PAGE_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "",
			Usage:     "URL of Statuspage to Monitor",
			Value:     &plugin.Url,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Url) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--url or STATUS_PAGE_URL environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	res, err := http.Get(plugin.Url)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("error fetching status page: %v", err)
	}
	defer res.Body.Close()

	issues, err := parseInstatusHTML(res.Body)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("error parsing HTML: %v", err)
	}

	if len(issues) > 0 {
		fmt.Println("Non-operational components:")
		for _, issue := range issues {
			fmt.Printf(" - %s\n", issue)
		}
		return sensu.CheckStateCritical, nil
	}

	fmt.Println("All components operational.")
	return sensu.CheckStateOK, nil
}

func parseInstatusHTML(body io.Reader) ([]string, error) {
	z := html.NewTokenizer(body)
	var issues []string
	var insideComponent bool
	var componentName, componentStatus string

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return issues, nil
			}
			return nil, z.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "div" {
				var class string
				for _, a := range t.Attr {
					if a.Key == "class" {
						class = a.Val
						break
					}
				}
				if strings.Contains(class, "component") {
					insideComponent = true
				} else if strings.Contains(class, "component-name") {
					z.Next()
					componentName = strings.TrimSpace(html.UnescapeString(string(z.Text())))
				} else if strings.Contains(class, "component-status") {
					z.Next()
					componentStatus = strings.TrimSpace(html.UnescapeString(string(z.Text())))
					if componentStatus != "Operational" {
						issues = append(issues, fmt.Sprintf("%s: %s", componentName, componentStatus))
					}
					componentName, componentStatus = "", "" // reset
				}
			}
		case html.EndTagToken:
			if insideComponent && z.Token().Data == "div" {
				insideComponent = false
			}
		}
	}
}
