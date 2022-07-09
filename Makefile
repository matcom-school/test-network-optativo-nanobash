build:
	cd chaincodes/chaincode-go && go mod vendor && cd ../..

install:
	source peer1admin.sh && peer lifecycle chaincode package mycc.tar.gz --path ./chaincodes/chaincode-go --lang golang --label mycc
	peer lifecycle chaincode install mycc.tar.gz
	peer lifecycle chaincode queryinstalled

accept-cc:
	peer lifecycle chaincode approveformyorg -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1 --package-id ${CHAINCODE_ID} --sequence 1 --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 

deploys:	
	peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

ledger:
	# Inicializar el ledger con datos de prueba
	peer chaincode invoke -o 127.0.0.1:7050 -C mychannel -n mycc -c '{"Args":["InitLedger"]}' --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

create:
	# Crear un activo
	peer chaincode invoke -o 127.0.0.1:7050 -C mychannel -n mycc -c '{"Args":["CreateAsset","1","blue","35","tom","1000"]}' --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

read:
	# Consultar un activo
	peer chaincode query -C mychannel -n mycc -c '{"Args":["ReadAsset","1"]}'

update:
	# Actualizar un activo
	peer chaincode invoke -o 127.0.0.1:7050 -C mychannel -n mycc -c '{"Args":["UpdateAsset","1","blue","35","jerry","1000"]}' --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt



help:
	echo "> make build\n> make install\n> export CHAINCODE_ID=mycc:id\n> make accept-cc\n> make deploys\n> make ledger"