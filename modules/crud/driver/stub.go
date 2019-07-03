package driver

type stub struct {
	c     Crud
	count int
}

func newStub(c Crud) *stub {
	return &stub{c: c, count: 1}
}

func (s *stub) addCount() {
	s.count++
}

func (s *stub) subtractCount() {
	s.count--
}

func (s *stub) getCount() int {
	return s.count
}

func (s *stub) getCrud() Crud {
	return s.c
}