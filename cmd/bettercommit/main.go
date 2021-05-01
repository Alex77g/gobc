package main

import (
	"os"

	"github.com/gobc/internal/cfg"
	"github.com/gobc/internal/jira"
	"github.com/gobc/internal/scm"
	"github.com/gobc/internal/tui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	jirauser := os.Getenv("JIRA_USER")
	jiratoken := os.Getenv("JIRA_TOKEN")

	if _, err := os.Stat(".git"); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("Directory is no git repository: %s", err)
		}
	}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Viper error: %s \n", err)
	}

	log.Info("bettercommit starting ...")

	p := cfg.LoadCfg()
	var jiraNumbers []string
	if p.Jira.Enable {
		issues := jira.Issues(p, jirauser, jiratoken)
		for _, v := range issues.Issues {
			jiraNumbers = append(jiraNumbers, v.Key+" ("+v.Fields.Summary+")")
		}
	}
	stagedFiles, err := scm.StagedFiles()
	if err != nil {
		log.Fatalf("No staged files given: %+v", err)
	}
	var to tui.Options
	_, err = tui.Run(&to, stagedFiles, jiraNumbers)
	if err != nil {
		log.Fatalf("Failed to load prompt menu: %+v", err)
	}

	if len(to.CommitMsg) >= 9 {
		scm.Commit(to.CommitMsg)

		if to.Push {
			scm.Push()
		}
	} else {
		log.Infoln("commit skipped")
	}

	log.Info("Initiate shutdown ........")
	os.Exit(0)
}
