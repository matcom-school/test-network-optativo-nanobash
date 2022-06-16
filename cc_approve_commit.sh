#!/usr/bin/env sh
set -eu

echo "Este script es solo para el proceso de instalacion en la ORG-1"
echo "Para instalar el chaincode en la ORG-2 no es necesario ejecutar este script"
echo "Solo ejecutar:"
echo "source peer2admin.sh"
echo "peer lifecycle chaincode install mycc.tar.gz"

# look for binaries in local dev environment /build/bin directory and then in local samples /bin directory
export PATH="${PWD}"/../bin:"$PATH"
export FABRIC_CFG_PATH="${PWD}"/../config

export CORE_PEER_ID=peer0.org1.example.com
export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=127.0.0.1:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/organizations/org1.example.com/users/Admin@org1.example.com/msp

peer lifecycle chaincode approveformyorg -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1 --package-id $CHAINCODE_ID --sequence 1 --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt


# --init-required