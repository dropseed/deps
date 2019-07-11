package runner

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Local runs a full interactive update process
func Local() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	allUpdates, err := collectUpdates(cfg, []string{})
	if err != nil {
		return err
	}

	newUpdates, _, _, err := organizeUpdates(allUpdates)
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
		refs := map[int]string{}
		items := []string{}

		updateIndex := 0
		for _, update := range updates {
			if !update.completed {
				items = append(items, update.title)
				refs[updateIndex] = update.id
				updateIndex++
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

		if i+1 == len(items) {
			// Chose to skip
			break
		}

		update := updates[refs[i]]
		if _, err := update.runner.Act(update.dependencies); err != nil {
			return err
		}
		update.completed = true
	}

	return nil
}
