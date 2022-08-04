package main

import "fmt"

/*
	Паттерн позволяет реализовать передачу запросов по последовательности
	обработчиков, где каждый обработчик решает может ли он обработать запрос
	и нужно ли передавать его дальше по цепи.
	Набор обработчиков может задаваться динамически, звено запроса всегда
	можно добавить, заменить или удалить.
	Цепь обработки запросов при этом может состоять лишь
	из одного обработчика и запрос может быть не обработан совсем.

	Плюсы:
	- Разнесение клиента и обработчиков, уменьшение их зависимости
	- Реализация принципа единственной ответственности

	Минусы:
	- Создание дополнительных объектов, усложнение кода
	- Запрос может быть не обработан
	Примером реализации паттерна может служить последовательность
	движения больного по больнице: он попадает в приемное отделение,
	затем к врачу, затем к кассиру на оплату приема.
*/

type Service interface {
	Execute(*Data)
	SetNext(Service)
}

type Data struct {
	GetSource    bool
	UpdateSource bool
}

//service Get data
type Device struct {
	Name string
	Next Service
}

func (d *Device) Execute(data *Data) {
	if data.GetSource {
		fmt.Printf("Data from device %s already get\n", d.Name)
		d.Next.Execute(data)
		return
	}
	fmt.Printf("Sucsess of geting data from %s\n", d.Name)
	data.GetSource = true
	d.Next.Execute(data)
}

func (d *Device) SetNext(serv Service) {
	d.Next = serv
}

//service Update data

type UpdateDataService struct {
	Name string
	Next Service
}

func (u *UpdateDataService) Execute(data *Data) {
	if data.UpdateSource {
		fmt.Printf("Data from device %s already update\n", u.Name)
		u.Next.Execute(data)
		return
	}
	fmt.Printf("Sucsess of updating data from %s\n", u.Name)
	data.GetSource = true
	u.Next.Execute(data)
}

func (u *UpdateDataService) SetNext(serv Service) {
	u.Next = serv
}

// service Save data

type SaveDataService struct {
	Next Service
}

func (s *SaveDataService) Execute(data *Data) {
	if data.UpdateSource {
		fmt.Println("Data need to be updated!")
		return
	}
	fmt.Println("Data Saved")
}

func (s *SaveDataService) SetNext(serv Service) {
	s.Next = serv
}

func main() {
	dev := &Device{Name: "SomeDevice"}
	updSvc := &UpdateDataService{Name: "Update"}
	saveDataSvc := &SaveDataService{}
	dev.SetNext(updSvc)
	updSvc.SetNext(saveDataSvc)
	data := &Data{}
	dev.Execute(data)
}
