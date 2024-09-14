package logrus

type DummyPurifier struct {
}

func (f *DummyPurifier) Purify(original, derivative string) string {
	return derivative
}
