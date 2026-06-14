package server

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/loanem-backend/inventory-service/internal/entity"
	service_mock "github.com/loanem-backend/inventory-service/internal/mocks/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAddInstrument(t *testing.T) {
	sampleInstrument := &struct {
		id   int32
		name string
	}{
		id:   234,
		name: "Instrument Test",
	}

	tests := []struct {
		name         string
		mockBehavior func(m *service_mock.MockInstrumentService)
		inputName    string
		assertCase   func(t *testing.T, resp *pbinventory.AddInstrumentResponse, err error)
	}{
		{
			name: "Success_Created",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					AddInstrument(gomock.Any(), sampleInstrument.name).
					Return(int32(sampleInstrument.id), nil)
			},
			inputName: sampleInstrument.name,
			assertCase: func(t *testing.T, resp *pbinventory.AddInstrumentResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, sampleInstrument.id, resp.GetId())
			},
		},
		{
			name: "Failed_Internal",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					AddInstrument(gomock.Any(), sampleInstrument.name).
					Return(int32(0), status.Error(codes.Internal, ""))
			},
			inputName: sampleInstrument.name,
			assertCase: func(t *testing.T, resp *pbinventory.AddInstrumentResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockInstrumentService := service_mock.NewMockInstrumentService(ctrl)
			test.mockBehavior(mockInstrumentService)

			s := NewInstrumentServer(mockInstrumentService)

			resp, err := s.AddInstrument(t.Context(), &pbinventory.AddInstrumentRequest{
				Name: test.inputName,
			})

			test.assertCase(t, resp, err)
		})
	}
}

func TestGetAllInstruments(t *testing.T) {
	sampleInstruments := []*entity.Instrument{
		{
			ID:        9876,
			Name:      "Instrument 1",
			Picture:   "picture-1",
			CreatedAt: time.Now().Add(-36 * time.Hour).UTC(),
			UpdatedAt: time.Now().Add(-12 * time.Hour).UTC(),
		},
		{
			ID:        6789,
			Name:      "Instrument 2",
			Picture:   "picture-2",
			CreatedAt: time.Now().Add(-24 * time.Hour).UTC(),
			UpdatedAt: time.Now().Add(-6 * time.Hour).UTC(),
		},
	}

	tests := []struct {
		name         string
		mockBehavior func(m *service_mock.MockInstrumentService)
		assertCase   func(t *testing.T, resp *pbinventory.GetAllInstrumentsResponse, err error)
	}{
		{
			name: "Success_Fetched",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					GetAllInstruments(gomock.Any()).
					Return(sampleInstruments, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(sampleInstruments), len(resp.GetInstruments()))

				lastIndex := len(sampleInstruments) - 1
				expectedTime := sampleInstruments[lastIndex].CreatedAt
				actualTime := resp.GetInstruments()[lastIndex].GetCreatedAt().AsTime()
				if !expectedTime.Equal(actualTime) {
					t.Errorf("Expected time %v, got %v", expectedTime, actualTime)
				}
			},
		},
		{
			name: "Failed_Internal",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					GetAllInstruments(gomock.Any()).
					Return(nil, status.Error(codes.Internal, ""))
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllInstrumentsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "Success_EmptyWhenNil",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					GetAllInstruments(gomock.Any()).
					Return(nil, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetAllInstrumentsResponse{
					Instruments: []*pbinventory.Instrument{},
				}
				diff := cmp.Diff(expected, resp, protocmp.Transform())
				if diff != "" {
					t.Errorf("Response mismatch: %s\n", diff)
				}
			},
		},
		{
			name: "Success_EmptyWhenNoRows",
			mockBehavior: func(m *service_mock.MockInstrumentService) {
				m.EXPECT().
					GetAllInstruments(gomock.Any()).
					Return([]*entity.Instrument{}, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetAllInstrumentsResponse{
					Instruments: []*pbinventory.Instrument{},
				}
				diff := cmp.Diff(expected, resp, protocmp.Transform())
				if diff != "" {
					t.Errorf("Response mismatch: %s\n", diff)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockInstrumentService := service_mock.NewMockInstrumentService(ctrl)
			test.mockBehavior(mockInstrumentService)

			s := NewInstrumentServer(mockInstrumentService)

			resp, err := s.GetAllInstruments(t.Context(), nil)

			test.assertCase(t, resp, err)
		})
	}
}
