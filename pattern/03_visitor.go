package pattern

import "fmt"

/*
Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Visitor_pattern

Посетитель — это поведенческий паттерн проектирования,
который позволяет добавлять в программу новые операции,
не изменяя классы объектов (изменяя, но единожды для всех
новых фич), над которыми эти операции могут выполняться.

Применяемость:
1. Когда нужно выполнить какую-то операцию над всеми
элементами сложной структуры объектов, например, деревом.
2. Когда над объектами сложной структуры объектов надо
выполнять некоторые не связанные между собой операции,
но вы не хотите «засорять» классы такими операциями.
3. Когда новое поведение имеет смысл только для некоторых
классов из существующей иерархии.

Преимущества:
1. Упрощает добавление операций, работающих со сложными структурами объектов.
2. Объединяет родственные операции в одном классе.
3. Посетитель может накапливать состояние при обходе структуры элементов.

Недостатки:
1. Паттерн не оправдан, если иерархия элементов часто меняется.
2. Может привести к нарушению инкапсуляции элементов.
*/

//Мой пример:
//Допустим имеется некий метеоцентр, который принимает информацию о
//погоде от некоторой группы датчиков, расположенной удаленно друг от друга.
//Наши коллеги с соседнего метеоцентра, в рамках проводимого ими эксперимента,
//попросили предоставить доп.информацию с датчиков (замер ультрафиолетового
//излучения и количества осадков). Нам эти данные в рамках наших задач не нужны
//и мы не хотим добавлять их к системе управления метеоцентром.

//Интерфейс нашего метеоцентра.
type meteoCenter interface {
	getTemperature() string
	//Метод, который нам необходимо добавить для реализации паттерна.
	accept(visitor)
}

//Интерфейс посетителя с методами для каждого типа объекта.
//Позволяет добавить новый функционал, не затрагивая существующую
//иерархию.
type visitor interface {
	visitForSensorType1(*sensorType1)
	visinForSensorType2(*sensorType2)
}

//Тип существующего элемента.
type sensorType1 struct {
	id int
}

//Метод, имплементирующий интерфейс метеоцентра.
func (s *sensorType1) getTemperature() string {
	return "some temperature from sensor type 1"
}

//Метод, позволяющий добавить непредусмотренный ранее функционал.
//Существующему типу (классу).
func (s *sensorType1) accept(v visitor) {
	v.visitForSensorType1(s)
}

//Тип существующего элемента.
type sensorType2 struct {
	id       int
	someData string
}

//Метод, имплементирующий интерфейс метеоцентра.
func (s *sensorType2) getTemperature() string {
	return "some temperature from sensor type 2"
}

//Метод, позволяющий добавить непредусмотренный ранее функционал.
//Существующему типу (классу).
func (s *sensorType2) accept(v visitor) {
	v.visinForSensorType2(s)
}

//Структрура посетителя.
type uvCureLevel struct{}

//Методы, добавляющий новый функционал к существующему типу (классу).
func (uv *uvCureLevel) visitForSensorType1(s *sensorType1) {
	fmt.Println("Some UV cure level info from sensor type 1.")
}

func (uv *uvCureLevel) visinForSensorType2(s *sensorType2) {
	fmt.Println("Some UV cure level info from sensor type 2.")
}

//Структура еще одного посетителя.
type precipitation struct{}

//Методы, добавляющий новый функционал к существующему типу (классу).
func (p *precipitation) visitForSensorType1(s *sensorType1) {
	fmt.Println("Some info about the amount of precipitation from sensor type 1")
}

func (p *precipitation) visinForSensorType2(s *sensorType2) {
	fmt.Println("Some info about the amount of precipitation from sensor type 2")
}

func testVisitor() {
	//Инициализируем существующие типы.
	sensorType1 := &sensorType1{id: 236478}
	sensorType2 := &sensorType2{
		id:       1874121,
		someData: "made in China",
	}
	//Их обычное поведение.
	sensorType1.getTemperature()
	sensorType2.getTemperature()

	//Инициализируем структуру доп.функционала.
	uvCure := &uvCureLevel{}
	sensorType1.accept(uvCure)
	sensorType2.accept(uvCure)

	//Инициализируем еще один доп.функционал,
	//используя теже методы.
	precipitation := &precipitation{}
	sensorType1.accept(precipitation)
	sensorType2.accept(precipitation)
}
