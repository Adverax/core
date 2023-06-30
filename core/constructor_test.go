package core

type myClass struct {
}

var xxx = Constructor(func() *myClass {
	return &myClass{}
})
