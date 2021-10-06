package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"patiencedesert/bromel/meshwork"
	"testing"
)

func TestCustomPacker(t *testing.T) {
	packer := &CustomPacker{}

	entry := &meshwork.Entry{
		ID:   "test",
		Data: []byte("data"),
	}
	msg, err := packer.Pack(entry)
	assert.NoError(t, err)
	assert.NotNil(t, msg)

	r := bytes.NewBuffer(msg)
	newEntry, err := packer.Unpack(r)
	assert.NoError(t, err)
	assert.NotNil(t, newEntry)
	assert.Equal(t, newEntry.ID, entry.ID)
	assert.Equal(t, newEntry.Data, entry.Data)
}
