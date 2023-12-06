package impl

type HumanImpl2 struct{}

// Returns the name of our HumanImpl.
func (_ *HumanImpl2) GetName() string {
	panic("not implemented") // TODO: Implement
}

// Our HumanImpl just had a birthday! Increase its age.
func (_ *HumanImpl2) Birthday() {
	panic("not implemented") // TODO: Implement
}

// Make the HumanImpl say hello.
func (_ *HumanImpl2) SayHello() {
	panic("not implemented") // TODO: Implement
}

var human Human

func init() {
	human = HumanImpl2{}
}
