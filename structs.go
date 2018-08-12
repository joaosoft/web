package types

type Date string
type Time string
type Status string
type Day string
type ListDay []Day

func (l *ListDay) Contains(day Day) bool {
	for _, value := range *l {
		if day == Day(value) {
			return true
		}
	}

	return false
}
