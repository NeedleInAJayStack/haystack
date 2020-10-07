package haystack

// Val represents a haystack tag value.
type Val interface {
	ToZinc() string
}
