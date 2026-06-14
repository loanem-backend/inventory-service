package server

import (
	"testing"

	service_mock "github.com/loanem-backend/inventory-service/internal/mocks/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
