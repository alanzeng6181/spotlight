package algorithm

type Location struct {
	XVal Float
	YVal Float
}

func (l Location) X() Float {
	return l.XVal
}

func (l Location) Y() Float {
	return l.YVal
}
