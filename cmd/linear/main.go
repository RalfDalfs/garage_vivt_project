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

// Вся обработка команд и модель стоянки находятся в main.
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
			fmt.Println("Стоянка от дальнего конца к выезду:")
			for _, car := range cars {
				fmt.Printf("%s: перемещений %d\n", car.Number, car.Moves)
			}
			continue
		}
		if len(fields) != 2 || (fields[0] != "+" && fields[0] != "-") {
			fmt.Println("Ошибка: используйте + НОМЕР, - НОМЕР, list или exit")
			continue
		}
		number := fields[1]
		index := -1
		for i, car := range cars {
			if car.Number == number {
				index = i
				break
			}
		}
		if fields[0] == "+" {
			if index != -1 {
				fmt.Println("Ошибка: автомобиль уже на стоянке")
				continue
			}
			cars = append(cars, Car{Number: number})
			fmt.Printf("Прибыл %s\n", number)
			continue
		}
		if index == -1 {
			fmt.Println("Ошибка: автомобиль не найден")
			continue
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
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}
}
