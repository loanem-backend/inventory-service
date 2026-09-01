package mapper

import (
	"time"

	"github.com/loanem-backend/inventory-service/internal/entity"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func GetLoansByDateRequestToTimeDate(req *pbinventory.GetLoansByDateRequest) *time.Time {
	if req == nil {
		return nil
	}

	t := time.Date(
		int(req.GetDate().GetYear()),
		time.Month(req.GetDate().GetMonth()),
		int(req.GetDate().GetDay()),
		0, 0, 0, 0, time.UTC,
	)

	return &t
}

func LoanToPBLoan(loan *entity.Loan) *pbinventory.Loan {
	if loan == nil {
		return nil
	}

	return &pbinventory.Loan{
		Id: loan.ID.String(),
		Toolkit: &pbinventory.Toolkit{
			Id:       int32(loan.Toolkit.ID),
			KitName:  loan.Toolkit.KitName,
			CourseId: int32(loan.Toolkit.Course.ID),
		},
		TeamNumber:    int32(loan.Team.Number),
		ClassId:       int32(loan.Team.ClassID),
		ClassName:     loan.Team.ClassName,
		SubmitterNim:  loan.Submitter.Nim,
		SubmitterName: loan.Submitter.Name,
		Date: &date.Date{
			Year:  int32(loan.Date.Year()),
			Month: int32(loan.Date.Month()),
			Day:   int32(loan.Date.Day()),
		},
		SessionNumber: int32(loan.SessionNumber),
		Status:        string(loan.Status),
		Note:          loan.Note,
		CreatedAt:     timestamppb.New(loan.CreatedAt),
		UpdatedAt:     timestamppb.New(loan.UpdatedAt),
	}
}

func LoansToGetLoansByDateResponse(loans []*entity.Loan) *pbinventory.GetLoansByDateResponse {
	pbLoans := make([]*pbinventory.Loan, 0, len(loans))

	for _, l := range loans {
		if pbLoan := LoanToPBLoan(l); pbLoan != nil {
			pbLoans = append(pbLoans, pbLoan)
		}
	}

	return &pbinventory.GetLoansByDateResponse{
		Loans: pbLoans,
	}
}
