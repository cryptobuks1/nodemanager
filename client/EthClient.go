package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/istanbul/backend"
	"time"

	"github.com/ybbus/jsonrpc"
	"github.com/nordicenergy/powerchain-maker-nodemanager/contracthandler"
	"github.com/ethereum/go-ethereum/common"
)

type AdminInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Enode      string    `json:"enode"`
	IP         string    `json:"ip"`
	Ports      Ports     `json:"ports"`
	ListenAddr string    `json:"listenAddr"`
	Protocols  Protocols `json:"protocols"`
}

type Ports struct {
	Discovery int `json:"discovery"`
	Listener  int `json:"listener"`
}

type AdminPeers struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Caps      []string  `json:"caps"`
	Network   Network   `json:"network"`
	Protocols Protocols `json:"protocols"`
}

type Protocols struct {
	Eth Eth `json:"eth"`
}

type Eth struct {
	Network    int    `json:"network"`
	Version    int    `json:"version"`
	Difficulty int    `json:"difficulty"`
	Genesis    string `json:"genesis"`
	Head       string `json:"head"`
}

type Network struct {
	LocalAddress  string `json:"localAddress"`
	RemoteAddress string `json:"remoteAddress"`
}

type BlockDetailsResponse struct {
	Number           string                       `json:"number"`
	Hash             string                       `json:"hash"`
	ParentHash       string                       `json:"parentHash"`
	Nonce            string                       `json:"nonce"`
	Sha3Uncles       string                       `json:"sha3Uncles"`
	LogsBloom        string                       `json:"logsBloom"`
	TransactionsRoot string                       `json:"transactionsRoot"`
	StateRoot        string                       `json:"stateRoot"`
	Miner            string                       `json:"miner"`
	Difficulty       string                       `json:"difficulty"`
	TotalDifficulty  string                       `json:"totalDifficulty"`
	ExtraData        string                       `json:"extraData"`
	Size             string                       `json:"size"`
	GasLimit         string                       `json:"gasLimit"`
	GasUsed          string                       `json:"gasUsed"`
	Timestamp        string                       `json:"timestamp"`
	Transactions     []TransactionDetailsResponse `json:"transactions"`
	Uncles           []string                     `json:"uncles"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type TransactionReceiptResponse struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	GasUsed           string `json:"gasUsed"`
	Logs              []Logs `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Root              string `json:"root"`
	To                string `json:"to"`
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
}

type Logs struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type Payload struct {
	From       string   `json:"from"`
	To         string   `json:"to,omitempty"`
	Data       string   `json:"data"`
	Gaslimit   string   `json:"gas"`
	PrivateFor []string `json:"privateFor,omitempty"`
}

type CallPayload struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type EthClient struct {
	Url string
}

// AccountStats
// in case it refers to validators, value means number of mined blocks
// in case it refers to users, value means tootal users's gas consumption
type AccountStats struct {
	Account common.Address
	Value   uint64
}

type IstanbulStats struct {
	Validators  []common.Address `json:"validators"`
	BlocksMined []uint32         `json:"blocks_mined"`

	Users           []common.Address `json:"users"`
	GasConsumptions []uint64         `json:"gas_consumptions"`

	MaxGas uint64 `json:"max_gas_used"`
}

func (ec *EthClient) GetTransactionByHash(txNo string) TransactionDetailsResponse {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getTransactionByHash", txNo)

	if err != nil {
		fmt.Println(err)
	}
	txResponse := TransactionDetailsResponse{}

	if response == nil {
		fmt.Println("No response returned!")
		return txResponse
	}

	err = response.GetObject(&txResponse)
	if err != nil {
		fmt.Println(err)
	}
	return txResponse
}

func (ec *EthClient) ProposeValidator(address string, vote bool) error {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("istanbul_propose", address, vote)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if response == nil {
		fmt.Println("No response returned!")
		return nil
	}

	return nil
}

func (ec *EthClient) AddPeer(enode string) error {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("admin_addPeer", enode)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if response == nil {
		fmt.Println("No response returned!")
		return nil
	}

	return nil
}

func (ec *EthClient) GetCandidates() map[common.Address]bool {
	rpcClient := jsonrpc.NewClient(ec.Url)

	txResponse := make(map[common.Address]bool)

	response, err := rpcClient.Call("istanbul_getCandidates")

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	err = response.GetObject(&txResponse)

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	return txResponse
}

func (ec *EthClient) GetValidators(blockNumber string) []common.Address {
	rpcClient := jsonrpc.NewClient(ec.Url)

	txResponse := []common.Address{}
	response, err := rpcClient.Call("istanbul_getValidators", blockNumber)

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	if response == nil {
		fmt.Println("No response returned!")
		return txResponse
	}

	err = response.GetObject(&txResponse)

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	return txResponse
}

func (ec *EthClient) GetSnapshot(blockNumber string) backend.Snapshot {
	rpcClient := jsonrpc.NewClient(ec.Url)

	txResponse := backend.Snapshot{}
	response, err := rpcClient.Call("istanbul_getSnapshot", blockNumber)

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	if response == nil {
		fmt.Println("No response returned!")
		return txResponse
	}

	err = response.GetObject(&txResponse)

	if err != nil {
		fmt.Println(err)
		return txResponse
	}

	return txResponse
}

func (ec *EthClient) GetStatistics(start string, end string) IstanbulStats {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("istanbul_getStatistics", start, end)

	if err != nil {
		fmt.Println(err)
		return IstanbulStats{}
	}

	if response == nil {
		fmt.Println("No response returned!")
		return IstanbulStats{}
	}

	txResponse := IstanbulStats{}

	err = response.GetObject(&txResponse)
	if err != nil {
		fmt.Println(err)
	}
	return txResponse
}

func (ec *EthClient) GetBlockByNumber(blockNo string) BlockDetailsResponse {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getBlockByNumber", blockNo, true)

	if err != nil {
		fmt.Println(err)
		return BlockDetailsResponse{}
	}

	if response == nil {
		fmt.Println("No response returned!")
		return BlockDetailsResponse{}
	}

	blockResponse := BlockDetailsResponse{}
	err = response.GetObject(&blockResponse)
	if err != nil {
		fmt.Println(err)
	}
	return blockResponse
}

func (ec *EthClient) PendingTransactions() []TransactionDetailsResponse {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_pendingTransactions")
	if err != nil {
		fmt.Println(err)
	}

	pendingTxResponse := []TransactionDetailsResponse{}

	if response == nil {
		fmt.Println("No response returned!")
		return pendingTxResponse
	}

	err = response.GetObject(&pendingTxResponse)
	if err != nil {
		fmt.Println(err)
	}
	return pendingTxResponse
}

func (ec *EthClient) AdminPeers() []AdminPeers {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("admin_peers")
	if err != nil {
		fmt.Println(err)
	}
	otherPeersResponse := []AdminPeers{}

	if response == nil {
		fmt.Println("No response returned!")
		return otherPeersResponse
	}

	err = response.GetObject(&otherPeersResponse)
	if err != nil {
		fmt.Println(err)
	}
	return otherPeersResponse
}

func (ec *EthClient) AdminNodeInfo() AdminInfo {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("admin_nodeInfo")
	if err != nil {
		fmt.Println(err)
	}

	thisAdminInfo := AdminInfo{}

	if response == nil {
		fmt.Println("No response returned!")
		return thisAdminInfo
	}

	err = response.GetObject(&thisAdminInfo)
	return thisAdminInfo
}

func (ec *EthClient) BlockNumber() string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_blockNumber")
	if err != nil {
		fmt.Println(err)
	}

	var blockNumber string

	if response == nil {
		fmt.Println("No response returned!")
		return blockNumber
	}

	if err == nil {
		err = response.GetObject(&blockNumber)
	}
	if err != nil {
		fmt.Println(err)
	}
	return blockNumber
}

func (ec *EthClient) Coinbase() string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_coinbase")
	if err != nil {
		fmt.Println(err)
	}
	var coinbase string

	if response == nil {
		fmt.Println("No response returned!")
		return coinbase
	}

	if err == nil {
		err = response.GetObject(&coinbase)
	}
	if err != nil {
		fmt.Println(err)
	}
	return coinbase
}

func (ec *EthClient) GetTransactionReceipt(txNo string) TransactionReceiptResponse {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getTransactionReceipt", txNo)

	if err != nil {
		fmt.Println(err)
	}

	txResponse := TransactionReceiptResponse{}

	if response == nil {
		fmt.Println("No response returned!")
		return txResponse
	}

	if response == nil {
		fmt.Println("No response returned!")
		return txResponse
	}

	err = response.GetObject(&txResponse)
	if err != nil {
		fmt.Println(err)
	}
	return txResponse
}

func (ec *EthClient) SendTransaction(param contracthandler.ContractParam, rh contracthandler.RequestHandler) string {
	rpcClient := jsonrpc.NewClient(ec.Url)

	fmt.Println("Entering SendTransaction")

	response, err := rpcClient.Call("personal_unlockAccount", param.From, param.Passwd, nil)
	if err != nil || response.Error != nil {
		fmt.Println(err)
	}

	fmt.Println(param.From, param.Passwd)
	fmt.Println(response)

	p := Payload{
		param.From,
		param.To,
		rh.Encode(), "0x1312d00", param.Parties}

	response, err = rpcClient.Call("eth_sendTransaction", []interface{}{p})

	fmt.Println(response)

	if err != nil || response.Error != nil {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response.Error)
		}
	}

	//fmt.Printf("%s", response.Result)
	return fmt.Sprintf("%s", response.Result)

}

func (ec *EthClient) EthCall(param contracthandler.ContractParam, encoder contracthandler.RequestHandler, decoder contracthandler.ResponseHandler) {
	rpcClient := jsonrpc.NewClient(ec.Url)

	p := CallPayload{param.To, encoder.Encode()}
	response, err := rpcClient.Call("eth_call", p, "latest")
	if err != nil {
		fmt.Println(err)
	}

	if response == nil {
		fmt.Println("No response returned!")
		return
	}

	decoder.Decode(fmt.Sprintf("%v", response.Result)[2:])
}

func (ec *EthClient) DeployContracts(byteCode string, pubKeys []string, private bool) string {
	coinbase := ec.Coinbase()
	var params contracthandler.ContractParam
	if private == true {
		params = contracthandler.ContractParam{From: coinbase, Passwd: "", Parties: pubKeys}
	} else {
		params = contracthandler.ContractParam{From: coinbase, Passwd: ""}
	}

	cont := contracthandler.DeployContractHandler{byteCode}
	txHash := ec.SendTransaction(params, cont)

	time.Sleep(1 * time.Second)

	contractAdd := ec.GetTransactionReceipt(txHash).ContractAddress
	for contractAdd == "" {
		time.Sleep(1 * time.Second)
		contractAdd = ec.GetTransactionReceipt(txHash).ContractAddress
	}
	return contractAdd
}

func (ec *EthClient) NetListening() bool {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("net_listening")
	if err != nil {
		fmt.Println(err)
	}
	var listening bool

	if response == nil {
		fmt.Println("No response returned!")
		return listening
	}

	err = response.GetObject(&listening)
	if err != nil {
		fmt.Println(err)
	}
	return listening
}

func (ec *EthClient) GetQuorumPayload(input string) string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getQuorumPayload", input)
	if err != nil {
		fmt.Println(err)
	}
	var payload string

	if response == nil {
		fmt.Println("No response returned!")
		return payload
	}

	err = response.GetObject(&payload)
	if err != nil {
		fmt.Println(err)
	}
	return payload
}

func (ec *EthClient) GetCode(address string) string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getCode", address, "latest")
	if err != nil {
		fmt.Println(err)
	}

	var bytecode string

	if response == nil {
		fmt.Println("No response returned!")
		return bytecode
	}

	err = response.GetObject(&bytecode)
	if err != nil {
		fmt.Println(err)
	}
	return bytecode
}

func (ec *EthClient) CreateAccount(password string) string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("personal_newAccount", password)
	if err != nil {
		fmt.Println(err)
	}
	var accountAddress string
	if response == nil {
		fmt.Println("No response returned!")
		return accountAddress
	}

	err = response.GetObject(&accountAddress)
	if err != nil {
		fmt.Println(err)
	}
	return accountAddress
}

func (ec *EthClient) GetAccounts() []string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_accounts")
	if err != nil {
		fmt.Println(err)
	}
	var accounts []string
	if response == nil {
		fmt.Println("No response returned!")
		return accounts
	}
	err = response.GetObject(&accounts)
	if err != nil {
		fmt.Println(err)
	}
	return accounts
}

func (ec *EthClient) GetBalance(account string) string {
	rpcClient := jsonrpc.NewClient(ec.Url)
	response, err := rpcClient.Call("eth_getBalance", account, "latest")
	if err != nil {
		fmt.Println(err)
	}
	var balance string
	if response == nil {
		fmt.Println("No response returned!")
		return balance
	}
	err = response.GetObject(&balance)
	if err != nil {
		fmt.Println(err)
	}
	return balance
}
