package component

// PoisonousComponent .
type PoisonousComponent struct {
	Duration int
}

func (pc PoisonousComponent) GetType() string {
	return "PoisonousComponent"
}
