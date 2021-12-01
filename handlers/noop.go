package handlers

func Noop() func(err error) {
	return func(err error) {}
}
