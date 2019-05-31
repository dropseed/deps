package versionfilter

func fillSlices(a *[]string, b *[]string) {
	longer := a
	shorter := b
	if len(*b) > len(*a) {
		longer = b
		shorter = a
	}
	toAdd := len(*longer) - len(*shorter)
	for i := 0; i < toAdd; i++ {
		*shorter = append(*shorter, "")
	}
}
