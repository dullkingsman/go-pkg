package writable

type Formatter interface {
	Format(value Writable) string
}
