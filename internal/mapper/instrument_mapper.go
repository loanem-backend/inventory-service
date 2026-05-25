package mapper

import pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"

func IntToAddInstrumentResponse(id int32) *pbinventory.AddInstrumentResponse {
	return &pbinventory.AddInstrumentResponse{
		Id: id,
	}
}
