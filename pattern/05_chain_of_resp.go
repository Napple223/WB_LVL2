package pattern

import "fmt"

/*
Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

Цепочка вызовов (обязанностей) — это поведенческий паттерн проектирования,
который позволяет передавать запросы последовательно по цепочке обработчиков.
Каждый последующий обработчик решает, может ли он обработать запрос
сам и стоит ли передавать запрос дальше по цепи.

Применяемость:
1. Когда программа должна обрабатывать разнообразные запросы несколькими
способами, но заранее неизвестно, какие конкретно запросы будут приходить
и какие обработчики для них понадобятся.
2. Когда важно, чтобы обработчики выполнялись один за другим в строгом порядке.
3. Когда набор объектов, способных обработать запрос, должен задаваться динамически.

Преимущества:
1. Уменьшает зависимость между клиентом и обработчиками.
2. Реализует принцип единственной обязанности.
3. Реализует принцип открытости/закрытости.

Недостатки:
1. Запрос может остаться никем не обработанным.
*/

//Мой пример:
//Есть структура Car. Клиент привозит автомобиль в сервис, т.к. у него
//имеется неисправность. Клиент указывает какого рода неисправность (подвеска,
//кузов, двигатель), а мы передаем его по цепочке вызовов до тех пор,
//пока не найдем тот отдел в сервисе, который сможет починить неисправность.

//Структура автомобиля.
type Car struct {
	problem string
}

//Функция - конструктор, оборачивающая запрос в объект.
func NewCar(someProblem string) *Car {
	return &Car{someProblem}
}

//Интерфейс цепочки. С помощью него мы можем добавлять
//в цепочку вызово все нужные обработчики, если они
//удовлетворяют интерфейсу.
type CheckProblem interface {
	execute(*Car)
	setNext(CheckProblem)
}

//Звено цепочки вызовов.
type Reseption struct {
	next CheckProblem
}

//Метод для проверки способности данного звена
//обработать запрос. Если способен - обрываем
//цепочку вызовов, если нет - передаем запрос
//следующему обработчику.
func (r *Reseption) execute(c *Car) {
	if c.problem == "i just want to wash my car." {
		fmt.Println("We can wash your car.")
		//car wash
		return
	}
	if r.next == nil {
		return
	}
	r.next.execute(c)

}

//Метод для установки следующего обработчика.
func (r *Reseption) setNext(next CheckProblem) {
	r.next = next
}

//Звено цепочки вызовов.
type SuspensionDep struct {
	next CheckProblem
}

//Метод для проверки способности данного звена
//обработать запрос. Если способен - обрываем
//цепочку вызовов, если нет - передаем запрос
//следующему обработчику.
func (s *SuspensionDep) execute(c *Car) {
	if c.problem == "I have some problem with suspension." {
		fmt.Println("I can fix your problem with suspension.")
		//check suspension
		return
	}
	if s.next == nil {
		return
	}
	s.next.execute(c)

}

//Метод для установки следующего обработчика.
func (s *SuspensionDep) setNext(next CheckProblem) {
	s.next = next
}

type EngineDep struct {
	next CheckProblem
}

func (e *EngineDep) execute(c *Car) {
	if c.problem == "I have some problem with engine." {
		fmt.Println("I can fix your problem with engine.")
		//check engine
		return
	}
	if e.next == nil {
		return
	}
	e.next.execute(c)

}

func (e *EngineDep) setNext(next CheckProblem) {
	e.next = next
}

//И т.д. и т.п. наращивать цепочку вызовов можно очень долго...

func testChainOfResp() {
	//Удобнее всего в этой ситуации строить цепочку вызовов с конца.
	//Инициализируем последнее звено и передаем его предыдущему обработчику.
	engineDep := &EngineDep{}
	suspensionDep := &SuspensionDep{}
	suspensionDep.setNext(engineDep)
	reseption := &Reseption{}
	reseption.setNext(suspensionDep)
	//Создаем запрос и оборачиваем его в объект.
	suv := NewCar("I have some problem with engine.")
	//Передаем объект запроса первому звену обработчика.
	reseption.execute(suv) //I can fix your problem with engine.
	anotherCar := NewCar("I don't know what happened.")
	reseption.execute(anotherCar) //Запрос не будет обработан, так как нет нужного обработчика.
}
