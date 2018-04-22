package mock

type Tailor struct {
	AddCalled, RemoveCalled, MoveCalled, ReplaceCalled bool
	AddReturn, RemoveReturn, MoveReturn, ReplaceReturn error
}

func (m *Tailor) Add(obj interface{}, path string, value interface{}) error {
	m.AddCalled = true
	return m.AddReturn
}

func (m *Tailor) Remove(obj interface{}, path string) error {
	m.RemoveCalled = true
	return m.RemoveReturn
}

func (m *Tailor) Move(obj interface{}, path string, indexA uint64, indexB uint64) error {
	m.MoveCalled = true
	return m.MoveReturn
}

func (m *Tailor) Replace(obj interface{}, path string, value interface{}) error  {
	m.ReplaceCalled = true
	return m.ReplaceReturn
}