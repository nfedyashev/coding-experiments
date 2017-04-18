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
	isBusy bool
}

const (
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

func customerProducer(customers chan Customer) {
	nextId := initIDSeq()

	for {
		var customer = Customer{nextId()}

		customers <- customer

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1500)))
	}
}

func customerListener(barber Barber, customers chan Customer, customersWaitingInWaitingRoom chan Customer) {
	for {
		select {
		case customer := <-customers:
			fmt.Println("visited the barbershop", customer)

			if barber.isBusy {
				fmt.Println("barber is busy, lets see what we can do with the customer", customer)

				select {
				case customersWaitingInWaitingRoom <- customer:
					fmt.Println("customer has occupied a seat at waiting room", customer)
				default:
					fmt.Println("no customer sent to waiting line. Say sorry and show him a door?", customer)
				}

			} else {

				select {
				// kind of an overkill right now, but customer may have a relax-time in future versions
				case customer := <-customersWaitingInWaitingRoom:
					fmt.Println("barber started its work with a client from waiting room", customer)
					barber.isBusy = true
					go func() {
						time.Sleep(time.Second * time.Duration(rand.Intn(MAX_SECONDS_NEEDED_FOR_HAIR_CUT)))
						barber.isBusy = false
						fmt.Println("barber is done with the customer who had to wait a bit (", customer, ") and ready to take another")
					}()
				default:
					fmt.Println("barber started its work with a client who just entered the shop", customer)
					barber.isBusy = true
					go func() {
						time.Sleep(time.Second * time.Duration(rand.Intn(MAX_SECONDS_NEEDED_FOR_HAIR_CUT)))
						barber.isBusy = false
						fmt.Println("barber is done with the lucky customer ", customer, "and ready to take another")
					}()
				}
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	customers := make(chan Customer)

	customersWaitingInWaitingRoom := make(chan Customer, NUMBER_OF_CHAIRS_IN_WAITING_ROOM)

	var barber = Barber{false}

	go customerListener(barber, customers, customersWaitingInWaitingRoom)
	go customerProducer(customers)

	<-time.After(time.Duration(math.MaxInt64))
	close(customers)
}
