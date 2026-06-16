package mapper

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
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
				assert.Nil(t, result)
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

func TestInstrumentsToGetAllInstrumentsResponse(t *testing.T) {
	sampleInstruments := []*entity.Instrument{
		{
			ID:        9876,
			Name:      "Instrument 1",
			Picture:   "picture-1",
			CreatedAt: time.Now().Add(-36 * time.Hour).UTC(),
			UpdatedAt: time.Now().Add(-12 * time.Hour).UTC(),
		},
		nil,
		{
			ID:        6789,
			Name:      "Instrument 3",
			Picture:   "picture-3",
			CreatedAt: time.Now().Add(-24 * time.Hour).UTC(),
			UpdatedAt: time.Now().Add(-6 * time.Hour).UTC(),
		},
	}

	tests := []struct {
		name       string
		input      []*entity.Instrument
		assertCase func(t *testing.T, result *pbinventory.GetAllInstrumentsResponse)
	}{
		{
			name:  "NonNilInput",
			input: sampleInstruments,
			assertCase: func(t *testing.T, result *pbinventory.GetAllInstrumentsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Instruments)
				resultInstruments := result.GetInstruments()

				expectedInstruments := make([]*pbinventory.Instrument, 0)
				for _, i := range sampleInstruments {
					if i != nil {
						expectedInstruments = append(expectedInstruments, &pbinventory.Instrument{
							Id:        int32(i.ID),
							Name:      i.Name,
							Picture:   i.Picture,
							CreatedAt: timestamppb.New(i.CreatedAt),
							UpdatedAt: timestamppb.New(i.UpdatedAt),
						})
					}
				}

				assert.Equal(t, len(expectedInstruments), len(resultInstruments))

				diff := cmp.Diff(expectedInstruments, resultInstruments, protocmp.Transform())
				if diff != "" {
					t.Errorf("Mapping result mismatch: %s\n", diff)
				}
			},
		},
		{
			name:  "NilInput",
			input: nil,
			assertCase: func(t *testing.T, result *pbinventory.GetAllInstrumentsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Instruments)
				assert.Empty(t, result.Instruments)

				expected := &pbinventory.GetAllInstrumentsResponse{
					Instruments: []*pbinventory.Instrument{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
		{
			name:  "EmptyNotNilInput",
			input: []*entity.Instrument{},
			assertCase: func(t *testing.T, result *pbinventory.GetAllInstrumentsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Instruments)
				assert.Empty(t, result.Instruments)

				expected := &pbinventory.GetAllInstrumentsResponse{
					Instruments: []*pbinventory.Instrument{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
		{
			name:  "AllNilInput",
			input: []*entity.Instrument{nil, nil},
			assertCase: func(t *testing.T, result *pbinventory.GetAllInstrumentsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Instruments)
				assert.Empty(t, result.Instruments)

				expected := &pbinventory.GetAllInstrumentsResponse{
					Instruments: []*pbinventory.Instrument{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := InstrumentsToGetAllInstrumentsResponse(test.input)

			test.assertCase(t, result)
		})
	}
}
