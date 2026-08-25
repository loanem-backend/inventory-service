package mapper

import (
	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AddToolkitRequestToToolkit(req *pbinventory.AddToolkitRequest) *entity.Toolkit {
	if req == nil {
		return nil
	}

	return &entity.Toolkit{
		KitName:    req.GetKitName(),
		TotalCount: int(req.GetTotalCount()),
	}
}

func IntToAddToolkitResponse(id int16) *pbinventory.AddToolkitResponse {
	return &pbinventory.AddToolkitResponse{
		Id: int32(id),
	}
}

func ToolkitToPBToolkit(t *entity.Toolkit) *pbinventory.Toolkit {
	if t == nil {
		return nil
	}

	return &pbinventory.Toolkit{
		Id:              int32(t.ID),
		KitName:         t.KitName,
		CourseId:        int32(t.Course.ID),
		TotalCount:      int32(t.TotalCount),
		OutOfOrderCount: int32(t.TotalCount),
		CreatedAt:       timestamppb.New(t.CreatedAt),
		UpdatedAt:       timestamppb.New(t.UpdatedAt),
	}
}

func ToolkitsToGetAllToolkitsResponse(toolkits []*entity.Toolkit) *pbinventory.GetAllToolkitsResponse {
	pbToolkits := make([]*pbinventory.Toolkit, 0, len(toolkits))

	for _, t := range toolkits {
		if pbToolkit := ToolkitToPBToolkit(t); pbToolkit != nil {
			pbToolkits = append(pbToolkits, pbToolkit)
		}
	}

	return &pbinventory.GetAllToolkitsResponse{
		Toolkits: pbToolkits,
	}
}
