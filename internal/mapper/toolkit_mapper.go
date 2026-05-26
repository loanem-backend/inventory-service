package mapper

import (
	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

func AddToolkitRequestToToolkit(req *pbinventory.AddToolkitRequest) *entity.Toolkit {
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
