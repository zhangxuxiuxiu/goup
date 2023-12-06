package util

// NoCopy : go vet to see if stopped
type NoCopy struct{}

func (*NoCopy) Lock()   {}
func (*NoCopy) Unlock() {}
