package component

// GoblinAIComponent .
type GoblinAIComponent struct {
	Energy          int
	SightRange      int
	HungerThreshold int
	SocialThreshold int
	State           string
	TargetX         int
	TargetY         int
}

func (pc GoblinAIComponent) GetType() string {
	return "GoblinAIComponent"
}
