package generic

type GenericCI struct {
}

func (generic *GenericCI) Autoconfigure() error {
	return nil
}

func (generic *GenericCI) Branch() string {
	return ""
}
