package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/common"

	"github.com/xudaquan2003/storage-slot/withdraw"
)

type network struct {
	l2RPC         string
	portalAddress string
	l2OOAddress   string
}

var networks = map[string]network{
	"mega-testnet": {
		l2RPC:         "http://127.0.0.1:19545",
		portalAddress: "0x978e3286EB805934215a88694d80b09aDed68D90",
		l2OOAddress:   "0xD31598c909d9C935a9e35bA70d9a3DD47d4D5865",
	},
}

func main() {
	var networkKeys []string
	for n := range networks {
		networkKeys = append(networkKeys, n)
	}

	// var rpcFlag string
	var networkFlag string
	var l2RpcFlag string
	var withdrawalFlag string
	flag.StringVar(&networkFlag, "network", "mega-testnet", fmt.Sprintf("op-stack network to withdraw.go from (one of: %s)", strings.Join(networkKeys, ", ")))
	flag.StringVar(&l2RpcFlag, "l2-rpc", "", "Custom network L2 RPC url")
	flag.StringVar(&withdrawalFlag, "withdrawal", "", "TX hash of the L2 withdrawal transaction")
	flag.Parse()

	log.SetDefault(oplog.NewLogger(os.Stderr, oplog.DefaultCLIConfig()))

	n, ok := networks[networkFlag]
	if !ok {
		log.Crit("Unknown network", "network", networkFlag)
	}

	if l2RpcFlag != "" {
		if l2RpcFlag == "" {
			log.Crit("Missing --l2-rpc flag")
		}

		n = network{
			l2RPC: l2RpcFlag,
		}
	}

	if withdrawalFlag == "" {
		log.Crit("Missing --withdrawal flag")
	}
	withdrawal := common.HexToHash(withdrawalFlag)

	ctx := context.Background()

	l2Client, err := rpc.DialContext(ctx, n.l2RPC)
	if err != nil {
		log.Crit("Error dialing L2 client", "error", err)
	}

	slot, err := withdraw.TxSlot(ctx, l2Client, withdrawal)
	if err != nil {
		log.Crit("Error querying withdrawal proof", "error", err)
	}

	fmt.Println(slot)
}
