package main

import (
	"errors"
	"fmt"
)

// Trip struct
type Trip struct {
	destination string  //destinatino of the trip
	weight      float32 //weight of the load
	deadline    int     // in hours
}

// Vehicle struct
type Vehicle struct {
	vehicle     string
	name        string
	destination string
	speed       float32
	capacity    float32
	load        float32
}

// Truck struct
type Truck struct {
	Vehicle
}

// NewTruck function
func NewTruck() Truck {
	return Truck{Vehicle{"Truck", "Truck", "", 40, 10, 0}}
}

// Pickup struct
type Pickup struct {
	Vehicle
	isPrivate bool
}

// NewPickUp function
func NewPickUp() Pickup {
	return Pickup{Vehicle{"Pickup", "Pickup", "", 60, 2, 0}, true}
}

// TrainCar struct
type TrainCar struct {
	Vehicle
	railway string
}

// NewTrainCar function
func NewTrainCar() TrainCar {
	return TrainCar{Vehicle{"TrainCar", "TrainCar", "", 30, 30, 0}, "CNR"}
}

// Transporter interface
type Transporter interface {
	/* Description
	 Returning a error if the transporter has insufficient
	capacity to carry the weight, has a different destination or cannot make the destination on
	time. If the current destination is empty, the destination needs to be updated to the trip’s
	destination.
	*/
	addLoad(Trip) error

	/* Example:
	Truck B to Montreal with 8.000000 tons
	Pickup A to with 0.000000 tons (Private: true)
	TrainCar A to Montreal with 8.000000 tons (CNR)
	*/
	print()
}

// NewTorontoTrip function
func NewTorontoTrip(weight float32, deadline int) *Trip {
	var x Trip
	x.destination = "Toronto"
	return &x
}

// NewMontrealTrip function
func NewMontrealTrip(weight float32, deadline int) *Trip {
	var x Trip
	x.destination = "Montreal"
	return &x
}

func (p *Vehicle) addLoad(trip Trip) error {

	if p.destination != trip.destination {
		if p.destination == "" {
			p.destination = trip.destination
		} else {
			return errors.New("Error: Other destination")
		}
	}
	if p.capacity-p.load < trip.weight {
		return errors.New("Error: Out of capacity ")
	}
	if trip.destination == "Toronto" {
		if float32(trip.deadline) < 400.0*(1.0/p.speed) {
			return errors.New("Error: Vehicle not fast enough")
		}
	} else if trip.destination == "Montreal" {
		if float32(trip.deadline) < 200.0*(1.0/p.speed) {
			return errors.New("Error: Vehicle not fast enough")
		}
	}
	p.load = p.load + trip.weight
	return nil
}

func (p *Truck) print() {
	fmt.Printf("Truck "+p.name+" to "+p.destination+" with %f tons\n", p.load)
}

func (p *Pickup) print() {

	fmt.Printf("Pickup "+p.name+" to "+p.destination+" with %f tons (Private: %t)\n", p.load, p.isPrivate)
}

func (p *TrainCar) print() {
	fmt.Printf("TrainCar "+p.name+" to "+p.destination+" with %f tons ("+p.railway+")\n", p.load)
}

func main() {
	truck1 := NewTruck()
	truck1.name = "A"
	truck2 := NewTruck()
	truck2.name = "B"
	pickup1 := NewPickUp()
	pickup1.name = "A"
	pickup2 := NewPickUp()
	pickup2.name = "B"
	pickup3 := NewPickUp()
	pickup3.name = "C"
	traincar1 := NewTrainCar()
	traincar1.name = "A"
	vehicles := [6]Transporter{&truck1, &truck2, &pickup1, &pickup2, &pickup3, &traincar1}
	var trips []Trip

running:
	for {
		var dest string
		var weight float32
		var deadline int
		//get destination
		fmt.Printf("Destination: (t)oronto, (m)ontreal, else exit? ")
		fmt.Scanf("%s \n", &dest)
		runes := []rune(dest)
		// ... Convert back into a string from rune slice.
		dest = string(runes[0:1])

		if dest == "t" {
			dest = "Toronto"
		} else if dest == "m" {
			dest = "Montreal"
		} else {
			fmt.Printf("Not going to TO or Montreal, bye!\n")
			temp := "Trips: [ "
			for i := range trips {
				w := fmt.Sprint(trips[i].weight)
				d := fmt.Sprint(trips[i].deadline)
				temp = temp + "{ " + trips[i].destination + " " + w + " " + d + " }"
			}
			fmt.Printf(temp + " ]\n")
			fmt.Printf("Vehicles:\n")
			for i := range vehicles {
				vehicles[i].print()
			}

			break running
		}

		//get weight
		fmt.Printf("Weight: ")
		fmt.Scanf("%f \n", &weight)
		//get deadline
		fmt.Printf("Deadline (in hours): ")
		fmt.Scanf("%d \n", &deadline)

		added := false
	checking:

		for i := range vehicles {
			err := vehicles[i].addLoad(Trip{dest, weight, deadline})
			if err == nil {
				trips = append(trips, Trip{dest, weight, deadline})
				vehicles[i].print()
				added = true
				break checking

			} else {
				fmt.Printf(err.Error() + "\n")

			}

		}
		if !added {
			fmt.Printf("Unable to add trip\n")
		}

	}
	//still unclear

}
