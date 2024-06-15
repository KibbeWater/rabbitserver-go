package rabbit

type GnuRand struct {
	r [344]uint32
	n int
}

func Srand(seed uint32) *GnuRand {
	if seed == 0 {
		seed = 1
	}
	r := [344]uint32{seed}
	for i := 1; i < 31; i++ {
		r[i] = (16807 * r[i-1]) % 0x7fffffff
	}
	for i := 31; i < 34; i++ {
		r[i] = r[i-31]
	}
	for i := 34; i < 344; i++ {
		r[i] = (r[i-31] + r[i-3]) & 0xffffffff
	}

	return &GnuRand{r: r, n: 0}
}

func (gr *GnuRand) Rand() uint32 {
	x := (gr.r[(gr.n+313)%344] + gr.r[(gr.n+341)%344]) & 0xffffffff
	gr.r[gr.n%344] = x
	gr.n = (gr.n + 1) % 344
	return x >> 1
}
