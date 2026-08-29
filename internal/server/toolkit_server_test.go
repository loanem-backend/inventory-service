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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddToolkit(t *testing.T) {
	var sampleID int32 = 42
	sampleReq := &pbinventory.AddToolkitRequest{
		KitName:    "Toolkit Test",
		TotalCount: 100,
	}

	tests := []struct {
		name         string
		mockBehavior func(m *service_mock.MockToolkitService)
		input        *pbinventory.AddToolkitRequest
		assertCase   func(t *testing.T, resp *pbinventory.AddToolkitResponse, err error)
	}{
		{
			name: "Success_Created",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					Create(gomock.Any(), &entity.Toolkit{
						KitName:    sampleReq.GetKitName(),
						TotalCount: int(sampleReq.GetTotalCount()),
					}).
					Return(int16(sampleID), nil)
			},
			input: sampleReq,
			assertCase: func(t *testing.T, resp *pbinventory.AddToolkitResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, sampleID, resp.GetId())
			},
		},
		{
			name: "Failed_Internal",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(int16(0), status.Error(codes.Internal, ""))
			},
			input: sampleReq,
			assertCase: func(t *testing.T, resp *pbinventory.AddToolkitResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockToolkitService := service_mock.NewMockToolkitService(ctrl)
			test.mockBehavior(mockToolkitService)

			s := NewToolkitServer(mockToolkitService)

			resp, err := s.AddToolkit(t.Context(), test.input)

			test.assertCase(t, resp, err)
		})
	}
}

func TestGetAllToolkits(t *testing.T) {
	sampleToolkits := []*entity.Toolkit{
		{
			ID:              1,
			KitName:         "Toolkit 1",
			TotalCount:      50,
			OutOfOrderCount: 5,
			CreatedAt:       time.Date(2026, time.August, 7, 7, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2026, time.August, 7, 7, 0, 0, 0, time.UTC),
		},
		{
			ID:              2,
			KitName:         "Toolkit 2",
			TotalCount:      100,
			OutOfOrderCount: 10,
			CreatedAt:       time.Date(2026, time.August, 7, 7, 20, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2026, time.August, 9, 6, 50, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name         string
		mockBehavior func(m *service_mock.MockToolkitService)
		assertCase   func(t *testing.T, resp *pbinventory.GetAllToolkitsResponse, err error)
	}{
		{
			name: "Success_Fetched",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(sampleToolkits, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllToolkitsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				resultToolkits := resp.GetToolkits()
				assert.Equal(t, len(sampleToolkits), len(resultToolkits))

				for i := range resultToolkits {
					expected := &pbinventory.Toolkit{
						Id:              int32(sampleToolkits[i].ID),
						KitName:         sampleToolkits[i].KitName,
						TotalCount:      int32(sampleToolkits[i].TotalCount),
						OutOfOrderCount: int32(sampleToolkits[i].OutOfOrderCount),
						CreatedAt:       timestamppb.New(sampleToolkits[i].CreatedAt),
						UpdatedAt:       timestamppb.New(sampleToolkits[i].UpdatedAt),
					}
					diff := cmp.Diff(expected, resultToolkits[i], protocmp.Transform())
					if diff != "" {
						t.Errorf("Mapping result mismatch: %s\n", diff)
					}
				}
			},
		},
		{
			name: "Failed_Internal",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, status.Error(codes.Internal, ""))
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllToolkitsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "Success_EmptyWhenNil",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllToolkitsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetAllToolkitsResponse{
					Toolkits: []*pbinventory.Toolkit{},
				}
				diff := cmp.Diff(expected, resp, protocmp.Transform())
				if diff != "" {
					t.Errorf("Response mismatch: %s\n", diff)
				}
			},
		},
		{
			name: "Success_EmptyWhenNoRows",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return([]*entity.Toolkit{}, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetAllToolkitsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetAllToolkitsResponse{
					Toolkits: []*pbinventory.Toolkit{},
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

			mockToolkitService := service_mock.NewMockToolkitService(ctrl)
			test.mockBehavior(mockToolkitService)

			s := NewToolkitServer(mockToolkitService)

			resp, err := s.GetAllToolkits(t.Context(), &pbinventory.GetAllToolkitsRequest{})

			test.assertCase(t, resp, err)
		})
	}
}
