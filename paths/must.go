package paths

// must panics if err is an error.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// must1 panics if err is an error. It is like must, but passes one argument through.
func must1(val interface{}, err error) interface{} {
	must(err)
	return val
}
