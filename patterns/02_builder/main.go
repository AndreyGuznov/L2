package main

import "fmt"

/*
	Строитель (англ. Builder) — порождающий шаблон проектирования предоставляет способ создания составного объекта.
Отделяет конструирование сложного объекта от его представления так, что в результате одного и того
же процесса конструирования могут получаться разные представления.

Паттерн Строитель предлагает вынести конструирование объекта за пределы его собственного класса, поручив это дело
отдельным объектам, называемым строителями.

	плюсы:
Позволяет создавать продукты пошагово.
Позволяет использовать один и тот же код для создания различных продуктов.
Изолирует сложный код сборки продукта от его основной бизнес-логики.

	минусы:
    Усложняет код программы из-за введения дополнительных классов.
	Клиент будет привязан к конкретным классам строителей, так как в интерфейсе директора может не быть метода получения результата.
*/

type Builder interface {
	BuildBrand(v string) Builder
	BuildCore(v int) Builder
	BuildMemory(v int) Builder
	BuildGraphic(v int) Builder
	BuildComputer() Computer
}

type Computer struct {
	Brand   string
	Core    int
	Memory  int
	Graphic int
}

type computerBuilder struct {
	brand   string
	core    int
	memory  int
	graphic int
}

func (cB computerBuilder) BuildBrand(v string) Builder {
	cB.brand = v
	return cB
}

func (cB computerBuilder) BuildCore(v int) Builder {
	cB.core = v
	return cB
}

func (cB computerBuilder) BuildMemory(v int) Builder {
	cB.memory = v
	return cB
}

func (cB computerBuilder) BuildGraphic(v int) Builder {
	cB.graphic = v
	return cB
}

func (cB computerBuilder) BuildComputer() Computer {
	return Computer{
		Brand:   cB.brand,
		Core:    cB.core,
		Memory:  cB.memory,
		Graphic: cB.graphic,
	}
}
func NewComputerBuilder() Builder {
	return computerBuilder{}
}

func main() {
	compBuild := NewComputerBuilder()
	comp1 := compBuild.BuildBrand("Sony").BuildCore(4).BuildMemory(8).BuildGraphic(2)
	comp2 := compBuild.BuildBrand("HP").BuildCore(6).BuildMemory(32).BuildGraphic(4)
	fmt.Println(comp1, comp2)
}
