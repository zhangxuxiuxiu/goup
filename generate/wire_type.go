//go:build wireinject
// +build wireinject

//https://github.com/google/wire/blob/main/_tutorial/README.md

package gen

import (
	"fmt"

	"github.com/google/wire"
)

type Message string

func NewMessage(phrase string) Message {
	return Message(phrase)
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

type Greeter struct {
	Message Message // <- adding a Message field
}

func (g Greeter) Greet() Message {
	return g.Message
}

func NewEvent(g Greeter) (Event, error) {
	return Event{Greeter: g}, nil
}

type Event struct {
	Greeter Greeter // <- adding a Greeter field
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func InitializeEvent(phrase string) (Event, error) {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}, nil
}
