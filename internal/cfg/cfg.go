package cfg

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/gobc/internal/scm"
)

type option struct {
	github struct {
		enable bool   `yaml:"enable"`
		auth   string `yaml:"auth"`
	}
	gitlab struct {
		enable bool   `yaml:"enable"`
		auth   string `yaml:"auth"`
	}
}

func LoadCfg() {

	o := option{}
	path, _ := scm.GitRoot()

	if _, err := os.Stat(path + ".betterconfig"); err == nil {
		log.Debugln(".betterconfig existiert")
	} else if os.IsNotExist(err) {
		log.Debugln(".betterconfig existiert nicht")
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
	log.Debugf("%+v", o)
}
