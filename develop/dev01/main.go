package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===
Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.
Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {

	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Println(err)
	}
	err = response.Validate()
	if err != nil {
		log.Fatal(err) // response data is suitable for synchronization purposes
	}
	opt1 := response.RootDelay
	opt2 := response.RootDistance
	time := time.Now().Add(response.ClockOffset)
	fmt.Println(time)
	fmt.Println(opt1, opt2)
}
