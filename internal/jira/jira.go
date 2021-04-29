package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gobc/internal/cfg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JiraIssue struct {
	Expand     string `json:"expand"`
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []struct {
		Expand string `json:"expand"`
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields struct {
			Summary  string   `json:"summary"`
			Assignee Assignee `json:"assignee"`
			Status   Status   `json:"status"`
		} `json:"fields"`
	} `json:"issues"`
}

type Assignee struct {
	Self       string `json:"self"`
	AccountID  string `json:"accountId"`
	AvatarUrls struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	TimeZone    string `json:"timeZone"`
	AccountType string `json:"accountType"`
}

type Status struct {
	Self           string `json:"self"`
	Description    string `json:"description"`
	IconURL        string `json:"iconUrl"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	StatusCategory struct {
		Self      string `json:"self"`
		ID        int    `json:"id"`
		Key       string `json:"key"`
		ColorName string `json:"colorName"`
		Name      string `json:"name"`
	} `json:"statusCategory"`
}

func init() {

}

func Issues(p cfg.Parameter) JiraIssue {
	var issues JiraIssue
	url := p.Jira.URL + "/rest/api/2/search?jql=project=KUB&fields=summary,assignee,status"
	log.Debugf("%s", url)
	resp := httpReq(nil, url, http.MethodGet, p)

	err := json.Unmarshal(resp, &issues)
	if err != nil {
		log.Errorf("error when unmarshall issue: %s", err)
	}
	sortedIssues := sortIssues(issues, p)

	return sortedIssues
}

func sortIssues(i JiraIssue, p cfg.Parameter) JiraIssue {
	var sortedIssues JiraIssue

	for _, j := range i.Issues {
		if j.Fields.Assignee.DisplayName == p.Jira.Issue.UserName {
			for _, k := range p.Jira.Issue.Status {
				if k == j.Fields.Status.Name {
					sortedIssues.Issues = append(sortedIssues.Issues, j)
				}
			}
		}
	}
	return sortedIssues
}

func httpReq(i interface{}, url, httpMethod string, p cfg.Parameter) []byte {

	const ConnectMaxWaitTime = 1 * time.Second
	const RequestMaxWaitTime = 5 * time.Second

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), RequestMaxWaitTime)
	defer cancel()

	jsn, _ := json.Marshal(i)
	log.Debugln(string(jsn))

	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewBuffer(jsn))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(viper.Get("JIRA_USER").(string), viper.Get("JIRA_TOKEN").(string))

	if err != nil {
		log.Fatalf("Cannot create request: %s\n", err)
	}

	rsp, err := client.Do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if e, ok := err.(net.Error); ok && e.Timeout() {
		log.Panicf("Do request timeout: %s\n", err)
	} else if err != nil {
		log.Panicf("Cannot do request: %s\n", err)
	}

	startRead := time.Now()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Panicf("Cannot read all response body: %s\n", err)
	}
	endRead := time.Now()
	log.Debugf("Read response took: %s\n to url %s", endRead.Sub(startRead), url)

	// log.Debugf("%s\n", string(body))

	return body
}
