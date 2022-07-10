package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/vendor/github.com/hyperledger/fabric-chaincode-go/shim"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
//Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type File struct {
	AssetType string   `json:"assetType`
	CreatedAt string   `json:"createdAt"`
	Customers []string `json:"customers"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Owner     string   `json:"owner"`
	Size      int      `json:"size"`
	State     string   `json:"State"`
	Type      string   `json:"type"`
	Url       string   `json:"url"`
}

type User struct {
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	AssetType string `json:"assetType"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	files := []File{
		{
			Size: 5, Owner: "tomoko@gmail.com", ID: "mockAssect1",
			CreatedAt: "2022-03-10", Customers: make([]string, 0),
			Type: "mp3", Name: "Im not", Url: "http//:localhost:9000", AssetType: "File"},
	}

	for _, asset := range files {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	users := []User{
		{
			ID:        "tomoko@gmail.com",
			CreatedAt: "2022-03-10",
			Name:      "tomoko",
			AssetType: "User",
		},
		{
			ID:        "emilio@gmail.com",
			CreatedAt: "2022-03-10",
			Name:      "emilio",
			AssetType: "User",
		},
	}

	for _, asset := range users {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

////////////////////////////// UserMethods //////////////////////////////////////////

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface,
	id string, name, string, createdAt string) error {

	exists, err := s.UserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := User{
		CreatedAt: createdAt,
		ID:        id,
		Name:      name,
		AssetType: "User",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		return false, nil
	}

	var user User
	err = json.Unmarshal(assetJSON, &user)
	if err != nil {
		return false, err
	}

	return user.AssetType == "User", nil
}

//////////////////////////// FilesMethods ///////////////////////////////////////////////////////

func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface,
	id string, url string, name, string, createdAt string, size int, owner string, type_ string) error {

	exists, err := s.FileExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := File{
		CreatedAt: createdAt,
		ID:        id,
		Name:      name,
		Owner:     owner,
		Customers: nil,
		Url:       url,
		Size:      size,
		Type:      type_,
		AssetType: "File",
		State:     "created",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) FileExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		return false, nil
	}

	var file File
	err = json.Unmarshal(assetJSON, &file)
	if err != nil {
		return false, err
	}

	return file.AssetType == "File", nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*File, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset File
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string,
	name string, url string, size int, type_ string) error {

	current_file, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	// overwriting original asset with new asset
	current_file.Url = url
	current_file.Name = name
	current_file.Type = type_
	current_file.Size = size
	current_file.State = "modified"

	assetJSON, err := json.Marshal(current_file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	file, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	file.State = "deleted"

	assetJSON, err := json.Marshal(file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner
	asset.State = "transferd"

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*File, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*File
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset File
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*File, error) {
	var assets []*File
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset File
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*File, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func (t *SmartContract) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*File, error) {
	queryString := fmt.Sprintf(`{"selector":{"AssetType":"File","owner":"%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) AssetHisoty(ctx contractapi.TransactionContextInterface, id string) ([]*File, error) {
	exists, err := s.FileExists(ctx, id)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	resultHistory, err := ctx.GetStub().GetHistoryForKey(id)

	var assets []*File

	for resultHistory.HasNext() {
		queryResponse, err := resultHistory.Next()
		if err != nil {
			return nil, err
		}
		var asset File
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)

	}

	return assets, nil
}
