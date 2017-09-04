# Local blockchain network

This repository contains source code for our local blockchain network used on http://iotplatformiicljubljana.mybluemix.net.

## Components

* [IBM Secure Gateway](https://www.ibm.com/blogs/bluemix/2017/03/secure-gateway-everything-ever-wanted-know/) 
* [Hyperledger on Docker](https://www.docker.com/)

This example is based on fabric sample balance-transfer (https://github.com/hyperledger/fabric-samples/tree/release/balance-transfer).

## Hyperledger network components

* 2 CAs
* A SOLO orderer
* 4 peers (2 peers per Org)
* IoT smart contract (https://github.com/evader1337/Blockchain-IoT_1)
* Datacenter smart contract (https://github.com/evader1337/Blockchain-DC)
* Spica smart contract (https://github.com/evader1337/Blockchain-Spica)

## How to use

```
sudo ./runApp.sh
```

This script will:
- destroy all previous hyperledger networks,
- restart the network,
- install node modules,
- start Secure Gateway service,
- start node.js REST api service,
- login users,
- create channel and join all peers on it,
- install and instantiate all smart contracts.

Please be patient, this can take a few minutes.

## Notes
All features except secure gateway should work out of the box. This instance of secure gateway is specifically tied to our bluemix account. If you want to set up your own, open runApp.sh and change your network id. After that, you need to access the gateway's console and add ACL rule to allow connections. In this case, you need to add yourip:4000 to fully allow all connections (route also needs to be configured on bluemix service).

*Adding secure gateway is optional. If your machine is publicly accessible on port 4000, this step is not required.*
