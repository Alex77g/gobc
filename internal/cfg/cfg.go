package cfg

import (
	"log"
	"os"

	"github.com/gobc/internal/scm"
)

func LoadCfg() {
	path, _ := scm.GitRoot()

	if _, err := os.Stat(path + "/.betterconfig"); err == nil {
		log.Println("existiert")
	} else if os.IsNotExist(err) {
		log.Println("existiert nicht")
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Println("existiert oder ?")
	}
}
