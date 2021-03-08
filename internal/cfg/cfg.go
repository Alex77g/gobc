package cfg

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/gobc/internal/scm"
)

type Parameter struct {
	Github struct {
		Enable bool   `yaml:"enable"`
		Auth   string `yaml:"auth"`
		User   string `yaml:"user"`
	} `yaml:"github"`
	Gitlab struct {
		Enable bool   `yaml:"enable"`
		Auth   string `yaml:"auth"`
		User   string `yaml:"user"`
	} `yaml:"gitlab"`
	Gitflow struct {
		Enable          bool `yaml:"enable"`
		VersionHandling struct {
			Enable bool   `yaml:"enable"`
			Tag    string `yaml:"tag"`
		} `yaml:"versionHandling"`
	} `yaml:"gitflow"`
	Jira struct {
		Enable bool   `yaml:"enable"`
		URL    string `yaml:"url"`
		Issue  struct {
			User   string   `yaml:"user"`
			Status []string `yaml:"status"`
		} `yaml:"issue"`
		Auth struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"jira"`
	Bitbucket struct {
		Enable bool `yaml:"enable"`
	} `yaml:"bitbucket"`
	Githooks struct {
		Enable bool `yaml:"enable"`
	} `yaml:"githooks"`
	Timetracker struct {
		Enable bool `yaml:"enable"`
	} `yaml:"timetracker"`
}

func LoadCfg() Parameter {

	o := Parameter{}
	path, _ := scm.GitRoot()

	if _, err := os.Stat(path + ".betterconfig.yml"); err == nil {
		log.Debugln(".betterconfig.yml existiert")
	} else if os.IsNotExist(err) {
		log.Debugln(".betterconfig.yml existiert nicht")
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Debugln("existiert oder ?")
	}

	file, err := ioutil.ReadFile(".betterconfig.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(file, &o)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Debugf("cfg line 47 %+v", o)

	return o
}
