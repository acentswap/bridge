// Package bridge init crosschain bridges.
package bridge

import (
	"fmt"
	"strings"

	"github.com/acentswap/bridge/dcrm"
	"github.com/acentswap/bridge/log"
	"github.com/acentswap/bridge/params"
	"github.com/acentswap/bridge/tokens"
	"github.com/acentswap/bridge/tokens/block"
	"github.com/acentswap/bridge/tokens/btc"
	"github.com/acentswap/bridge/tokens/colx"
	"github.com/acentswap/bridge/tokens/etc"
	"github.com/acentswap/bridge/tokens/eth"
	"github.com/acentswap/bridge/tokens/fsn"
	"github.com/acentswap/bridge/tokens/ltc"
	"github.com/acentswap/bridge/tokens/tools"
	"github.com/acentswap/bridge/tokens/okex"
)

// NewCrossChainBridge new bridge according to chain name
func NewCrossChainBridge(id string, isSrc bool) tokens.CrossChainBridge {
	blockChainIden := strings.ToUpper(id)
	switch {
	case strings.HasPrefix(blockChainIden, "BITCOIN"):
		return btc.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "LITECOIN"):
		return ltc.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "BLOCK"):
		return block.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "ETHCLASSIC"):
		return etc.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "ETHEREUM"):
		return eth.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "OKEX"):
		return okex.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "FUSION"):
		return fsn.NewCrossChainBridge(isSrc)
	case strings.HasPrefix(blockChainIden, "COLOSSUS") || strings.HasPrefix(blockChainIden, "COLX"):
		return colx.NewCrossChainBridge(isSrc)
	default:
		log.Fatalf("Unsupported block chain %v", id)
		return nil
	}
}

// InitCrossChainBridge init bridge
func InitCrossChainBridge(isServer bool) {
	cfg := params.GetConfig()
	srcChain := cfg.SrcChain
	dstChain := cfg.DestChain
	srcGateway := cfg.SrcGateway
	dstGateway := cfg.DestGateway

	srcID := srcChain.BlockChain
	dstID := dstChain.BlockChain
	srcNet := srcChain.NetID
	dstNet := dstChain.NetID

	tokens.AggregateIdentifier = fmt.Sprintf("%s:%s", params.GetIdentifier(), tokens.AggregateIdentifier)

	tokens.SrcBridge = NewCrossChainBridge(srcID, true)
	tokens.DstBridge = NewCrossChainBridge(dstID, false)
	log.Info("New bridge finished", "source", srcID, "sourceNet", srcNet, "dest", dstID, "destNet", dstNet)

	tokens.SrcBridge.SetChainAndGateway(srcChain, srcGateway)
	log.Info("Init bridge source", "source", srcID, "gateway", srcGateway)

	tokens.DstBridge.SetChainAndGateway(dstChain, dstGateway)
	log.Info("Init bridge destation", "dest", dstID, "gateway", dstGateway)

	tokens.SrcNonceSetter, _ = tokens.SrcBridge.(tokens.NonceSetter)
	tokens.DstNonceSetter, _ = tokens.DstBridge.(tokens.NonceSetter)

	tokens.SrcForkChecker, _ = tokens.SrcBridge.(tokens.ForkChecker)
	tokens.DstForkChecker, _ = tokens.DstBridge.(tokens.ForkChecker)

	tokens.SrcStableConfirmations = *tokens.SrcBridge.GetChainConfig().Confirmations
	tokens.DstStableConfirmations = *tokens.DstBridge.GetChainConfig().Confirmations

	tools.AdjustGatewayOrder(true)
	tools.AdjustGatewayOrder(false)

	tokens.IsDcrmDisabled = cfg.Dcrm.Disable
	tokens.LoadTokenPairsConfig(true)

	BlockChain := strings.ToUpper(srcChain.BlockChain)
	switch BlockChain {
	case "BITCOIN":
		btc.Init(cfg.BtcExtra)
	case "LITECOIN":
		ltc.Init(cfg.BtcExtra)
	case "BLOCK":
		block.Init(cfg.BtcExtra)
	case "COLX":
		colx.Init(cfg.BtcExtra)
	default:
		cfg.BtcExtra = nil
	}

	dcrm.Init(cfg.Dcrm, isServer)

	log.Info("Init bridge success", "isServer", isServer, "dcrmEnabled", !cfg.Dcrm.Disable)
}
