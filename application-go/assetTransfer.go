/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	log.Println("============ application-golang starts ============")

	var channelID = "mychannel"
	// Identidad con la que se conecta al nodo-par
	var identity = "User1@org1.example.com"

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists(identity) {
		log.Fatalf("Failed to populate wallet contents: %v", err)
	}

	ccpPath := filepath.Join(
		"ccp.yaml",
	)

	fmt.Println(ccpPath)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, identity),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelID)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("mycc")

	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
	result, err = contract.SubmitTransaction("CreateAsset", "asset13", "yellow", "5", "Tom", "1300")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadAsset, function returns an asset with a given assetID")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset13")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
	result, err = contract.EvaluateTransaction("AssetExists", "asset1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
	_, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}

	log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
	result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	log.Println("============ application-golang ends ============")


}