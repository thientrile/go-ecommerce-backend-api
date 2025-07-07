package main

import (
	"fmt"
	"time"
)

// type Course struct {
// 	Title string
// 	Price int
// }

// func main() {
// 	// add channel to
// 	ch := make(chan Course)
// 	// 2. create a goroutine that sends data to the channel
// 	go func() {
// 		course := Course{Title: "Goroutine and Channel", Price: 100}
// 		ch <- course // send data to the channel
// 	}()
// 	c := <-ch
// 	fmt.Printf("Received Course: %s, Price: %d\n", c.Title, c.Price)
// }

// pub/sub using channel

type Message struct {
	OrderId int
	Title   string
	Price   int
}

func buyTicket(ch chan<- Message, orders []Message) {
	for _, order := range orders {
		fmt.Printf("Processing Order: %d, Title: %s, Price: %d\n", order.OrderId, order.Title, order.Price)
		time.Sleep(1 * time.Second) // simulate processing time
		ch <- order                 // send order to the channel
	}
	fmt.Println("Buying ticket...")
	close(ch) // close the channel after sending all orders
}
func canncelTicket(ch chan<- int, orders []int) {
	for _, orderId := range orders {
		time.Sleep(5 * time.Second) // simulate processing time
		fmt.Printf("Cancelling Order: %d\n", orderId)
		ch <- orderId // send cancellation as Message
	}
	fmt.Println("Cancelling ticket...")
	close(ch) // close the channel after sending all cancellations
}

func handlderOrder(orderChannel <-chan Message, cancelChannel <-chan int) {
	for {
		select {
		case order, orderOke := <-orderChannel:
			if orderOke {

				fmt.Printf("Handle Received Order: %d, Title: %s, Price: %d\n", order.OrderId, order.Title, order.Price)
			} else {
				fmt.Println("No more orders to process.")
				orderChannel = nil // set to nil to avoid further processing

			}
		case cancel, cancelOke := <-cancelChannel:
			if cancelOke {
				fmt.Printf("Handle Received Cancellation for Order: %d\n", cancel)
			} else {
				fmt.Println("No more cancellations to process.")
				orderChannel = nil // set to nil to avoid further processing

			}
		}
		// check if both channels are closed
		if orderChannel == nil && cancelChannel == nil {
			break // exit the loop if both channels are closed
		}
	}
}
func main() {
	buyChannel := make(chan Message)
	cancelChannel := make(chan int)
	// simulate order processing
	buyOrders := []Message{
		{OrderId: 1, Title: "Goroutine and Channel", Price: 100},
		{OrderId: 2, Title: "Go Programming", Price: 200},
		{OrderId: 3, Title: "Concurrency in Go", Price: 300},
	}
	canncelOrders := []int{1, 2, 3}
	go buyTicket(buyChannel, buyOrders)
	go canncelTicket(cancelChannel, canncelOrders)
	// handle orders and cancellations
	// this will block until the channels are closed
	go handlderOrder(buyChannel, cancelChannel)
	time.Sleep(15 * time.Second) // wait for goroutines to finish
	fmt.Println("All orders processed.")

}

// func subscribe(ch <-chan Message, user string) {
// 	for msg := range ch {
// 		fmt.Printf("User: %s, Received Order: %d, Title: %s, Price: %d\n", user, msg.OrderId, msg.Title, msg.Price)
// 		time.Sleep(1 * time.Second) // simulate processing time
// 	}
// }

// func publish(ch chan<- Message, orders []Message) {
// 	for _, order := range orders {
// 		fmt.Printf("Publishing Order: %d, Title: %s, Price: %d\n", order.OrderId, order.Title, order.Price)
// 		ch <- order // send order to the channel
// 	}
// 	close(ch) // close the channel after sending all orders
// }

// func main() {
// 	// create a channel order
// 	orderChannel := make(chan Message)
// 	orders := []Message{
// 		{OrderId: 1, Title: "Goroutine and Channel", Price: 100},
// 		{OrderId: 2, Title: "Go Programming", Price: 200},
// 		{OrderId: 3, Title: "Concurrency in Go", Price: 300},
// 	}
// 	go publish(orderChannel, orders)
// 	go subscribe(orderChannel, "Gfj user")
// 	time.Sleep(5 * time.Second) // wait for goroutines to finish
// 	fmt.Println("All orders processed.")
// }
