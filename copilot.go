/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/8 2:20 PM
 * @Desc:
 */

package main

import "fmt"

func mainaaa() {
	f := NewAnimalFactory(&Person{
		Name: "b",
		Age:  0,
	})
	fmt.Println(f.GetName())
}
func Hello(who string) {
	fmt.Printf("hello %s", who)

}

type AnimalFactory struct {
	SomeAnimal Animal
}

func (a *AnimalFactory) GetName() string {
	return a.SomeAnimal.GetName()
}
func NewAnimalFactory(a Animal) *AnimalFactory {
	return &AnimalFactory{
		SomeAnimal: a}
}

type Animal interface {
	SetName(name string) bool
	GetName() string
	SetAge(age int) bool
	GetAge() int
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) SetName(name string) bool {
	p.Name = name
	return true
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) SetAge(age int) bool {
	p.Age = age
	return true
}

func (p *Person) GetAge() int {
	return p.Age
}
