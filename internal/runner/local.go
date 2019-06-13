package runner

// Run a full interactive update process
func Local() error {
	newUpdates, _, _, err := collectUpdates(-1)
	if err != nil {
		return err
	}

	if err := newUpdates.Prompt(); err != nil {
		return err
	}

	return nil
}
