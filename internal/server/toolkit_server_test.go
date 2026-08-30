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

func TestGetToolkitsWithInstruments(t *testing.T) {
	sampleToolkits := []*entity.Toolkit{
		{
			ID:              1,
			KitName:         "Toolkit 1",
			TotalCount:      50,
			OutOfOrderCount: 5,
			Course: entity.Course{
				ID: 1,
			},
			CreatedAt: time.Date(2026, time.August, 7, 7, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, time.August, 7, 7, 0, 0, 0, time.UTC),
			Instruments: []entity.Instrument{
				{
					ID:        1,
					Name:      "Instrument 1",
					Picture:   "pic1.jpg",
					CreatedAt: time.Date(2026, time.August, 5, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2026, time.August, 5, 10, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					Name:      "Instrument 2",
					Picture:   "pic2.jpg",
					CreatedAt: time.Date(2026, time.August, 6, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2026, time.August, 6, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID:              2,
			KitName:         "Toolkit 2",
			TotalCount:      100,
			OutOfOrderCount: 10,
			Course: entity.Course{
				ID: 2,
			},
			CreatedAt: time.Date(2026, time.August, 7, 7, 20, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, time.August, 9, 6, 50, 0, 0, time.UTC),
			Instruments: []entity.Instrument{
				{
					ID:        3,
					Name:      "Instrument 3",
					Picture:   "pic3.jpg",
					CreatedAt: time.Date(2026, time.August, 1, 14, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2026, time.August, 1, 14, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	tests := []struct {
		name         string
		mockBehavior func(m *service_mock.MockToolkitService)
		assertCase   func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error)
	}{
		{
			name: "Success_FetchedWithInstruments",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetWithInstruments(gomock.Any()).
					Return(sampleToolkits, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				resultToolkits := resp.GetToolkits()
				assert.Equal(t, len(sampleToolkits), len(resultToolkits))

				for i := range resultToolkits {
					assert.Equal(t, int32(sampleToolkits[i].ID), resultToolkits[i].Toolkit.GetId())
					assert.Equal(t, sampleToolkits[i].KitName, resultToolkits[i].Toolkit.GetKitName())
					assert.Equal(t, int32(sampleToolkits[i].TotalCount), resultToolkits[i].Toolkit.GetTotalCount())
					assert.Equal(t, int32(sampleToolkits[i].OutOfOrderCount), resultToolkits[i].Toolkit.GetOutOfOrderCount())
					assert.Equal(t, len(sampleToolkits[i].Instruments), len(resultToolkits[i].GetInstruments()))
				}
			},
		},
		{
			name: "Failed_InternalError",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetWithInstruments(gomock.Any()).
					Return(nil, status.Error(codes.Internal, "database connection failed"))
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "Success_EmptyWhenNil",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				m.EXPECT().
					GetWithInstruments(gomock.Any()).
					Return(nil, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetToolkitsWithInstrumentsResponse{
					Toolkits: []*pbinventory.ToolkitWithInstruments{},
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
					GetWithInstruments(gomock.Any()).
					Return([]*entity.Toolkit{}, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				expected := &pbinventory.GetToolkitsWithInstrumentsResponse{
					Toolkits: []*pbinventory.ToolkitWithInstruments{},
				}
				diff := cmp.Diff(expected, resp, protocmp.Transform())
				if diff != "" {
					t.Errorf("Response mismatch: %s\n", diff)
				}
			},
		},
		{
			name: "Success_ToolkitWithoutInstruments",
			mockBehavior: func(m *service_mock.MockToolkitService) {
				toolkitWithoutInstruments := []*entity.Toolkit{
					{
						ID:              3,
						KitName:         "Toolkit 3",
						TotalCount:      75,
						OutOfOrderCount: 0,
						Course: entity.Course{
							ID:   3,
							Name: "Course 3",
							Year: 2026,
						},
						CreatedAt:   time.Date(2026, time.August, 8, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2026, time.August, 8, 0, 0, 0, 0, time.UTC),
						Instruments: []entity.Instrument{},
					},
				}
				m.EXPECT().
					GetWithInstruments(gomock.Any()).
					Return(toolkitWithoutInstruments, nil)
			},
			assertCase: func(t *testing.T, resp *pbinventory.GetToolkitsWithInstrumentsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				resultToolkits := resp.GetToolkits()

				assert.Equal(t, 1, len(resultToolkits))
				assert.Equal(t, 0, len(resultToolkits[0].GetInstruments()))
				assert.Equal(t, "Toolkit 3", resultToolkits[0].Toolkit.GetKitName())
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

			resp, err := s.GetToolkitsWithInstruments(t.Context(), &pbinventory.GetToolkitsWithInstrumentsRequest{})

			test.assertCase(t, resp, err)
		})
	}
}
