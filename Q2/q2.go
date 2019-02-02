package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NumRoutines = 3
	NumRequests = 1000
)

// global semaphore monitoring the number of routines
var semRout = make(chan int, NumRoutines)

// global semaphore monitoring console
var semDisp = make(chan int, 1)

// Waitgroups to ensure that main does not exit until all done
var wgRout sync.WaitGroup
var wgDisp sync.WaitGroup

type Task struct {
	a, b float32
	disp chan float32
}

/*A function that sleeps for a random time between 1 and 15 seconds, adds
the numbers a and b and sends the result on the display channel. */
func solve(t *Task) {
	time.Sleep(time.Duration(rand.Intn(15)+1) * time.Second)
	t.disp <- t.a + t.b

}

/* A function that acts as intermediary between ComputeServer and
solve.*/
func handleReq(t *Task) {
	wgRout.Add(1)
	solve(t)
	<-semRout
	wgRout.Done()

}

/* A function that uses the channel factory pattern
(lambda) and listens for requests on the created channel for tasks. It calls the handleReq function.*/
func ComputeServer() chan *Task {
	reqs := make(chan *Task)
	go func() {
		for {
			req, ok := <-reqs
			if !ok {
				break
			}
			semRout <- 1
			go handleReq(req)
		}
	}()

	return reqs
}

/* A function that uses the channel factory pattern
(lambda) and listens for requests on the created channel for results to print to the console.*/
func DisplayServer() chan float32 {
	reqs := make(chan float32)
	go func() {
		for {
			req, ok := <-reqs
			if !ok {
				break
			}
			semDisp <- 1
			fmt.Printf("-------\nResult: %f\n-------\n", req)
			<-semDisp
			wgDisp.Done()
		}
	}()

	return reqs
}

func main() {
	dispChan := DisplayServer()
	reqChan := ComputeServer()
	defer close(dispChan)
	defer close(reqChan)

	for {
		var a, b float32
		// make sure to use semDisp
		// …
		wgDisp.Add(1)
		semDisp <- 1
		fmt.Print("Enter two numbers: ")
		fmt.Scanf("%f %f \n", &a, &b)
		fmt.Printf("%f %f \n", a, b)
		<-semDisp
		wgDisp.Done()
		if a == 0 && b == 0 {
			break
		}
		// Create task and send to ComputeServer
		wgDisp.Add(1)
		reqChan <- &Task{a, b, dispChan}
		// …
		time.Sleep(1e9)
	}
	// Don’t exit until all is done

	wgDisp.Wait()
	wgRout.Wait()
}
