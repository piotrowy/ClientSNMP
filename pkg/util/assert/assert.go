package assert

type Predicate func() bool

func AssertP(p func() bool) {
	if !p() {
		panic("assertion failed")
	}
}

func Assert(p bool) {
	if !p {
		panic("assertion failed")
	}
}
