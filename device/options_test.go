package device

import (
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/wrp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestOptionsDefault(t *testing.T) {
	assert := assert.New(t)

	for _, o := range []*Options{nil, new(Options)} {
		t.Log(o)

		assert.Equal(DefaultDeviceNameHeader, o.deviceNameHeader())
		assert.Equal(DefaultConveyHeader, o.conveyHeader())
		assert.Equal(DefaultDeviceMessageQueueSize, o.deviceMessageQueueSize())
		assert.Equal(DefaultHandshakeTimeout, o.handshakeTimeout())
		assert.Equal(DefaultInitialCapacity, o.initialCapacity())
		assert.Equal(DefaultIdlePeriod, o.idlePeriod())
		assert.Equal(DefaultPingPeriod, o.pingPeriod())
		assert.Equal(DefaultWriteTimeout, o.writeTimeout())
		assert.Equal(DefaultReadBufferSize, o.readBufferSize())
		assert.Equal(DefaultWriteBufferSize, o.writeBufferSize())
		assert.Empty(o.subprotocols())
		assert.NotNil(o.keyFunc())
		assert.NotNil(o.logger())
		assert.NotNil(o.messageListener())
		assert.NotNil(o.connectListener())
		assert.NotNil(o.disconnectListener())
		assert.NotNil(o.pongListener())
	}
}

func TestOptions(t *testing.T) {
	assert := assert.New(t)

	expectedLogger := &logging.LoggerWriter{os.Stdout}

	expectedKey := Key("TestOptions key")
	expectedKeyFunc := func(ID, Convey, *http.Request) (Key, error) {
		return expectedKey, nil
	}

	var called map[string]bool
	expectedMessageListener := func(Interface, *wrp.Message) { called["expectedMessageListener"] = true }
	expectedConnectListener := func(Interface) { called["expectedConnectListener"] = true }
	expectedDisconnectListener := func(Interface) { called["expectedDisconnectListener"] = true }
	expectedPongListener := func(Interface, string) { called["expectedPongListener"] = true }

	o := Options{
		DeviceNameHeader:       "X-TestOptions-Device-Name",
		ConveyHeader:           "X-TestOptions-Convey",
		HandshakeTimeout:       DefaultHandshakeTimeout + 12377123*time.Second,
		InitialCapacity:        DefaultInitialCapacity + 4719,
		ReadBufferSize:         DefaultReadBufferSize + 48729,
		WriteBufferSize:        DefaultWriteBufferSize + 926,
		Subprotocols:           []string{"foobar"},
		DeviceMessageQueueSize: DefaultDeviceMessageQueueSize + 287342,
		IdlePeriod:             DefaultIdlePeriod + 3472*time.Minute,
		PingPeriod:             DefaultPingPeriod + 384*time.Millisecond,
		WriteTimeout:           DefaultWriteTimeout + 327193*time.Second,
		KeyFunc:                expectedKeyFunc,
		Logger:                 expectedLogger,
		MessageListener:        expectedMessageListener,
		ConnectListener:        expectedConnectListener,
		DisconnectListener:     expectedDisconnectListener,
		PongListener:           expectedPongListener,
	}

	assert.Equal(o.DeviceNameHeader, o.deviceNameHeader())
	assert.Equal(o.ConveyHeader, o.conveyHeader())
	assert.Equal(o.DeviceMessageQueueSize, o.deviceMessageQueueSize())
	assert.Equal(o.HandshakeTimeout, o.handshakeTimeout())
	assert.Equal(o.InitialCapacity, o.initialCapacity())
	assert.Equal(o.IdlePeriod, o.idlePeriod())
	assert.Equal(o.PingPeriod, o.pingPeriod())
	assert.Equal(o.WriteTimeout, o.writeTimeout())
	assert.Equal(o.ReadBufferSize, o.readBufferSize())
	assert.Equal(o.WriteBufferSize, o.writeBufferSize())
	assert.Equal(o.Subprotocols, o.subprotocols())
	assert.Equal(expectedLogger, o.logger())

	expectedDevice := new(mockDevice)
	expectedMessage := new(wrp.Message)
	pongData := "some pong data"

	called = make(map[string]bool)
	o.messageListener()(expectedDevice, expectedMessage)
	assert.Len(called, 1)
	assert.True(called["expectedMessageListener"])

	called = make(map[string]bool)
	o.connectListener()(expectedDevice)
	assert.Len(called, 1)
	assert.True(called["expectedConnectListener"])

	called = make(map[string]bool)
	o.disconnectListener()(expectedDevice)
	assert.Len(called, 1)
	assert.True(called["expectedDisconnectListener"])

	called = make(map[string]bool)
	o.pongListener()(expectedDevice, pongData)
	assert.Len(called, 1)
	assert.True(called["expectedPongListener"])

	expectedDevice.AssertExpectations(t)

	actualKeyFunc := o.keyFunc()
	if assert.NotNil(actualKeyFunc) {
		actualKey, err := actualKeyFunc(ID(""), nil, nil)
		assert.Equal(expectedKey, actualKey)
		assert.Nil(err)
	}
}
