package mapper

import (
	"testing"

	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
)

func TestAddToolkitRequestToToolkit(t *testing.T) {
	sampleReq := &pbinventory.AddToolkitRequest{
		KitName:    "Kit Test",
		TotalCount: 20,
	}

	tests := []struct {
		name       string
		input      *pbinventory.AddToolkitRequest
		assertCase func(t *testing.T, result *entity.Toolkit)
	}{
		{
			name:  "Success",
			input: sampleReq,
			assertCase: func(t *testing.T, result *entity.Toolkit) {
				assert.NotNil(t, result)
				assert.Equal(t, sampleReq.GetKitName(), result.KitName)
				assert.Equal(t, int(sampleReq.GetTotalCount()), result.TotalCount)
			},
		},
		{
			name:  "NilInput",
			input: nil,
			assertCase: func(t *testing.T, result *entity.Toolkit) {
				assert.Nil(t, result)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := AddToolkitRequestToToolkit(test.input)

			test.assertCase(t, result)
		})
	}
}

func TestIntToAddToolkitResponse(t *testing.T) {
	var sampleID int16 = 678

	tests := []struct {
		name       string
		input      int16
		assertCase func(t *testing.T, result *pbinventory.AddToolkitResponse)
	}{
		{
			name:  "Success",
			input: sampleID,
			assertCase: func(t *testing.T, result *pbinventory.AddToolkitResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, int32(sampleID), result.GetId())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IntToAddToolkitResponse(test.input)

			test.assertCase(t, result)
		})
	}
}
