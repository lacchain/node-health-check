package main

func divide(a, b float64) float64 {
	if b > 0 {
		return a / b
	}
	return 0.0
}

func fanIn(i1, i2 <-chan bool) <-chan bool { //this func only returns a receiver channel
	c := make(chan bool)
	go func() {
		for {
			c <- <-i1
		}
	}()

	go func() {
		for {
			c <- <-i2
		}
	}()

	return c
}

// func fanIn(i1 <-chan string) <-chan string { //this func only returns a receiver channel
// 	c := make(chan string)
// 	go func() {
// 		for {
// 			c <- <-i1 //<-i1 means extracting the value of  channel i1; c<- passing the value to channel 'c'
// 		}
// 	}()

// 	// go func() {
// 	// 	for {
// 	// 		c <- <-i2
// 	// 	}
// 	// }()

// 	return c
// }
