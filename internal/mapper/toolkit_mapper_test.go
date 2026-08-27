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

func TestToolkitToPBToolkit(t *testing.T) {
	sampleToolkit := &entity.Toolkit{
		ID:      5678,
		KitName: "Kit Test",
		Course: entity.Course{
			ID: 234,
		},
		TotalCount:      20,
		OutOfOrderCount: 3,
		CreatedAt:       time.Date(2026, time.July, 25, 14, 20, 5, 0, time.UTC),
		UpdatedAt:       time.Date(2026, time.July, 25, 14, 20, 5, 0, time.UTC),
	}

	tests := []struct {
		name       string
		input      *entity.Toolkit
		assertCase func(t *testing.T, result *pbinventory.Toolkit)
	}{
		{
			name:  "Success",
			input: sampleToolkit,
			assertCase: func(t *testing.T, result *pbinventory.Toolkit) {
				assert.NotNil(t, result)
				assert.Equal(t, int32(sampleToolkit.ID), result.GetId())
				assert.Equal(t, sampleToolkit.KitName, result.GetKitName())
				assert.Equal(t, int32(sampleToolkit.Course.ID), result.GetCourseId())
				assert.Equal(t, int32(sampleToolkit.TotalCount), result.GetTotalCount())
				assert.Equal(t, int32(sampleToolkit.OutOfOrderCount), result.GetOutOfOrderCount())
				assert.Equal(t, sampleToolkit.CreatedAt, result.GetCreatedAt().AsTime())
				assert.Equal(t, sampleToolkit.UpdatedAt, result.GetUpdatedAt().AsTime())
			},
		},
		{
			name:  "NilInput",
			input: nil,
			assertCase: func(t *testing.T, result *pbinventory.Toolkit) {
				assert.Nil(t, result)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ToolkitToPBToolkit(test.input)

			test.assertCase(t, result)
		})
	}
}

func TestToolkitsToGetAllToolkitsResponse(t *testing.T) {
	sampleToolkits := []*entity.Toolkit{
		{
			ID:      12,
			KitName: "Kit Test 1",
			Course: entity.Course{
				ID: 7,
			},
			TotalCount:      10,
			OutOfOrderCount: 0,
			CreatedAt:       time.Date(2026, time.July, 20, 8, 30, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2026, time.August, 4, 10, 15, 0, 0, time.UTC),
		},
		nil,
		{
			ID:      23,
			KitName: "Kit Test 2",
			Course: entity.Course{
				ID: 8,
			},
			TotalCount:      20,
			OutOfOrderCount: 3,
			CreatedAt:       time.Date(2026, time.July, 26, 8, 30, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2026, time.July, 26, 8, 30, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name       string
		input      []*entity.Toolkit
		assertCase func(t *testing.T, result *pbinventory.GetAllToolkitsResponse)
	}{
		{
			name:  "NonNilInput",
			input: sampleToolkits,
			assertCase: func(t *testing.T, result *pbinventory.GetAllToolkitsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Toolkits)
				resultToolkits := result.GetToolkits()

				expectedToolkits := make([]*pbinventory.Toolkit, 0)
				for _, t := range sampleToolkits {
					if t != nil {
						expectedToolkits = append(expectedToolkits, &pbinventory.Toolkit{
							Id:              int32(t.ID),
							KitName:         t.KitName,
							CourseId:        int32(t.Course.ID),
							TotalCount:      int32(t.TotalCount),
							OutOfOrderCount: int32(t.OutOfOrderCount),
							CreatedAt:       timestamppb.New(t.CreatedAt),
							UpdatedAt:       timestamppb.New(t.UpdatedAt),
						})
					}
				}

				assert.Equal(t, len(expectedToolkits), len(resultToolkits))

				diff := cmp.Diff(expectedToolkits, resultToolkits, protocmp.Transform())
				if diff != "" {
					t.Errorf("Mapping result mismatch: %s\n", diff)
				}
			},
		},
		{
			name:  "NilInput",
			input: nil,
			assertCase: func(t *testing.T, result *pbinventory.GetAllToolkitsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Toolkits)
				assert.Empty(t, result.Toolkits)

				expected := &pbinventory.GetAllToolkitsResponse{
					Toolkits: []*pbinventory.Toolkit{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
		{
			name:  "EmptyNotNilInput",
			input: []*entity.Toolkit{},
			assertCase: func(t *testing.T, result *pbinventory.GetAllToolkitsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Toolkits)
				assert.Empty(t, result.Toolkits)

				expected := &pbinventory.GetAllToolkitsResponse{
					Toolkits: []*pbinventory.Toolkit{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
		{
			name:  "AllNilInput",
			input: []*entity.Toolkit{nil, nil},
			assertCase: func(t *testing.T, result *pbinventory.GetAllToolkitsResponse) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Toolkits)
				assert.Empty(t, result.Toolkits)

				expected := &pbinventory.GetAllToolkitsResponse{
					Toolkits: []*pbinventory.Toolkit{},
				}
				if !proto.Equal(expected, result) {
					t.Errorf("Mapping result mismatch")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ToolkitsToGetAllToolkitsResponse(test.input)

			test.assertCase(t, result)
		})
	}
}
