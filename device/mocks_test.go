package device

import (
	"github.com/Comcast/webpa-common/wrp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// mockRandom provides an io.Reader mock for a source of random bytes
type mockRandom struct {
	mock.Mock
}

func (m *mockRandom) Read(b []byte) (int, error) {
	arguments := m.Called(b)
	return arguments.Int(0), arguments.Error(1)
}

// mockDevice mocks the Interface type
type mockDevice struct {
	mock.Mock
}

func (m *mockDevice) ID() ID {
	arguments := m.Called()
	return arguments.Get(0).(ID)
}

func (m *mockDevice) Key() Key {
	arguments := m.Called()
	return arguments.Get(0).(Key)
}

func (m *mockDevice) Convey() Convey {
	arguments := m.Called()
	return arguments.Get(0).(Convey)
}

func (m *mockDevice) ConnectedAt() time.Time {
	arguments := m.Called()
	return arguments.Get(0).(time.Time)
}

func (m *mockDevice) RequestClose() {
	m.Called()
}

func (m *mockDevice) Closed() bool {
	arguments := m.Called()
	return arguments.Bool(0)
}

func (m *mockDevice) Send(message *wrp.Message) error {
	arguments := m.Called(message)
	return arguments.Error(0)
}

// deviceSet is a convenient map type for capturing visited devices
// and asserting expectations.
type deviceSet map[*device]bool

func (s deviceSet) len() int {
	return len(s)
}

func (s deviceSet) add(d Interface) {
	s[d.(*device)] = true
}

func (s *deviceSet) reset() {
	*s = make(deviceSet)
}

// managerCapture returns a high-level visitor for Manager testing
func (s deviceSet) managerCapture() func(Interface) {
	return func(d Interface) {
		s.add(d)
	}
}

// registryCapture returns a low-level visitor for registry testing
func (s deviceSet) registryCapture() func(*device) {
	return func(d *device) {
		s[d] = true
	}
}

func (s deviceSet) assertSameID(assert *assert.Assertions, expected ID) {
	for d, _ := range s {
		assert.Equal(expected, d.ID())
	}
}

func (s deviceSet) assertDistributionOfIDs(assert *assert.Assertions, expected map[ID]int) {
	actual := make(map[ID]int, len(expected))
	for d, _ := range s {
		actual[d.ID()] += 1
	}

	assert.Equal(expected, actual)
}

// drain copies whatever is available on the given channel into this device set
func (s deviceSet) drain(source <-chan Interface) {
	for d := range source {
		s.add(d)
	}
}

func expectsDevices(devices ...*device) deviceSet {
	result := make(deviceSet, len(devices))
	for _, d := range devices {
		result[d] = true
	}

	return result
}
