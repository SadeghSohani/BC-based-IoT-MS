package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

type Participant struct {
	PublicKey string `json:"publicKey"`
	Role      string `json:"role"`
	Link      string `json:"link"`
}

type IOTLocalNetwork struct {
	PublicKey string `json:"publicKey"`
	Owner     string `json:"owner"`
	AreaType  string `json:"areaType"`
}

type Asset struct {
	Id      string `json:"Id"`
	Holder  string `json:"Holder"`  // participant
	Owner   string `json:"owner"`   // participant with customer role
	Station string `json:"station"` // IOTLocalNetwork
}

type PartiHistoryModel struct {
	TxId        string       `json:"txId"`
	Participant *Participant `json:"paritcipant"`
	Timestamp   string       `json:"timestamp"`
	IsDelete    bool         `json:"isDelete"`
}

type LocalNetworkHistoryModel struct {
	TxId         string           `json:"txId"`
	LocalNetwork *IOTLocalNetwork `json:"localNetwork"`
	Timestamp    string           `json:"timestamp"`
	IsDelete     bool             `json:"isDelete"`
}

type AssetHistoryModel struct {
	TxId      string `json:"txId"`
	Asset     *Asset `json:"asset"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) AddParticipant(ctx contractapi.TransactionContextInterface,
	publicKey string, role string, link string) (*Participant, error) {

	parti := Participant{PublicKey: publicKey, Role: role, Link: link}
	assetAsBytes, _ := json.Marshal(parti)
	err := ctx.GetStub().PutState(publicKey, assetAsBytes)
	if err != nil {
		return nil, err
	}

	return &parti, nil
}

func (s *SmartContract) QueryParticipant(ctx contractapi.TransactionContextInterface, publicKey string) (*Participant, error) {
	assetAsBytes, err := ctx.GetStub().GetState(publicKey)

	if err != nil {
		return nil, err
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", publicKey)
	}

	parti := new(Participant)
	_ = json.Unmarshal(assetAsBytes, parti)

	return parti, nil
}

func (s *SmartContract) QueryAllParticipants(ctx contractapi.TransactionContextInterface) ([]Participant, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Participant{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		parti := new(Participant)
		if err = json.Unmarshal(queryResponse.Value, parti); err == nil && parti.Role != "" {
			results = append(results, *parti)
		}
	}

	return results, nil
}

func (t *SmartContract) GetParticipantHistory(ctx contractapi.TransactionContextInterface, publicKey string) ([]PartiHistoryModel, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(publicKey)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	results := []PartiHistoryModel{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		parti := new(Participant)
		_ = json.Unmarshal(queryResponse.Value, parti)

		historyItem := PartiHistoryModel{
			TxId:        queryResponse.TxId,
			Participant: parti,
			Timestamp:   time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String(),
			IsDelete:    queryResponse.IsDelete}
		results = append(results, historyItem)
	}

	return results, nil
}

func (s *SmartContract) ChangeParticipantLink(ctx contractapi.TransactionContextInterface,
	publicKey string, link string) (*Participant, error) {

	parti, err := s.QueryParticipant(ctx, publicKey)
	if err != nil {
		return nil, err
	}

	parti.Link = link

	assetAsBytes, _ := json.Marshal(parti)

	_err := ctx.GetStub().PutState(publicKey, assetAsBytes)

	if _err != nil {
		return nil, err
	}

	return parti, nil
}

func (s *SmartContract) AddIOTLocalNetwork(ctx contractapi.TransactionContextInterface,
	publicKey string, owner string, areaType string) (*IOTLocalNetwork, error) {

	localNet := IOTLocalNetwork{PublicKey: publicKey, Owner: owner, AreaType: areaType}

	assetAsBytes, _ := json.Marshal(localNet)
	err := ctx.GetStub().PutState(publicKey, assetAsBytes)

	if err != nil {
		return nil, err
	}

	return &localNet, nil
}

func (s *SmartContract) QueryIOTLocalNetwork(ctx contractapi.TransactionContextInterface, publicKey string) (*IOTLocalNetwork, error) {
	assetAsBytes, err := ctx.GetStub().GetState(publicKey)

	if err != nil {
		return nil, err
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", publicKey)
	}

	localNet := new(IOTLocalNetwork)
	_ = json.Unmarshal(assetAsBytes, localNet)

	return localNet, nil
}

func (s *SmartContract) QueryIOTLocalNetworkByOwner(ctx contractapi.TransactionContextInterface, publicKey string, owner string) (*IOTLocalNetwork, error) {
	localNet, err := s.QueryIOTLocalNetwork(ctx, publicKey)

	if err != nil {
		return nil, err
	}

	if localNet.Owner == owner {
		return localNet, nil
	} else {
		return nil, fmt.Errorf("Permission denied.")
	}
}

func (s *SmartContract) QueryAllLocalNetworks(ctx contractapi.TransactionContextInterface) ([]IOTLocalNetwork, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []IOTLocalNetwork{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		localNet := new(IOTLocalNetwork)
		if err = json.Unmarshal(queryResponse.Value, localNet); err == nil {
			if localNet.AreaType != "" {
				results = append(results, *localNet)
			}
		}
	}

	return results, nil
}

func (t *SmartContract) GetLocalNetworkHistory(ctx contractapi.TransactionContextInterface,
	publicKey string) ([]LocalNetworkHistoryModel, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(publicKey)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	results := []LocalNetworkHistoryModel{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		localNet := new(IOTLocalNetwork)
		_ = json.Unmarshal(queryResponse.Value, localNet)

		historyItem := LocalNetworkHistoryModel{
			TxId:         queryResponse.TxId,
			LocalNetwork: localNet,
			Timestamp:    time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String(),
			IsDelete:     queryResponse.IsDelete}
		results = append(results, historyItem)
	}

	return results, nil
}

func (s *SmartContract) ChangeLocalNetworkOwner(ctx contractapi.TransactionContextInterface,
	publicKey string, owner string, newOwner string) (*IOTLocalNetwork, error) {
	localNet, err := s.QueryIOTLocalNetwork(ctx, publicKey)

	if err != nil {
		return nil, err
	}

	if localNet.Owner == owner {
		localNet.Owner = newOwner
	}

	assetAsBytes, _ := json.Marshal(localNet)

	_err := ctx.GetStub().PutState(publicKey, assetAsBytes)

	if _err != nil {
		return nil, _err
	}

	return localNet, nil
}

func (s *SmartContract) ChangeLocalNetworkAreaType(ctx contractapi.TransactionContextInterface,
	publicKey string, owner string, newAreaType string) (*IOTLocalNetwork, error) {
	localNet, err := s.QueryIOTLocalNetwork(ctx, publicKey)

	if err != nil {
		return nil, err
	}

	if localNet.Owner == owner {
		localNet.AreaType = newAreaType
	}

	assetAsBytes, _ := json.Marshal(localNet)

	_err := ctx.GetStub().PutState(publicKey, assetAsBytes)

	if _err != nil {
		return nil, _err
	}

	return localNet, nil
}

func (s *SmartContract) AddAsset(ctx contractapi.TransactionContextInterface,
	id string, holederKey string, ownerKey string, iotLNKey string) (*Asset, error) {

	asset := Asset{Id: id, Holder: holederKey, Owner: ownerKey, Station: iotLNKey}

	assetAsBytes, _ := json.Marshal(asset)
	err := ctx.GetStub().PutState(id, assetAsBytes)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) QueryAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, err
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	asset := new(Asset)
	_ = json.Unmarshal(assetAsBytes, asset)

	return asset, nil
}

// Change it
func (s *SmartContract) QueryAssetByOwnerOrHolder(ctx contractapi.TransactionContextInterface,
	id string, publicKey string) (*Asset, error) {

	asset, err := s.QueryAsset(ctx, id)

	if err != nil {
		return nil, err
	}

	if asset.Owner == publicKey || asset.Holder == publicKey {
		return asset, nil
	} else {
		return nil, fmt.Errorf("Permission denied.")
	}
}

func (s *SmartContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]Asset, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Asset{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		asset := new(Asset)
		if err = json.Unmarshal(queryResponse.Value, asset); err == nil && asset.Id != "" {
			results = append(results, *asset)
		}
	}
	return results, nil
}

func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface,
	id string) ([]AssetHistoryModel, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	results := []AssetHistoryModel{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		asset := new(Asset)
		_ = json.Unmarshal(queryResponse.Value, asset)

		historyItem := AssetHistoryModel{
			TxId:      queryResponse.TxId,
			Asset:     asset,
			Timestamp: time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String(),
			IsDelete:  queryResponse.IsDelete}
		results = append(results, historyItem)
	}

	return results, nil
}

func (s *SmartContract) ChangeAssetOwner(ctx contractapi.TransactionContextInterface,
	id string, owner string, newOwner string) (*Asset, error) {

	asset, err := s.QueryAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	if asset.Owner == owner {
		newOwnerP, err := s.QueryParticipant(ctx, newOwner)
		if err != nil {
			return nil, err
		}
		oldOwnerP, err := s.QueryParticipant(ctx, owner)
		if err != nil {
			return nil, err
		}

		station, err := s.QueryIOTLocalNetwork(ctx, asset.Station)
		if err != nil {
			return nil, err
		}

		if station != nil {
			if err := ctx.GetStub().SetEvent(station.PublicKey,
				[]byte(`Stop:`+oldOwnerP.Link)); err != nil {
				return nil, err
			}
			if err := ctx.GetStub().SetEvent(station.PublicKey,
				[]byte(`Send:`+newOwnerP.Link)); err != nil {
				return nil, err
			}
		}

		asset.Owner = newOwner
		if assetAsBytes, err := json.Marshal(asset); err == nil {
			if err := ctx.GetStub().PutState(asset.Id, assetAsBytes); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return asset, nil
}

func (s *SmartContract) ChangeAssetHolder(ctx contractapi.TransactionContextInterface,
	id string, holder string, newHolder string) (*Asset, error) {

	asset, err := s.QueryAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	if asset.Holder == holder {
		newHolderP, err := s.QueryParticipant(ctx, newHolder)
		if err != nil {
			return nil, err
		}
		oldHolderP, err := s.QueryParticipant(ctx, holder)
		if err != nil {
			return nil, err
		}

		station, err := s.QueryIOTLocalNetwork(ctx, asset.Station)
		if err != nil {
			return nil, err
		}
		
		if station != nil {
			if err := ctx.GetStub().SetEvent(station.PublicKey,
				[]byte(`Stop:`+oldHolderP.Link)); err != nil {
				return nil, err
			}
			if err := ctx.GetStub().SetEvent(station.PublicKey,
				[]byte(`Send:`+newHolderP.Link)); err != nil {
				return nil, err
			}
		}

		asset.Holder = newHolder
		if assetAsBytes, err := json.Marshal(asset); err == nil {
			if err := ctx.GetStub().PutState(id, assetAsBytes); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return asset, nil
}

func (s *SmartContract) ChangeAssetStation(ctx contractapi.TransactionContextInterface,
	id string, holder string, newIOTLNKey string) (*Asset, error) {

	asset, err := s.QueryAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	if asset.Holder == holder {
		partiH, err := s.QueryParticipant(ctx, holder)
		if err != nil {
			return nil, err
		}
		partiO, err := s.QueryParticipant(ctx, asset.Owner)
		if err != nil {
			return nil, err
		}
		holderLink := partiH.Link
		ownerLink := partiO.Link

		newStation, err := s.QueryIOTLocalNetwork(ctx, newIOTLNKey)
		if err != nil {
			return nil, err
		}
		oldStation, err := s.QueryIOTLocalNetwork(ctx, asset.Station)
		if err != nil {
			return nil, err
		}
		
		if oldStation != nil {
			if err := ctx.GetStub().SetEvent(oldStation.PublicKey,
				[]byte(`Stop:`+ownerLink)); err != nil {
				return nil, err
			}
			if err := ctx.GetStub().SetEvent(oldStation.PublicKey,
				[]byte(`Stop:`+holderLink)); err != nil {
				return nil, err
			}
		}
		if newStation != nil {
			if err := ctx.GetStub().SetEvent(newStation.PublicKey,
				[]byte(`Send:`+ownerLink)); err != nil {
				return nil, err
			}
			if err := ctx.GetStub().SetEvent(newStation.PublicKey,
				[]byte(`Send:`+holderLink)); err != nil {
				return nil, err
			}
		}

		asset.Station = newIOTLNKey
		if assetAsBytes, err := json.Marshal(asset); err != nil {
			return nil, err
		} else {
			if err := ctx.GetStub().PutState(id, assetAsBytes); err != nil {
				return nil, err
			}
		}
	}

	return asset, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting broilerChickenCC Smart Contract: %s", err.Error())
	}
}