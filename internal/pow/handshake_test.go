package pow

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandshake(t *testing.T) {
	var difficulty byte = 10
	randTokenSize := 64

	powReceive := NewReceiver(difficulty, randTokenSize)

	message := "hello"

	clientConn, serverConn := net.Pipe()

	go func() {
		_, receiveErr := powReceive(serverConn)
		assert.NoError(t, receiveErr)
		serverConn.Write([]byte(message))
		serverConn.Close()
	}()

	receivedDifficulty, _, establishErr := Establish(clientConn)
	assert.NoError(t, establishErr)
	assert.Equal(t, difficulty, receivedDifficulty)

	receivedData, err := ioutil.ReadAll(clientConn)
	assert.NoError(t, err)
	receivedMessage := string(receivedData)
	assert.Equal(t, message, receivedMessage)
}
