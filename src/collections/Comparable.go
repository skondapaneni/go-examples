package collections

type Comparable interface {
    Compare(b Comparable) bool
}
