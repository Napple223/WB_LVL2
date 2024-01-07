package pattern

import (
	"errors"
	"fmt"
)

/*
Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/State_pattern

Состояние — это поведенческий паттерн проектирования,
который позволяет объектам менять поведение в
зависимости от своего состояния.
Извне создаётся впечатление, что изменился класс объекта.

Применяемость:
1. Когда есть объект, поведение которого кардинально
меняется в зависимости от внутреннего состояния,
причём типов состояний много, и их код часто меняется.
2. Когда код класса содержит множество больших,
похожих друг на друга, условных операторов, которые
выбирают поведения в зависимости от текущих значений полей класса.
3. Когда вы сознательно используете табличную машину состояний,
построенную на условных операторах, но вынуждены мириться с
дублированием кода для похожих состояний и переходов.

Преимущества:
1. Избавляет от множества больших условных операторов машины состояний.
2. Концентрирует в одном месте код, связанный с определённым состоянием.

Недостатки:
1. Может неоправданно усложнить код, если состояний мало и они редко меняются.
*/

//Мой пример:
//Допустим у нас есть вендинговый автомат, выдающий 1 товар,
//автомат может пребывать в 4 состояних:
//товар в наличие, товар не в наличие, товар в процессе выдачи,
//деньги оплачен.

// Интерфейс состояний.
type State interface {
	addItem(int) error
	requestItem() error
	insetMoney(int) error
	dispenseItem() error
}

// Структура вендингового автомата.
// В отличие от паттерна статегия,
// состояния должны знать друг о друге
// и самостоятельно устанавливать
// необходимое состояние в той
// или иной ситуации.
type VendingMachine struct {
	hasItem       State
	itemRequested State
	hasMoney      State
	noItem        State
	curState      State
	itemCount     int
	itemPrise     int
}

// При вызове методов, будет произведен вызов соответствующего метода
// текущего состояния.
func (v *VendingMachine) addItem(count int) error {
	return v.curState.addItem(count)
}

func (v *VendingMachine) requestItem() error {
	return v.curState.requestItem()
}

func (v *VendingMachine) insetMoney(count int) error {
	return v.curState.insetMoney(count)
}

func (v *VendingMachine) dispenseItem() error {
	return v.curState.dispenseItem()
}

// Метод для установки состояния.
func (v *VendingMachine) setState(s State) {
	v.curState = s
}

// Метод для добавления товаров в автомат.
func (v *VendingMachine) incrementItemsCount(count int) {
	v.itemCount += count
	fmt.Printf("%d items added\n", count)
}

// Функция - конструктор.
func NewVM(count, price int) *VendingMachine {
	v := &VendingMachine{
		itemCount: count,
		itemPrise: price,
	}
	hasItemState := &HasItemState{v: v}
	itemRequstedState := &ItemRequestedState{v: v}
	hasMoneyState := &HasMoneyState{v: v}
	noItemState := &NoItemState{v: v}
	v.setState(hasItemState)
	v.hasItem = hasItemState
	v.itemRequested = itemRequstedState
	v.hasMoney = hasMoneyState
	v.noItem = noItemState
	return v
}

//Блок с описанием состояний и их методов.

type NoItemState struct {
	v *VendingMachine
}

func (noState *NoItemState) addItem(count int) error {
	noState.v.incrementItemsCount(count)
	noState.v.setState(noState.v.hasItem)
	return nil
}

func (noState *NoItemState) requestItem() error {
	return errors.New("item out of stock")
}

func (noState *NoItemState) insetMoney(count int) error {
	return errors.New("item out of stock")
}

func (noState *NoItemState) dispenseItem() error {
	return errors.New("item out of stock")
}

type HasItemState struct {
	v *VendingMachine
}

func (hasState *HasItemState) addItem(count int) error {
	hasState.v.incrementItemsCount(count)
	return nil
}

func (hasState *HasItemState) requestItem() error {
	if hasState.v.itemCount == 0 {
		hasState.v.setState(hasState.v.noItem)
		return errors.New("item out of stock")
	}
	fmt.Printf("Item requestd\n")
	hasState.v.setState(hasState.v.itemRequested)
	return nil
}

func (hasState *HasItemState) insetMoney(count int) error {
	return errors.New("first select an item")
}

func (hasState *HasItemState) dispenseItem() error {
	return errors.New("first select an item")
}

type ItemRequestedState struct {
	v *VendingMachine
}

func (req *ItemRequestedState) addItem(count int) error {
	return errors.New("item dispense in progress, try later")
}

func (req *ItemRequestedState) requestItem() error {
	return errors.New("item already requested")
}

func (req *ItemRequestedState) insetMoney(count int) error {
	if count < req.v.itemPrise {
		return errors.New("not enough money")
	}
	req.v.setState(req.v.hasMoney)
	return nil
}

func (req *ItemRequestedState) dispenseItem() error {
	return errors.New("insert money first")
}

type HasMoneyState struct {
	v *VendingMachine
}

func (m *HasMoneyState) addItem(count int) error {
	return errors.New("item dispense in progress, try later")
}

func (m *HasMoneyState) requestItem() error {
	return errors.New("item already requested")
}

func (m *HasMoneyState) insetMoney(count int) error {
	return errors.New("money already inserted")
}

func (m *HasMoneyState) dispenseItem() error {
	m.v.itemCount--
	fmt.Println("Dispensing item.")
	if m.v.itemCount == 0 {
		m.v.setState(m.v.noItem)
		return nil
	}
	m.v.setState(m.v.hasItem)
	return nil
}

//Очень большой блок...

func testState() {
	v := NewVM(1, 5) //has item state
	v.requestItem()  //item requested state
	v.insetMoney(5)  //has money state
	v.dispenseItem() //no item state
	v.addItem(1)     //has item state
}
