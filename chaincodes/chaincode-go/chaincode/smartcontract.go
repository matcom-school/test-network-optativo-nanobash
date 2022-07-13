package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
//Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
type File struct {
	AssetType string `json:"assetType"`
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Size      int    `json:"size"`
	State     string `json:"state"`
	Type      string `json:"type"`
	Url       string `json:"url"`
}

type User struct {
	CreatedAt string `json:"createdAt"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	AssetType string `json:"assetType"`
	State     string `json:"state"`
}

type FileHistory struct {
	Record    *File     `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	files := []File{
		{
			Size: 5, Owner: "tomoko@gmail.com", ID: "mockAssect1",
			CreatedAt: "2022-03-10",
			Type:      "mp3", Name: "Im not", Url: "http//:localhost:9000", AssetType: "File", State: "created"},
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
	id string, name string, createdAt string) error {

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

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the file %s does not exist", id)
	}

	var user User
	err = json.Unmarshal(assetJSON, &user)
	if err != nil {
		return nil, err
	}

	if user.AssetType != "User" {
		return nil, fmt.Errorf("the asset isn't file, it's a %s", user.AssetType)
	}

	return &user, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	queryString := `{"selector":{"assetType":"User"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return buildUsersListByQueryResponse(resultsIterator)
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func buildUsersListByQueryResponse(resultsIterator shim.StateQueryIteratorInterface) ([]*User, error) {
	var assets []*User
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset User
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}

		if asset.State != "deleted" {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

// Transaccionalidad en la blockchain ????
// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	user, err := s.ReadUser(ctx, id)
	if err != nil {
		return err
	}

	user.State = "deleted"

	fileJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	files, err := s.GetAllFilesByOwner(ctx, id)
	for _, file := range files {
		err := s.DeleteFile(ctx, file.ID)
		if err != nil {
			return err
		}
	}

	return ctx.GetStub().PutState(id, fileJSON)
}

//////////////////////////// FilesMethods ///////////////////////////////////////////////////////

func (s *SmartContract) CreateFile(ctx contractapi.TransactionContextInterface,
	id string, url string, name string, createdAt string, size int, owner string, type_ string) error {

	exists, err := s.FileExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	exists, err = s.UserExists(ctx, owner)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("the user %s doesn't exist", owner)
	}

	asset := File{
		CreatedAt: createdAt,
		ID:        id,
		Name:      name,
		Owner:     owner,
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
func (s *SmartContract) ReadFile(ctx contractapi.TransactionContextInterface, id string) (*File, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the file %s does not exist", id)
	}

	var file File
	err = json.Unmarshal(assetJSON, &file)
	if err != nil {
		return nil, err
	}

	if file.AssetType != "File" {
		return nil, fmt.Errorf("the asset isn't file, it's a %s", file.AssetType)
	}

	return &file, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateFile(ctx contractapi.TransactionContextInterface, id string,
	name string, url string, size int, type_ string) error {

	current_file, err := s.ReadFile(ctx, id)
	if err != nil {
		return err
	}

	// overwriting original asset with new asset
	current_file.Url = url
	current_file.Name = name
	current_file.Type = type_
	current_file.Size = size
	current_file.State = "modified"

	fileJSON, err := json.Marshal(current_file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, fileJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteFile(ctx contractapi.TransactionContextInterface, id string) error {
	file, err := s.ReadFile(ctx, id)
	if err != nil {
		return err
	}

	file.State = "deleted"

	fileJSON, err := json.Marshal(file)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, fileJSON)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferFile(ctx contractapi.TransactionContextInterface, id string, userId string) (string, error) {
	file, err := s.ReadFile(ctx, id)
	if err != nil {
		return "", err
	}

	exists, err := s.UserExists(ctx, userId)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("the user %s doesn't exist", userId)
	}

	oldOwner := file.Owner
	file.Owner = userId
	file.State = "transferd"

	fileJSON, err := json.Marshal(file)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, fileJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllFiles(ctx contractapi.TransactionContextInterface) ([]*File, error) {
	queryString := `{"selector":{"assetType":"File"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return buildFilesListByQueryResponse(resultsIterator)
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func buildFilesListByQueryResponse(resultsIterator shim.StateQueryIteratorInterface) ([]*File, error) {
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

		if asset.State != "deleted" {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

func (t *SmartContract) GetAllFilesByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*File, error) {
	queryString := fmt.Sprintf(`{"selector":{"assetType":"File","owner":"%s"}}`, owner)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return buildFilesListByQueryResponse(resultsIterator)
}

func (s *SmartContract) FilesHistory(ctx contractapi.TransactionContextInterface, id string) ([]FileHistory, error) {
	exists, err := s.FileExists(ctx, id)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	resultHistory, err := ctx.GetStub().GetHistoryForKey(id)
	defer resultHistory.Close()

	var history []FileHistory

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

		timestamp, err := ptypes.Timestamp(queryResponse.Timestamp)
		if err != nil {
			return nil, err
		}

		assetHist := FileHistory{
			TxId:      queryResponse.TxId,
			Timestamp: timestamp,
			Record:    &asset,
		}

		history = append(history, assetHist)

	}

	return history, nil
}
