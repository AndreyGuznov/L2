package main

import (
	"fmt"
)

// Паттерн Фабричный метод является пораждающим паттерном. Он определяет общий интерфейс поведения
// для создаваемых объектов. Решает проблему дополнения бизнес-логики.

// Плюсы:
// 1. Избавляет от привязки к конкретному типу объекта при помощи конструктора;
// 2. Упрощает добавление новых объектов от базового класса;
// 3. Реализует принцип open-closed.

// Минусы:
// 1. Может привести к созданию больших иерархий объектов (большое количество структур, которые сложны при сопровождении);
// 2. Появляется один главный конструктор, к которому будет привязана вся логика программы.

const (
	WorkComp   = "working"
	PersonComp = "personal"
)

func New(typeName string) Computer {
	switch typeName {
	case WorkComp:
		return NewWork()
	case PersonComp:
		return NewForMe()
	default:
		fmt.Printf("%s is unknown\n", typeName)
		return nil
	}
}

type Computer interface {
	Info()
	GetType() string
}

type ForWork struct {
	Type   string
	Core   int
	Memory int
}

// Создание базового конструктора для фабричного метода
func NewWork() Computer {
	return ForWork{
		Type:   WorkComp,
		Core:   4,
		Memory: 256,
	}
}

func (o ForWork) GetType() string {
	return o.Type
}

func (o ForWork) Info() {
	fmt.Printf("%s Core: [%d] Memory: [%d]\n", o.Type, o.Core, o.Memory)
}

type MyOwn struct {
	Type        string
	Core        int
	Memory      int
	GraphicCard bool
}

func NewForMe() Computer {
	return MyOwn{
		Type:        PersonComp,
		Core:        8,
		Memory:      16,
		GraphicCard: true,
	}
}

func (m MyOwn) GetType() string {
	return m.Type
}

func (pc MyOwn) Info() {
	fmt.Printf("%s Core: [%d] Memory: [%d] GraphicCard: [%v]\n", pc.Type, pc.Core, pc.Memory, pc.GraphicCard)
}
