package component

// MyTurnComponent .
type AppearanceComponent struct {
	SpriteX  int32
	SpriteY  int32
	Resource string
}

func (pc AppearanceComponent) GetType() string {
	return "AppearanceComponent"
}
