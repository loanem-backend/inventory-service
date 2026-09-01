package mapper

import (
	"time"

	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

func AddLoanRequestToLoan(req *pbinventory.AddLoanRequest) *entity.Loan {
	if req == nil {
		return nil
	}

	return &entity.Loan{
		Toolkit: entity.Toolkit{
			ID: int(req.GetToolkitId()),
		},
		Team: entity.Team{
			ID: req.GetTeamId(),
		},
		Submitter: entity.Participant{
			ID: req.GetSubmitterId(),
		},
		Date: time.Date(
			int(req.GetDate().GetYear()),
			time.Month(req.GetDate().GetMonth()),
			int(req.GetDate().GetDay()),
			0, 0, 0, 0, time.UTC,
		),
		SessionNumber: int(req.GetSessionNumber()),
	}
}

func StringToAddLoanResponse(id string) *pbinventory.AddLoanResponse {
	return &pbinventory.AddLoanResponse{
		Id: id,
	}
}
