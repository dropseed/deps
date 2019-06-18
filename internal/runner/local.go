package runner

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Local runs a full interactive update process
func Local() error {
	newUpdates, _, _, err := collectUpdates(-1)
	if err != nil {
		return err
	}

	if err := newUpdates.prompt(); err != nil {
		return err
	}

	return nil
}

func (updates Updates) prompt() error {
	for {
		items := []string{}
		for _, update := range updates {
			if !update.completed {
				title, err := update.dependencies.GenerateTitle()
				if err != nil {
					return err
				}
				items = append(items, title)
			}
		}

		if len(items) < 1 {
			// No updates left
			break
		}

		items = append(items, "Skip")

		prompt := promptui.Select{
			Label: fmt.Sprintf("Choose an update to make"),
			Items: items,
		}

		i, _, err := prompt.Run()
		if err != nil {
			return err
		}

		if i < len(updates) {
			update := updates[i]
			if err := update.runner.Act(update.dependencies, "", false); err != nil {
				return err
			}
			update.completed = true
		} else {
			// Chose skip
			break
		}
	}

	return nil
}
