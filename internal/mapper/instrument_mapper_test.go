package mapper

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestInstrumentToPBInstrument(t *testing.T) {
	sampleInstrument := &entity.Instrument{
		ID:        234,
		Name:      "Instrument Test",
		Picture:   "picture-1",
		CreatedAt: time.Now().Add(-36 * time.Hour).UTC(),
		UpdatedAt: time.Now().Add(-12 * time.Hour).UTC(),
	}

	tests := []struct {
		name       string
		input      *entity.Instrument
		assertCase func(t *testing.T, result *pbinventory.Instrument)
	}{
		{
			name:  "NonNilInput",
			input: sampleInstrument,
			assertCase: func(t *testing.T, result *pbinventory.Instrument) {
				expected := &pbinventory.Instrument{
					Id:        int32(sampleInstrument.ID),
					Name:      sampleInstrument.Name,
					Picture:   sampleInstrument.Picture,
					CreatedAt: timestamppb.New(sampleInstrument.CreatedAt),
					UpdatedAt: timestamppb.New(sampleInstrument.UpdatedAt),
				}

				diff := cmp.Diff(expected, result, protocmp.Transform())
				if diff != "" {
					t.Errorf("Mapping result mismatch: %s\n", diff)
				}
			},
		},
		{
			name:  "NilInput",
			input: nil,
			assertCase: func(t *testing.T, result *pbinventory.Instrument) {
				expected := &pbinventory.Instrument{}

				diff := cmp.Diff(expected, result, protocmp.Transform())
				if diff != "" {
					t.Errorf("Mapping result mismatch: %s\n", diff)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := InstrumentToPBInstrument(test.input)

			test.assertCase(t, result)
		})
	}
}
