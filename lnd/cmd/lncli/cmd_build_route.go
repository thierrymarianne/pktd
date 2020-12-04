package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pkt-cash/pktd/lnd/chainreg"
	"github.com/pkt-cash/pktd/lnd/lnrpc/routerrpc"
	"github.com/pkt-cash/pktd/lnd/routing/route"
	"github.com/urfave/cli"
)

var buildRouteCommand = cli.Command{
	Name:     "buildroute",
	Category: "Payments",
	Usage:    "Build a route from a list of hop pubkeys.",
	Action:   actionDecorator(buildRoute),
	Flags: []cli.Flag{
		cli.Int64Flag{
			Name: "amt",
			Usage: "the amount to send expressed in satoshis. If" +
				"not set, the minimum routable amount is used",
		},
		cli.Int64Flag{
			Name: "final_cltv_delta",
			Usage: "number of blocks the last hop has to reveal " +
				"the preimage",
			Value: chainreg.DefaultBitcoinTimeLockDelta,
		},
		cli.StringFlag{
			Name:  "hops",
			Usage: "comma separated hex pubkeys",
		},
		cli.Uint64Flag{
			Name: "outgoing_chan_id",
			Usage: "short channel id of the outgoing channel to " +
				"use for the first hop of the payment",
			Value: 0,
		},
	},
}

func buildRoute(ctx *cli.Context) er.R {
	conn := getClientConn(ctx, false)
	defer conn.Close()

	client := routerrpc.NewRouterClient(conn)

	if !ctx.IsSet("hops") {
		return er.New("hops required")
	}

	// Build list of hop addresses for the rpc.
	hops := strings.Split(ctx.String("hops"), ",")
	rpcHops := make([][]byte, 0, len(hops))
	for _, k := range hops {
		pubkey, err := route.NewVertexFromStr(k)
		if err != nil {
			return er.Errorf("error parsing %v: %v", k, err)
		}
		rpcHops = append(rpcHops, pubkey[:])
	}

	var amtMsat int64
	hasAmt := ctx.IsSet("amt")
	if hasAmt {
		amtMsat = ctx.Int64("amt") * 1000
		if amtMsat == 0 {
			return er.Errorf("non-zero amount required")
		}
	}

	// Call BuildRoute rpc.
	req := &routerrpc.BuildRouteRequest{
		AmtMsat:        amtMsat,
		FinalCltvDelta: int32(ctx.Int64("final_cltv_delta")),
		HopPubkeys:     rpcHops,
		OutgoingChanId: ctx.Uint64("outgoing_chan_id"),
	}

	rpcCtx := context.Background()
	route, err := client.BuildRoute(rpcCtx, req)
	if err != nil {
		return err
	}

	printRespJSON(route)

	return nil
}
