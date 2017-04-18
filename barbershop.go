package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

//keeping a barber working when there are customers, resting when there are none, and doing so in an orderly manner

//The barber has one barber chair and a waiting room with a number of chairs in it.
//When the barber finishes cutting a customer's hair, he dismisses the customer and then
//   goes to the waiting room to see if there are other customers waiting.
//
// If there are, he brings one of them back to the chair and cuts his hair.
// If there are no other customers waiting, he returns to his chair and sleeps in it.

type Customer struct {
	id int
}

type Barber struct {
	id int
}

func (b Barber) startedItsWork(customer Customer, done chan<- Barber) {
	fmt.Println("barber ", b, " started its work with a client from", customer)
	go func() {
		time.Sleep(time.Second * time.Duration(rand.Intn(MAX_SECONDS_NEEDED_FOR_HAIR_CUT)))
		fmt.Println("barber ", b, " is done with the customer ", customer, " and ready to take another")
		done <- b
	}()

}

const (
	BARBERS_QUANTITY                 = 3
	MAX_SECONDS_NEEDED_FOR_HAIR_CUT  = 12
	NUMBER_OF_CHAIRS_IN_WAITING_ROOM = 3
)

func initIDSeq() func() int {
	i := 0

	return func() int {
		i += 1
		return i
	}
}

func customerProducer(customers chan<- Customer) {
	nextId := initIDSeq()

	for {
		var customer = Customer{nextId()}

		customers <- customer

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1500)))
	}
}

func customerListener(availableBarbers chan Barber, customers <-chan Customer, customersWaitingInWaitingRoom chan Customer) {
	tempChanForBarbers := make(chan Barber, BARBERS_QUANTITY)

	for {
		select {
		case barber := <-tempChanForBarbers:
			select {
			case customer := <-customersWaitingInWaitingRoom:
				fmt.Println("customer ", customer, " from waiting room is taken by barber", barber)
				barber.startedItsWork(customer, tempChanForBarbers)
			default:
				fmt.Println("waiting line is empty, putting a barber", barber, " back to the list of availableBarbers")
				availableBarbers <- barber
			}
		case customer := <-customers:
			fmt.Println("visited the barbershop", customer)

			select {
			case barber := <-availableBarbers:
				barber.startedItsWork(customer, tempChanForBarbers)
			default:
				select {
				case customersWaitingInWaitingRoom <- customer:
					fmt.Println("customer has occupied a seat at waiting room", customer)
				default:
					fmt.Println("no customer sent to waiting line. Say sorry and show him a door?", customer)
				}
			}

		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	customers := make(chan Customer)

	customersWaitingInWaitingRoom := make(chan Customer, NUMBER_OF_CHAIRS_IN_WAITING_ROOM)

	availableBarbers := make(chan Barber, BARBERS_QUANTITY)

	for i := 1; i <= BARBERS_QUANTITY; i++ {
		availableBarbers <- Barber{i}
	}

	go customerListener(availableBarbers, customers, customersWaitingInWaitingRoom)
	go customerProducer(customers)

	<-time.After(time.Duration(math.MaxInt64))
	close(customers)
}
