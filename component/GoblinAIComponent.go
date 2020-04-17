package component

// GoblinAIComponent .
type GoblinAIComponent struct {
	Name            string
	Energy          int
	SightRange      int
	HungerThreshold int
	SocialThreshold int
	MateThreshold   int
	State           string
	TargetX         int
	TargetY         int
}

func (pc GoblinAIComponent) GetType() string {
	return "GoblinAIComponent"
}
