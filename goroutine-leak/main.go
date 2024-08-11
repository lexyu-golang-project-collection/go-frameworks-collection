package main

func Leak() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}

func main() {
	Leak()
}
