package rpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
)

func StartServer() {
	http.HandleFunc("/", handleJSONRPC)
	http.ListenAndServe(":5279", nil)
}

func sendResultResponse(w http.ResponseWriter, result any) {
	json.NewEncoder(w).Encode(map[string]any{
		"jsonrpc": "2.0",
		"result":  result,
	})
}

func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	json.NewEncoder(w).Encode(map[string]any{
		"jsonrpc": "2.0",
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

func handleJSONRPC(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if strings.EqualFold(req.Method, "POST") {
		var messageBatch []map[string]any
		var message map[string]any

		body, _ := ioutil.ReadAll(req.Body)

		errBatch := json.Unmarshal(body, &messageBatch)
		err := json.Unmarshal(body, &message)

		if errBatch != nil && err != nil {
			w.WriteHeader(http.StatusBadRequest)
			sendErrorResponse(w, -32700, "Cannot parse invalid JSON data.")
			return
		}
		if messageBatch != nil {
			w.WriteHeader(http.StatusBadRequest)
			sendErrorResponse(w, -32700, "Batches are not supported")
			return
		}

		handleJSONRPCMessage(w, message)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	sendErrorResponse(w, -32700, "HTTP method not allowed.")
}

var handlers = map[string]func(http.ResponseWriter, any){
	"account_add":             handleAccountAdd,
	"account_balance":         handleAccountBalance,
	"account_create":          handleAccountCreate,
	"account_deposit":         handleAccountDeposit,
	"account_fund":            handleAccountFund,
	"account_list":            handleAccountList,
	"account_max_address_gap": handleAccountMaxAddressGap,
	"account_remove":          handleAccountRemove,
	"account_send":            handleAccountSend,
	"account_set":             handleAccountSet,
	"address_is_mine":         handleAddressIsMine,
	"address_list":            handleAddressList,
	"address_unused":          handleAddressUnused,
	"blob_announce":           handleBlobAnnounce,
	"blob_clean":              handleBlobClean,
	"blob_delete":             handleBlobDelete,
	"blob_get":                handleBlobGet,
	"blob_list":               handleBlobList,
	"blob_reflect":            handleBlobReflect,
	"blob_reflect_all":        handleBlobReflectAll,
	"channel_abandon":         handleChannelAbandon,
	"channel_create":          handleChannelCreate,
	"channel_list":            handleChannelList,
	"channel_sign":            handleChannelSign,
	"channel_update":          handleChannelUpdate,
	"claim_list":              handleClaimList,
	"claim_search":            handleClaimSearch,
	"collection_abandon":      handleCollectionAbandon,
	"collection_create":       handleCollectionCreate,
	"collection_list":         handleCollectionList,
	"collection_resolve":      handleCollectionResolve,
	"collection_update":       handleCollectionUpdate,
	"ffmpeg_find":             handleFfmpegFind,
	"file_delete":             handleFileDelete,
	"file_list":               handleFileList,
	"file_reflect":            handleFileReflect,
	"file_save":               handleFileSave,
	"file_set_status":         handleFileSetStatus,
	"get":                     handleGet,
	"peer_list":               handlePeerList,
	"peer_ping":               handlePeerPing,
	"preference_get":          handlePreferenceGet,
	"preference_set":          handlePreferenceSet,
	"publish":                 handlePublish,
	"purchase_create":         handlePurchaseCreate,
	"purchase_list":           handlePurchaseList,
	"resolve":                 handleResolve,
	"routing_table_get":       handleRoutingTableGet,
	"settings_clear":          handleSettingsClear,
	"settings_get":            handleSettingsGet,
	"settings_set":            handleSettingsSet,
	"status":                  handleStatus,
	"stop":                    handleStop,
	"stream_abandon":          handleStreamAbandon,
	"stream_cost_estimate":    handleStreamCostEstimate,
	"stream_create":           handleStreamCreate,
	"stream_list":             handleStreamList,
	"stream_repost":           handleStreamRepost,
	"stream_update":           handleStreamUpdate,
	"support_abandon":         handleSupportAbandon,
	"support_create":          handleSupportCreate,
	"support_list":            handleSupportList,
	"support_sum":             handleSupportSum,
	"sync_apply":              handleSyncApply,
	"sync_hash":               handleSyncHash,
	"tracemalloc_disable":     handleTracemallocDisable,
	"tracemalloc_enable":      handleTracemallocEnable,
	"tracemalloc_top":         handleTracemallocTop,
	"transaction_list":        handleTransactionList,
	"transaction_show":        handleTransactionShow,
	"txo_list":                handleTxoList,
	"txo_plot":                handleTxoPlot,
	"txo_spend":               handleTxoSpend,
	"txo_sum":                 handleTxoSum,
	"utxo_list":               handleUtxoList,
	"utxo_release":            handleUtxoRelease,
	"version":                 handleVersion,
	"wallet_add":              handleWalletAdd,
	"wallet_balance":          handleWalletBalance,
	"wallet_create":           handleWalletCreate,
	"wallet_decrypt":          handleWalletDecrypt,
	"wallet_encrypt":          handleWalletEncrypt,
	"wallet_export":           handleWalletExport,
	"wallet_import":           handleWalletImport,
	"wallet_list":             handleWalletList,
	"wallet_lock":             handleWalletLock,
	"wallet_reconnect":        handleWalletReconnect,
	"wallet_remove":           handleWalletRemove,
	"wallet_send":             handleWalletSend,
	"wallet_status":           handleWalletStatus,
	"wallet_unlock":           handleWalletUnlock,
}

func handleJSONRPCMessage(w http.ResponseWriter, message map[string]any) {
	method, existsMethod := message["method"].(string)
	params, existsParams := message["params"]

	if !existsMethod {
		w.WriteHeader(http.StatusBadRequest)
		sendErrorResponse(w, -32600, "Method property is missing.")
		return
	}

	handler, exists := handlers[method]
	if exists {
		if existsParams {
			handler(w, params)
			return
		}
		handler(w, nil)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	sendErrorResponse(w, -32601, "Unknown JSON-RPC method.")
}

func handleStatus(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"jsonrpc": "2.0",
		"result":  map[string]any{},
	})
}

func handleStop(w http.ResponseWriter, params any) {
	sendResultResponse(w, "Shutting down")
}

func handleResolve(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "resolve not implemented"})
}

func handleGet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "get not implemented"})
}

func handlePublish(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "publish not implemented"})
}

func handleAccountList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleAccountCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "account_create not implemented"})
}

func handleAccountAdd(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "account_add not implemented"})
}

func handleAccountBalance(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"available": "0.0",
		"reserved":  "0.0",
		"reserved_subtotals": map[string]any{
			"claims":   "0.0",
			"supports": "0.0",
			"tips":     "0.0",
		},
		"total": "0.0",
	})
}

func handleAccountSend(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "account_send not implemented"})
}

func handleAccountFund(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "account_fund not implemented"})
}

func handleAccountDeposit(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "account_deposit not implemented"})
}

func handleAccountMaxAddressGap(w http.ResponseWriter, params any) {
	sendResultResponse(w, 1)
}

func handleAccountRemove(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleAccountSet(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleAddressList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleAddressIsMine(w http.ResponseWriter, params any) {
	sendResultResponse(w, false)
}

func handleAddressUnused(w http.ResponseWriter, params any) {
	sendResultResponse(w, "")
}

func handleChannelCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "channel_create not implemented"})
}

func handleChannelUpdate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "channel_update not implemented"})
}

func handleChannelList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleChannelAbandon(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "channel_abandon not implemented"})
}

func handleChannelSign(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "channel_sign not implemented"})
}

func handleClaimList(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"page":        1,
		"page_size":   20,
		"total_pages": 0,
		"total_items": 0,
		"items":       []any{},
	})
}

func handleClaimSearch(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"page":        1,
		"page_size":   20,
		"total_pages": 0,
		"total_items": 0,
		"items":       []any{},
	})
}

func handleCollectionList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleCollectionCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "collection_create not implemented"})
}

func handleCollectionUpdate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "collection_update not implemented"})
}

func handleCollectionAbandon(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "collection_abandon not implemented"})
}

func handleCollectionResolve(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "collection_resolve not implemented"})
}

func handleFileList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleFileDelete(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "file_delete not implemented"})
}

func handleFileReflect(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "file_reflect not implemented"})
}

func handleFileSave(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleFileSetStatus(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleBlobList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleBlobGet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "blob_get not implemented"})
}

func handleBlobAnnounce(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleBlobDelete(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleBlobReflect(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleBlobReflectAll(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleBlobClean(w http.ResponseWriter, params any) {
	sendResultResponse(w, 0)
}

func handleStreamCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "stream_create not implemented"})
}

func handleStreamUpdate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "stream_update not implemented"})
}

func handleStreamList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleStreamAbandon(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "stream_abandon not implemented"})
}

func handleStreamRepost(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "stream_repost not implemented"})
}

func handleStreamCostEstimate(w http.ResponseWriter, params any) {
	sendResultResponse(w, 0.0)
}

func handleSupportCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "support_create not implemented"})
}

func handleSupportAbandon(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "support_abandon not implemented"})
}

func handleSupportList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleSupportSum(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"support_amount": "0.0"})
}

func handleTransactionList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleTransactionShow(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "transaction_show not implemented"})
}

func handleUtxoList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleUtxoRelease(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleTxoList(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"page":        1,
		"page_size":   20,
		"total_pages": 0,
		"total_items": 0,
		"items":       []any{},
	})
}

func handleTxoSum(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"amount":   "0.0",
		"currency": "LBC",
	})
}

func handleTxoSpend(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "txo_spend not implemented"})
}

func handleTxoPlot(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "txo_plot not implemented"})
}

func handleWalletList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleWalletCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "wallet_create not implemented"})
}

func handleWalletAdd(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "wallet_add not implemented"})
}

func handleWalletBalance(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"available": "0.0",
		"reserved":  "0.0",
		"reserved_subtotals": map[string]any{
			"claims":   "0.0",
			"supports": "0.0",
			"tips":     "0.0",
		},
		"total": "0.0",
	})
}

func handleWalletSend(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "wallet_send not implemented"})
}

func handleWalletEncrypt(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletDecrypt(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletLock(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletUnlock(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletStatus(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"connected":         false,
		"encrypted":         false,
		"locked":            true,
		"items":             0,
		"keys":              0,
		"height":            0,
		"wallet_db_version": 0,
	})
}

func handleWalletExport(w http.ResponseWriter, params any) {
	sendResultResponse(w, "")
}

func handleWalletImport(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletRemove(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleWalletReconnect(w http.ResponseWriter, params any) {
	sendResultResponse(w, true)
}

func handleRoutingTableGet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "routing_table_get not implemented"})
}

func handlePeerList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handlePeerPing(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "peer_ping not implemented"})
}

func handleSyncApply(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "sync_apply not implemented"})
}

func handleSyncHash(w http.ResponseWriter, params any) {
	sendResultResponse(w, "")
}

func handleSettingsGet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"share_usage_data":     true,
		"download_timeout":     60,
		"max_key_fee":          map[string]any{},
		"search_timeout":       10,
		"max_search_results":   25,
		"filter_bluray":        false,
		"currency":             "USD",
		"exchange_rate_fiat":   "1.0",
		"run_reflector_server": false,
		"audio_bitrate":        128,
		"video_bitrate":        1000,
		"audio_normalize":      false,
		"data_rate":            0.5,
		"max_upload_rate":      0,
		"max_download_rate":    0,
		"upload_log":           true,
		"upload_log_filename":  "",
		"peer_port":            3333,
		"dht_node_port":        4444,
		"downloader_limit":     32,
		"uploader_limit":       64,
		"ui_theme":             "light",
		"language":             "en",
		"background_va":        false,
		"ffmpeg_path":          "",
		"bid_influence":        1.0,
		"cache_time":           60,
		"romfs_root":           "",
		"romfs_level":          0,
		"prometheus_port":      0,
		"allowed_origin":       "",
		"streaming_get":        true,
		"streaming_port":       0,
		"streaming_hosts":      []any{},
	})
}

func handleSettingsSet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "settings_set not implemented"})
}

func handleSettingsClear(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "settings_clear not implemented"})
}

func handlePreferenceGet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{})
}

func handlePreferenceSet(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "preference_set not implemented"})
}

func handlePurchaseList(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handlePurchaseCreate(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"error": "purchase_create not implemented"})
}

func handleFfmpegFind(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{
		"available":            false,
		"which":                "",
		"analyze_audio_volume": false,
	})
}

func handleTracemallocDisable(w http.ResponseWriter, params any) {
	sendResultResponse(w, nil)
}

func handleTracemallocEnable(w http.ResponseWriter, params any) {
	sendResultResponse(w, map[string]any{"tracing": true})
}

func handleTracemallocTop(w http.ResponseWriter, params any) {
	sendResultResponse(w, []any{})
}

func handleVersion(w http.ResponseWriter, params any) {
	info, _ := debug.ReadBuildInfo()

	sendResultResponse(w, map[string]any{
		"jsonrpc": "2.0",
		"result": map[string]any{
			"build":           nil,
			"lbrynet_version": nil,
			"os_release":      nil,
			"os_system":       nil,
			"platform":        nil,
			"processor":       nil,
			"python_version":  nil,
			"version":         info.Main.Version,
		},
	})
}
