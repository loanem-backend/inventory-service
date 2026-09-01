package server

import (
	"context"

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
	idData, err := s.loanServ.CreateLoan(ctx, mapper.AddLoanRequestToLoan(req))
	if err != nil {
		return nil, err
	}

	return mapper.StringToAddLoanResponse(idData), nil
}
