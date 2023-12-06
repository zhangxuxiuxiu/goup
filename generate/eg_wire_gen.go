// Code generated by Wire. DO NOT EDIT.

//go:build !wireinject
// +build !wireinject

package gen

import (
	"fmt"
)

// Injectors from wire_type.go:

func InitializeEvent(phrase string) (Event, error) {
	message := NewMessage(phrase)
	greeter := NewGreeter(message)
	event, err := NewEvent(greeter)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

// wire_type.go:

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
