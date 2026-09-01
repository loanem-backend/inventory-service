package server

import (
	"context"
	"errors"

	"github.com/loanem-backend/inventory-service/internal/mapper"
	"github.com/loanem-backend/inventory-service/internal/service"
	pbinventory "github.com/loanem-backend/protos/pb/proto/services/inventory/v1"
)

type LoanServer struct {
	pbinventory.UnimplementedLoanServiceServer
	loanServ service.LoanService
}

func NewLoanServer(ls service.LoanService) *LoanServer {
	return &LoanServer{
		loanServ: ls,
	}
}

func (s *LoanServer) AddLoan(ctx context.Context, req *pbinventory.AddLoanRequest) (*pbinventory.AddLoanResponse, error) {
	idData, err := s.loanServ.Create(ctx, mapper.AddLoanRequestToLoan(req))
	if err != nil {
		return nil, err
	}

	return mapper.StringToAddLoanResponse(idData), nil
}

func (s *LoanServer) GetLoansByDate(ctx context.Context, req *pbinventory.GetLoansByDateRequest) (*pbinventory.GetLoansByDateResponse, error) {
	reqDate := mapper.GetLoansByDateRequestToTimeDate(req)
	if reqDate == nil {
		return nil, errors.New("invalid argument")
	}

	loansData, err := s.loanServ.GetByDate(ctx, *reqDate)
	if err != nil {
		return nil, err
	}

	return mapper.LoansToGetLoansByDateResponse(loansData), nil
}
