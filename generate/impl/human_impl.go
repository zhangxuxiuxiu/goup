package impl

import "fmt"

//go:generate ifacemaker -f human_impl.go -s HumanImpl -i Human -p gen -y "Human makes human interaction easy" -c "DONT EDIT: Auto generated" -o ../human.go
type HumanImpl struct {
	name string
	age  int
}

// Returns the name of our HumanImpl.
func (h *HumanImpl) GetName() string {
	return h.name
}

// Our HumanImpl just had a birthday! Increase its age.
func (h *HumanImpl) Birthday() {
	h.age += 1
	fmt.Printf("I am now %d years old!\n", h.age)
}

// Make the HumanImpl say hello.
func (h *HumanImpl) SayHello() {
	fmt.Printf("Hello, my name is %s, and I am %d years old.\n", h.name, h.age)
}

//go:generate bash inject.sh Human "impl.NewHumanImpl()" ../human.go
func NewHumanImpl() *HumanImpl {
	return &HumanImpl{}
}

//func main() {
//	human := &HumanImpl{name: "Bob", age: 30}
//	human.GetName()
//	human.SayHello()
//	human.Birthday()
//}
