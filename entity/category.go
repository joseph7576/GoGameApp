package entity

type Category string

const (
	CategoryFootball Category = "football"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryFootball:
		return true
	default:
		return false
	}
}
