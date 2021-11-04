package worker

import (
	"time"

	"github.com/acentswap/bridge/rpc/client"
	"github.com/acentswap/bridge/tokens/bridge"
)

const interval = 10 * time.Millisecond

// StartWork start swap server work
func StartWork(isServer bool) {
	if isServer {
		logWorker("worker", "start server worker")
	} else {
		logWorker("worker", "start oracle worker")
	}

	client.InitHTTPClient()
	bridge.InitCrossChainBridge(isServer)

	StartScanJob(isServer)
	time.Sleep(interval)

	StartUpdateLatestBlockHeightJob()
	time.Sleep(interval)

	if !isServer {
		StartAcceptSignJob()
		time.Sleep(interval)
		AddTokenPairDynamically()
		return
	}

	StartSwapJob()
	time.Sleep(interval)

	StartVerifyJob()
	time.Sleep(interval)

	StartStableJob()
	time.Sleep(interval)

	StartReplaceJob()
	time.Sleep(interval)

	StartPassBigValueJob()
	time.Sleep(interval)

	StartAggregateJob()
}
