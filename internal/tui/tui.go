package tui

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/manifoldco/promptui"
)

type Options struct {
	CommitMsg string
	Push      bool
}

// Run starts the terminal
func Run(to *Options) (*Options, error) {
	t, _ := commitType()
	s, _ := commitScope()
	d, _ := commitDesc()
	te, _ := commitText()
	to.Push = commitPush()

	to.CommitMsg = t + ":" + s + d + te

	return to, nil
}

func commitType() (string, error) {
	prompt := promptui.Select{
		Label: "Select the type of change that you're committing",
		Items: []string{
			"feat:     A new feature",
			"fix:      A bug fix",
			"docs:     Documentation only changes",
			"style:    Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)",
			"refactor: A code change that neither fixes a bug nor adds a feature",
			"perf:     A code change that improves performance",
			"test:     Adding missing tests or correcting existing tests",
		},
	}

	_, result, err := prompt.Run()

	s := strings.Split(result, ":")

	if err != nil {
		return "", nil
	}

	return s[0], nil
}

func commitScope() (string, error) {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "What is the scope of this change (e.g. component or file name): (press enter to skip)",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		log.Errorf("Prompt failed %v\n", err)
		return "", nil
	}

	if len(result) == 0 {
		return "", nil
	}

	return "(" + result + ")", nil
}

func commitDesc() (string, error) {
	validate := func(input string) error {
		if len(input) > 0 && len(input) >= 95 {
			return errors.New("Message must have less than 95 characters")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Write a short description",
		Validate:  validate,
		Templates: templates,
	}

	result, err := prompt.Run()

	if err != nil {
		log.Errorf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func commitText() (string, error) {

	prompt := promptui.Prompt{
		Label: "Write a longer description",
	}

	result, err := prompt.Run()

	if err != nil {
		log.Errorf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func commitPush() bool {

	prompt := promptui.Prompt{
		Label:     "Push Commit?",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		return false
	}

	var ret bool
	switch result {
	case "y":
		ret = true
	case "n":
		ret = false
	default:
		ret = true
	}

	return ret
}
