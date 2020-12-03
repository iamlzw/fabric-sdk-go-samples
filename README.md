# fabric-sdk-go-samples
fabric sdk go samples

基于hyperledger fabric v1.4.2测试网络(其他版本应该也没什么问题，可能需要修改文件中的相关路径)

主要是参考fabric-sdk-go中end_to_end.go以及ledger_queries_test.go中的代码进行修改

测试了查询链，查询区块，解析区块等,没有测试创建channel安装以及实例化链码等功能,

首先确保能将first-network成功运行

克隆仓库

```
$ git clone https://github.com/iamlzw/fabric-sdk-go-samples.git
```

启动测试

```
$ cd fabric-sdk-go-samples
$ go run main.go setup.go
```

输出结果类似于
```
2020/12/03 01:41:26 Initialized channel client
&{0xc00007fe80 0xc000418140 0xc00043d5c0 0xc00043d080 0xc000414000}
&{0xc00041a740 <nil> 0xb237a0 0xc00007fa40}

Org1MSP�-----BEGIN CERTIFICATE-----
MIICKTCCAc+gAwIBAgIRAMPOsdVbFb4BB3d7DgyWzgYwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjAxMjAzMDkyOTAwWhcNMzAxMjAxMDkyOTAw
WjBqMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzENMAsGA1UECxMEcGVlcjEfMB0GA1UEAxMWcGVlcjAub3Jn
MS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABKaIQ6vv+9uq
28+5lnFQ1nmojqfh1n0FY/2we9NeiVSPXn3v45kHOwcq015PIfBeYjQVyGFQtSZS
Z0m9UGWFMDSjTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1Ud
IwQkMCKAIJu6j+LzRWojal3xQgIeFKG9s/Ylm3LLQSXrjI89y/03MAoGCCqGSM49
BAMCA0gAMEUCIQDeqQIOHIX75eQJlPPk7wyAYOQFm6z/W6uGEBowKqoP4wIgMHkz
SkDYBYVCQ29nvDuRMBjdqpRA5xQ0MOTUVJcco9c=
-----END CERTIFICATE-----

�M%�3lL�<&|�������ԃ�r�Eǋc}CQ�|b ,i��a������=aq�L%

```
就测试成功了。


### 遇到的问题
1、只测试了在ubuntu平台上，在windows上测试连接部署在虚拟机上的区块链网络时，报错"user not found",暂无解决方法
2、报错"Description: dialing connection on target [peer0.org1.example.com:7051]: connection is in TRANSIENT_FAILURE"
  解决：参考https://stackoverflow.com/questions/50051193/hyper-fabric-fabric-sdk-go-error-connection-failed-description-dialing-conn
  虽然报错不太一样，但是解决办法是一样的，修改"entityMatchers"下的属性
