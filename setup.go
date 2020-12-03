package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"io/ioutil"
	"log"
	"strconv"
)


var defaultQueryArgs = [][]byte{[]byte("b")}

// Initialize reads the configuration file and sets up the client, chain and event hub
func Initialize() error {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := resmgmt.New(clientContext)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}


	ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	log.Println("Initialized channel client")

	ledgerClient, err := ledger.New(ccp)

	// Test Query Info - retrieve values before transaction
	chainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Println("QueryInfo return error: %s", err)
	}
	height := chainInfo.BCI.Height
	for i := int64(0); i < int64(height); i++ {
		block ,err := ledgerClient.QueryBlock(uint64(i))
		if err != nil{
			fmt.Println(err)
		}
		buf := new (bytes.Buffer)
		err = protolator.DeepMarshalJSON(buf, block)
		if err != nil {
			fmt.Errorf("malformed block contents: %s", err)
		}
		filename := "mychannel_"+strconv.FormatInt(i,10)+".json"
		err = ioutil.WriteFile(filename,buf.Bytes(),0644)
		if err != nil{
			fmt.Println("write to file failure:",err)
		}

	}
	fmt.Println(cc)
	fmt.Println(client)

	response, err := cc.Query(channel.Request{ChaincodeID: "mycc", Fcn: "query", Args: [][]byte{[]byte("a")}},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(),
	)
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}

	txResponse := response.Responses
	endorsement := txResponse[0].Endorsement
	endorser := endorsement.Endorser
	signature := endorsement.Signature

	fmt.Println(string(endorser))
	fmt.Println(string(signature))

	return nil

	//// Channel management client is responsible for managing channels (create/update channel)
	//// Supply user that has privileges to create channel (in this case orderer admin)
	//clientContext := sdk.Context(fabsdk.WithUser("User1"), fabsdk.WithOrg("Org1"))
	//
	//resMgmtClient, err := resmgmt.New(clientContext)
	//
	//if err != nil {
	//	fmt.Println("failed to create resMgmt")
	//}
	//
	//mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("Org1"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//req := resmgmt.SaveChannelRequest{ChannelID: "mychannel",
	//	ChannelConfigPath: "channel-artifacts/mychannel.tx",
	//	SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	//txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	//
	//fmt.Println(txID)
	//
	//return nil
}
