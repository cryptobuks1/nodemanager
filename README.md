# PowerChain Maker Node Manager

PowerChain Network Manager allows users to create and manage a Nordic Energy network of nodes that use IBFT as consensus mechanism.

#### Prerequisites

##### solc
sudo add-apt-repository ppa:ethereum/ethereum  
sudo apt-get update  
sudo apt-get install solc  

##### abigen
go get -u github.com/ethereum/go-ethereum  
cd $GOPATH/src/github.com/ethereum/go-ethereum/  
make  
make devtools  

#### Create Smart Contract ABI
run
```
solc --abi --overwrite --optimize NetworkManagerContract.sol --output-dir contractclient/internalcontract
```

#### Create Smart Contract go class 
run
```
abigen --abi=contractclient/internalcontract/NetworkManagerContract.abi --pkg=ScClient --out=contractclient/internalcontract/NetworkManagerContract.go 
```
// Then replace imports in NetworkManagerContract.go from "github.com/ethereum/go-ethereum" to "github.com/ethereum/go-ethereum"
