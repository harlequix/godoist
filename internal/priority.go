package internal

type PRIORITY_LEVEL int

const (
	HIGH     PRIORITY_LEVEL = 4
	MEDIUM   PRIORITY_LEVEL = 3
	LOW      PRIORITY_LEVEL = 2
	VERY_LOW PRIORITY_LEVEL = 1
)

func (p PRIORITY_LEVEL) String() string {
	switch p {
	case HIGH:
		return "High"
	case MEDIUM:
		return "Medium"
	case LOW:
		return "Low"
	case VERY_LOW:
		return "Very Low"
	default:
		return "Unknown"
	}
}
