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

export FABRIC_LOGGING_SPEC=INFO
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ID=peer0.org1.example.com
export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/organizations/org1.example.com/peers/peer0.org1.example.com/msp
export CORE_PEER_ADDRESS=127.0.0.1:7051

peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')" --package-id mycc:1.0
peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')"
peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')" --peerAddresses 127.0.0.1:7051


# --init-required