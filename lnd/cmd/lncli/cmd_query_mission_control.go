package main

import (
	"context"

	"github.com/pkt-cash/pktd/lnd/lnrpc/routerrpc"

	"github.com/urfave/cli"
)

var queryMissionControlCommand = cli.Command{
	Name:     "querymc",
	Category: "Payments",
	Usage:    "Query the internal mission control state.",
	Action:   actionDecorator(queryMissionControl),
}

func queryMissionControl(ctx *cli.Context) er.R {
	conn := getClientConn(ctx, false)
	defer conn.Close()

	client := routerrpc.NewRouterClient(conn)

	req := &routerrpc.QueryMissionControlRequest{}
	rpcCtx := context.Background()
	snapshot, err := client.QueryMissionControl(rpcCtx, req)
	if err != nil {
		return err
	}

	printRespJSON(snapshot)

	return nil
}
