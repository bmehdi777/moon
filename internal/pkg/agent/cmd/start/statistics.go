package start

type Statistics map[string]int

func (s *Statistics) Set(path string, value int) {
	(*s)[path] = value
}

func (s *Statistics) Increment(path string) {
	(*s)[path] += 1
}

func (s *Statistics) Delete(path string) {
	delete(*s, path)
}

func (s *Statistics) Get(path string) int {
	if stat, found := (*s)[path]; found {
		return stat
	}

	return -1
}
