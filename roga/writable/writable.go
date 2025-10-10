package writable

type Writable interface {
	String(f Formatter) string
	Json() ([]byte, error)
}
