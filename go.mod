module github.com/lifegoeson/fabric-sdk-go-practice

go 1.15

require (
	github.com/Shopify/sarama v1.28.0 // indirect
	github.com/fsouza/go-dockerclient v1.7.2 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hyperledger/fabric v1.4.3
	github.com/hyperledger/fabric-amcl v0.0.0-20210319225857-000ace5745f9 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	go.uber.org/zap v1.19.1 // indirect
)

//replace github.com/hyperledger/fabric-sdk-go v1.0.0 => ../fabric-sdk-go
