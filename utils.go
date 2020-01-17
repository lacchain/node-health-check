package main

func divide(a, b float64) float64 {
	if b > 0 {
		return a / b
	}
	return 0.0
}

func fanIn(i1, i2 <-chan bool) <-chan bool {
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
