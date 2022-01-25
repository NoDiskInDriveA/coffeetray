package coffeetray

func AssertNoError(e error) {
	if e != nil {
		panic(e)
	}
}
