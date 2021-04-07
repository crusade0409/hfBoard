/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	// "strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Meet struct {
	// Mid int `json:"mid"`
	MTitle  string `json:"mtitle"`
	MContents string `json:"mcontents"`
	Owner  string `json:"owner"`
	DateTime  string `json:"datetime"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Meet
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	//now := time.Now().String()

	meets := []Meet{
		// Meet{MTitle: "1월 회의", MContents: "1월 회의내용", Owner: "이대국", DateTime: "2021-03-25"},
		// Meet{MTitle: "2월 회의", MContents: "2월 회의내용", Owner: "이대국", DateTime: "2021-03-25"},
		// Meet{MTitle: "3월 회의", MContents: "3월 회의내용", Owner: "이대국", DateTime: "2021-03-25"},
		// Meet{Mid:4,MTitle: "4월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:5,MTitle: "5월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:6,MTitle: "6월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:7,MTitle: "7월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:8,MTitle: "8월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:9, MTitle: "9월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:10,MTitle: "10월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
		// Meet{Mid:11,MTitle: "11월 회의", MContents: "회의내용", Owner: "이대국", DateTime: now},
	}

	for i, meet := range meets {
		meetAsBytes, _ := json.Marshal(meet)		
		//fmt.Printf("Item%4d", i) // "age: 0007"
		tmpKey := fmt.Sprintf("Item%4d", i)
		err := ctx.GetStub().PutState(tmpKey, meetAsBytes)
		//err := ctx.GetStub().PutState("List"+strconv.Itoa(i), meetAsBytes)
		
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, mNumber string, mtitle string, mcontents string, owner string, datetime string) error {
	fmt.Printf("Item%4d", 1);
	meet := Meet{		
		MTitle:   mtitle,
		MContents:  mcontents,
		Owner: owner,
		DateTime:  datetime,
	}

	meetAsBytes, _ := json.Marshal(meet)

	return ctx.GetStub().PutState(mNumber, meetAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, mNumber string) (*Meet, error) {
	meetAsBytes, err := ctx.GetStub().GetState(mNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if meetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", mNumber)
	}

	meet := new(Meet)
	_ = json.Unmarshal(meetAsBytes, meet)

	return meet, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		meet := new(Meet)
		_ = json.Unmarshal(queryResponse.Value, meet)

		queryResult := QueryResult{Key: queryResponse.Key, Record: meet}
		results = append(results, queryResult)		
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, mNumber string, newOwner string) error {
	meet, err := s.QueryCar(ctx, mNumber)

	if err != nil {
		return err
	}

	meet.Owner = newOwner

	meetAsBytes, _ := json.Marshal(meet)

	return ctx.GetStub().PutState(mNumber, meetAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
