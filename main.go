package main

func main(){
	sdk := initSDK()

	//createChannel(sdk)
	//joinChannel(sdk)
	//
	//createCC(sdk)

	invokeChaincode(sdk)
	queryChaincode(sdk)
}
