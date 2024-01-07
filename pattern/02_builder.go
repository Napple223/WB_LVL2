package pattern

/*
Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Builder_pattern

Строитель — это порождающий паттерн проектирования,
который позволяет создавать сложные объекты пошагово.
Строитель даёт возможность использовать один и тот же
код строительства для получения разных представлений объектов.

Применяемость:
1. Когда вы хотите избавиться от «телескопического конструктора».
2. Когда ваш код должен создавать разные представления какого-то объекта.
3. Когда вам нужно собирать сложные составные объекты.

Преимущества:
1. Позволяет создавать продукты пошагово.
2. Позволяет использовать один и тот же код для создания различных продуктов.
3. Изолирует сложный код сборки продукта от его основной бизнес-логики.

Недостатки:
1. Усложняет код программы из-за введения дополнительных классов.
2. Клиент будет привязан к конкретным классам строителей.
*/

//Мой пример:
//Есть структура пиццы, состав которой может изменяться (толстое/тонкое тесто,
//соус, начинка и т.д.).

//Структура пиццы.
type pizza struct {
	dough   string
	sause   string
	topping string
	cheese  string
}

//Интерфейс строителя пиццы.
type pizzaBuilder interface {
	setDoughType()    //Устанавливаем тип теста.
	setSause()        //Устанавливаем тип соуса.
	setToppingType()  //Устанавливаем тип топпинга.
	setCheese()       //Устанавливаем тип сыра.
	makePizza() pizza //Метод для создания конкретной реализации пиццы с нужными полями.
}

//Функция для установки конкретного типа строителя.
func getPizzaBuilder(pizzaType string) pizzaBuilder {
	if pizzaType == "pepperoni" {
		return newPepperoniPizza()
	}
	if pizzaType == "hawaiian" {
		return newHawaiianPizza()
	}
	return nil
}

//Структура конкретного объекта, реализуемого
//билдером.
type peperroniPizza struct {
	dough   string
	sause   string
	topping string
	cheese  string
}

//Функция-конструктор конкретной реализации.
func newPepperoniPizza() *peperroniPizza {
	return &peperroniPizza{}
}

//Методы, имплементирующие интерфейс билдера.
func (p *peperroniPizza) setDoughType() {
	p.dough = "thin"
}

func (p *peperroniPizza) setSause() {
	p.sause = "tomato"
}

func (p *peperroniPizza) setToppingType() {
	p.topping = "pepperoni"
}

func (p *peperroniPizza) setCheese() {
	p.cheese = "mozzarella"
}

func (p *peperroniPizza) makePizza() pizza {
	return pizza{
		dough:   p.dough,
		sause:   p.sause,
		topping: p.topping,
		cheese:  p.cheese,
	}
}

//Структрура другой конкретной реализации пиццы.
type hawaiianPizza struct {
	dough   string
	sause   string
	topping string
	cheese  string
}

//Функция-конструктор.
func newHawaiianPizza() *hawaiianPizza {
	return &hawaiianPizza{}
}

//Методы, имплементирующие интерфейс.
func (h *hawaiianPizza) setDoughType() {
	h.dough = "thick"
}

func (h *hawaiianPizza) setSause() {
	h.sause = "tomato"
}

func (h *hawaiianPizza) setToppingType() {
	h.topping = "pineapple"
}

func (h *hawaiianPizza) setCheese() {
	h.cheese = "dutch cheese"
}

func (h *hawaiianPizza) makePizza() pizza {
	return pizza{
		dough:   h.dough,
		sause:   h.sause,
		topping: h.topping,
		cheese:  h.cheese,
	}
}

//Структура "директора". Только директор знает как правильно
//собрать пиццу. И именно с ним взаимодействует клиент.
type chef struct {
	builder pizzaBuilder
}

//Функция-конструктор.
func newChef(b pizzaBuilder) *chef {
	return &chef{
		builder: b,
	}
}

//Метод, для установки конкретной реализации билдера.
//Подходят только те реализации, которые удовлетворяют
//интерфейсу билдера.
func (c *chef) setPizzaBuilder(b pizzaBuilder) {
	c.builder = b
}

//Функция, определяющая порядок сборки объекта пицца.
//Позволяет не допустить "недоделанных" реализаций структуры
//объекта.
func (c *chef) makeNewPizza() pizza {
	c.builder.setDoughType()
	c.builder.setSause()
	c.builder.setToppingType()
	c.builder.setCheese()
	return c.builder.makePizza()
}

//Пример реализации билдера.
func builderTest() {
	//Выбираем конкретную реализацию объекта пицца.
	peperroniBuilder := getPizzaBuilder("pepperoni")
	//Передаем "директору" инструкцию по сборке объекта.
	chef := newChef(peperroniBuilder)
	//Директор запускает цикл сборки.
	_ = chef.makeNewPizza()

	//Аналогично для другой реализации объекта.
	hawaiianBuilder := getPizzaBuilder("hawaiian")
	//Повторно создавать директора не нужно.
	//Можно просто установить новые инструкции по сборке объекта.
	chef.setPizzaBuilder(hawaiianBuilder)
	_ = chef.makeNewPizza()
}
