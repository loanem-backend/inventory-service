package mapper

import (
	"testing"

	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
)

func TestIntToAddInstrumentResponse(t *testing.T) {
	sampleResponse := &pbinventory.AddInstrumentResponse{
		Id: 98765,
	}

	tests := []struct {
		name       string
		input      int32
		assertCase func(t *testing.T, result *pbinventory.AddInstrumentResponse)
	}{
		{
			name:  "Success",
			input: sampleResponse.Id,
			assertCase: func(t *testing.T, result *pbinventory.AddInstrumentResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, sampleResponse.Id, result.GetId())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IntToAddInstrumentResponse(test.input)

			test.assertCase(t, result)
		})
	}
}
