package main

import (
  "fmt"
)

type human interface{
  speak()
}

func saySomething(h human){
  h.speak()
}

type person struct {
  fname string
  lname string
}

func (p person) speakName(x int) int {
  fmt.Println("My name is",  p.fname, p.lname)
  return x
}

func (p person) speak() {
  fmt.Println("My name is",  p.fname, p.lname)
}

type secreatAgent struct {
  person
  tatti bool
}

func (sa secreatAgent) speak() {
  fmt.Println("Am I tatti")
  if (sa.tatti){
    fmt.Println("Yes")
  }
  fmt.Println("My name is", sa.fname, sa.lname)
}

func main() {
  fmt.Println("Hello, playground")
  p := person {
    "tatti",
    "aadmi",
  }
  // p.speakName(2)

  s := secreatAgent{
    person{
      fname: "Yoyo",
      lname: "Honey",
    },
    true,
  }
  // s.speak()
  saySomething(p)
  saySomething(s)
}
