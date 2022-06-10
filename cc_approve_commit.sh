#!/usr/bin/env sh
set -eu

# imports
. ./env.sh

export CORE_PEER_MSPCONFIGPATH="${PWD}"/crypto-config/organizations/org1.example.com/users/Admin@org1.example.com/msp

peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')" --package-id mycc:1.0
peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')"
peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID mychannel --name mycc --version 1.0 --sequence 1 --signature-policy "OR ('Org1MSP.member')" --peerAddresses 127.0.0.1:7051


# --init-required