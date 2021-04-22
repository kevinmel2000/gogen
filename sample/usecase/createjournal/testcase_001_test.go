package createjournal

import (
	"accounting/gateway"
	"accounting/infrastructure/log"
	"context"
	"testing"
)

func TestCase001(t *testing.T) {

	ctx := log.ContextWithLogGroupID(context.Background())

	usecase := NewUsecase(gateway.NewInmemoryGateway())
	res, err := usecase.Execute(ctx, InportRequest{
		BusinessID:  "",
		Date:        "",
		Description: "",
		JournalType: "",
		UserID:      "",
		JournalBalances: []JournalBalanceRequest{
			{AccountCode: "100", Nominal: 200},
			{AccountCode: "101", Nominal: 200},
		},
	})

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(res)

}
