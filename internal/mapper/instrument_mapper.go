package mapper

import (
	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func IntToAddInstrumentResponse(id int32) *pbinventory.AddInstrumentResponse {
	return &pbinventory.AddInstrumentResponse{
		Id: id,
	}
}

func InstrumentToPBInstrument(i *entity.Instrument) *pbinventory.Instrument {
	if i == nil {
		return nil
	}

	return &pbinventory.Instrument{
		Id:        int32(i.ID),
		Name:      i.Name,
		Picture:   i.Picture,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
	}
}

func InstrumentsToGetAllInstrumentsResponse(instruments []*entity.Instrument) *pbinventory.GetAllInstrumentsResponse {
	pbInstruments := make([]*pbinventory.Instrument, 0, len(instruments))

	for _, i := range instruments {
		if pbInstrument := InstrumentToPBInstrument(i); pbInstrument != nil {
			pbInstruments = append(pbInstruments, pbInstrument)
		}
	}

	return &pbinventory.GetAllInstrumentsResponse{
		Instruments: pbInstruments,
	}
}

func SetInstrumentPictureRequestToInstrument(req *pbinventory.SetInstrumentPictureRequest) *entity.Instrument {
	return &entity.Instrument{
		ID:      int(req.GetId()),
		Picture: req.GetKey(),
	}
}
