// DONT EDIT: Auto generated

package gen

// Human makes human interaction easy
type Human interface {
	// Returns the name of our HumanImpl.
	GetName() string
	// Our HumanImpl just had a birthday! Increase its age.
	Birthday()
	// Make the HumanImpl say hello.
	SayHello()
}

var human Human

func RefHuman() Human {
	return human
}

//func init() {
//	human = impl.NewHumanImpl()
//}
