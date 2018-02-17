package main

import (
	"time"
	"fmt"
)

type Address struct {
	city     string
	district string
}

type Person struct {
	Name    string
	Age     uint8
	Address Address
}

type PersonHandler interface {
	Batch(origins <-chan Person) <-chan Person
	Handle(origins *Person)
}

var personTotal = 200
var persons = make([]Person, personTotal)
var personCount int

func init()  {
	for i :=0; i < 200; i++  {
		name := fmt.Sprintf("%s%d", "P", i)
		p := Person{name, 32, Address{"Beijing", "Haidian"}}
		persons[i] = p
	}
}

type PersonHandlerImpl struct {}

// 后面的G3和G4,是为了让此批处理流程完全的异步化，异步地获取和存储人员信息
//    G1             G2              G3                 G4
//  初始化和协调   从通道origins     把人员信息           存储人员信息
//               接收人员信息，变更  发送给通道origins
func main() {
	handler := getPersonHandler()
	origins := make(chan Person, 100)
	destinations := handler.Batch(origins)
	fetchPerson(origins)
	sign := savePerson(destinations)
	<-sign
}
func savePerson(dest <-chan Person) <-chan byte {
	sign := make(chan byte, 1)
	go func() {
		for {
			p, ok := <-dest
			if !ok {
				fmt.Println("All the information has been saved.")
				sign <- 0
				break
			}
			saveOneByOne(p)
		}
	}()
	return sign
}
func saveOneByOne(p Person) bool {
	fmt.Printf("%s\n", p.Name)
	return true
}

//3个细节
// 1.是goroutine池，循环放东西到缓冲区
// 2.在收发两端都有并发需求下，使用非缓冲通道作为元素值传输介质是不合适的
// 3.如果已经没有人员信息可取了（!ok），还要检查被启用的goroutine是否运行完毕（==）
func fetchPerson(origins chan<- Person) {
	originCapacity := cap(origins)
	buffered := originCapacity > 0
	goTicketTotal := originCapacity / 2
	goTicket := initGoTicket(goTicketTotal)
	go func() {
		for {
			p, ok := fetchOneByOne()
			if !ok {
				// this for is good
				for {
					if !buffered || len(goTicket) == goTicketTotal {
						break
					}
					time.Sleep(time.Nanosecond)
				}
				fmt.Println("All the infomation has been fetched")
				close(origins)
				break
			}
			if buffered {
				// this is very good
				<-goTicket
				go func() {
					origins <- p
					goTicket <- 1
				}()
			} else {
				origins <- p
			}
		}
	}()
}

func fetchOneByOne() (Person, bool) {
	if personCount < personTotal{
		p := persons[personCount]
		personCount++
		return p, true
	}

	return Person{}, false
}

// 1.goroutine池
// 不直接赋值给goTicket变量，首先这是个缓冲channel,循环放进缓冲,里面元素值的个数代表了还没有被获得
// 和已被归还的票数总和，那么在初始化的时候，其中所有的票都没有被获得。
// 如果不这么做，从该池中拿票都会被阻塞（缓冲通道的2个特性：1.缓冲区满了，别人发不了东西到这个通道 2.缓冲区空了，接收不到东西）
func initGoTicket(total int) chan byte {
	var goTicket chan byte
	if total == 0 {
		return goTicket
	}
	goTicket = make(chan byte, total)
	for i := 0; i < total; i++  {
		goTicket <- 1
	}

	return goTicket
}

func getPersonHandler() PersonHandler {
	return PersonHandlerImpl{}
}

func (handler PersonHandlerImpl) Batch(origins <-chan Person) <-chan Person {
	destinations := make(chan Person, 100)
	go func() {
		for item := range origins {
			handler.Handle(&item)
			destinations <- item
		}
		fmt.Println("All the information has been handled.")
		close(destinations)
	}()

	return destinations
}

func (handler PersonHandlerImpl) Handle(orig *Person) {
	if orig.Address.district == "Haidian" {
		orig.Address.district = "Shijingshan"
	}
}


