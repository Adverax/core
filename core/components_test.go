package core

type myClass struct {
}

var xxx = NewComponent(func() *myClass {
	return &myClass{}
})

func main() {
	// build components
	Components.Init()
	defer Components.Done()
	// run application
}
