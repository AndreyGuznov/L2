package main

import "fmt"

// Паттерн Посетитель является поведенческим паттерном уровня объекта. Он позволяет обойти
// набор элементов (объектов) с разнородными интерфейсами, а также позволяет добавить
// новый метод в тип объекта, при этом, не изменяя сам тип этого объекта.
// Решает задачу определения новой операции, не изменяя типы объектов, над которыми выполняется одна или более операций.
// Паттерн следует применять если:
// 1. Имеются различные объекты разных типов с разными интерфейсами, но над ними нужно совершать операции, зависящие от конкретных типов;
// 2. Необходимо над структурой выполнить различные, усложняющие структуру операции;
// 3. Часто добавляются новые операции над структурой.

// Плюсы:
// 1. Упрощается добавление новых операций;
// 2. Объединение родственных операции Посетителе;
// 3. Посетитель может запоминать в себе какое-то состояние по мере обхода контейнера.

// Минусы:
// 1. Затруднено добавление новых типов, поскольку нужно обновлять иерархию Посетителя и его сыновей.

type Man struct {
}

type Visitor interface {
	OwnHouse(p *Flat) string
	ParentsHouse(p *House) string
}

type Place interface {
	Accept(v Visitor) string
}

func (m *Man) OwnHouse(p *Flat) string {
	return p.TalkToWife()
}

func (m *Man) ParentsHouse(p *House) string {
	return p.TalkToParents()
}

type City struct {
	places []Place
}

func (c *City) Add(p Place) {
	c.places = append(c.places, p)
}

func (c *City) Accept(v Visitor) string {
	var result string
	for _, p := range c.places {
		result += p.Accept(v)
	}
	return result
}

type Flat struct {
}

func (f *Flat) Accept(v Visitor) string {
	return v.OwnHouse(f)
}

func (f *Flat) TalkToWife() string {
	return "Talking with wife..."
}

type House struct {
}

func (h *House) Accept(v Visitor) string {
	return v.ParentsHouse(h)
}

func (h *House) TalkToParents() string {
	return "Talking with parents..."
}

func main() {
	city := new(City)
	city.Add(&Flat{})
	city.Add(&House{})
	result := city.Accept(&Man{})
	fmt.Println(result)
}
