package component

// DescriptionComponent .
type DescriptionComponent struct {
	Name string
}

func (pc DescriptionComponent) GetType() string {
	return "DescriptionComponent"
}
