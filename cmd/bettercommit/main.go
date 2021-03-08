package main

import (
	"os"
	"os/signal"

	"github.com/gobc/internal/cfg"
	"github.com/gobc/internal/jira"
	"github.com/gobc/internal/scm"
	"github.com/gobc/internal/tui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {

	if _, err := os.Stat(".git"); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("Directory is no git repository: %s", err)
		}
	}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln("Viper error: ", err)
	}

	log.Info("bettercommit starting ...")

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, os.Kill)

	go func() {
		p := cfg.LoadCfg()
		var jiraNumbers []string
		if p.Jira.Enable {
			issues := jira.Issues(p)
			for _, v := range issues.Issues {
				jiraNumbers = append(jiraNumbers, v.Key+" ("+v.Fields.Summary+")")
			}
		}
		stagedFiles, err := scm.StagedFiles()
		if err != nil {
			log.Fatalf("Failed to accept incoming requests: %+v", err)
		}
		var to tui.Options
		_, err = tui.Run(&to, stagedFiles, jiraNumbers)
		if err != nil {
			log.Fatalf("Failed to accept incoming requests: %+v", err)
		}
		log.Debugf("test commit msg: %s", to.CommitMsg)
		if len(to.CommitMsg) >= 9 {
			scm.Commit(to.CommitMsg)

			if to.Push {
				scm.Push()
			}
		} else {
			log.Infoln("commit skipped")
		}
		os.Exit(0)
	}()

	<-shutdown

	log.Info("Initiate graceful shutdown here")
	os.Exit(0)
}
