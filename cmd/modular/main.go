package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Car struct {
	Number string
	Moves  int
}

func find(cars []Car, number string) int {
	for i, car := range cars {
		if car.Number == number {
			return i
		}
	}
	return -1
}

func arrive(cars []Car, number string) []Car {
	if find(cars, number) != -1 {
		fmt.Println("Ошибка: автомобиль уже на стоянке")
		return cars
	}
	fmt.Printf("Прибыл %s\n", number)
	return append(cars, Car{Number: number})
}

func depart(cars []Car, number string) []Car {
	index := find(cars, number)
	if index == -1 {
		fmt.Println("Ошибка: автомобиль не найден")
		return cars
	}
	var temporary []Car
	for len(cars)-1 > index {
		car := cars[len(cars)-1]
		cars = cars[:len(cars)-1]
		car.Moves++
		temporary = append(temporary, car)
		fmt.Printf("Временно выехал %s\n", car.Number)
	}
	car := cars[len(cars)-1]
	cars = cars[:len(cars)-1]
	fmt.Printf("Уехал %s: перемещений %d\n", car.Number, car.Moves)
	for len(temporary) > 0 {
		car := temporary[len(temporary)-1]
		temporary = temporary[:len(temporary)-1]
		cars = append(cars, car)
		fmt.Printf("Вернулся %s\n", car.Number)
	}
	return cars
}

func show(cars []Car) {
	fmt.Println("Стоянка от дальнего конца к выезду:")
	for _, car := range cars {
		fmt.Printf("%s: перемещений %d\n", car.Number, car.Moves)
	}
}

func main() {
	var cars []Car
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Команды: + НОМЕР, - НОМЕР, list, exit")
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) == 1 && fields[0] == "exit" {
			return
		}
		if len(fields) == 1 && fields[0] == "list" {
			show(cars)
			continue
		}
		if len(fields) != 2 || (fields[0] != "+" && fields[0] != "-") {
			fmt.Println("Ошибка: используйте + НОМЕР, - НОМЕР, list или exit")
			continue
		}
		if fields[0] == "+" {
			cars = arrive(cars, fields[1])
		} else {
			cars = depart(cars, fields[1])
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}
}
