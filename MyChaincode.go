/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type VoteChaincode struct{

}


func (t *VoteChaincode) Init(stub shim.ChaincodeStubInterface,function string,args []string) ([]byte,error) {
	if len(args)!=1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var maxCount int
	var err error
	maxCount,err=strconv.Atoi(args[0])
	if err!=nil {
		return nil, errors.New("Invalid transaction amount, expecting a integer value")
	}
	err=stub.PutState("max",[]byte(strconv.Itoa(maxCount)))
	if err!=nil {
		return nil,err
	}
	err = stub.CreateTable("Candidate", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err!=nil {
		return nil, errors.New("Failed creating Candidate table.")
	}
	err = stub.CreateTable("Vote", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Id", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "VId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name:"Time",Type: shim.ColumnDefinition_STRING,Key:false},
	})
	if err!=nil {
		return nil, errors.New("Failed creating Vote table.")
	}

	return nil,nil
}

func (t *VoteChaincode) Invoke(stub shim.ChaincodeStubInterface,function string,args []string) ([]byte,error){

	if function=="addCandidate" {
		return t.addCandidate(stub,args)
	}
	if function=="vote" {
		return t.vote(stub,args)

	}
	return nil,nil
}

func (t *VoteChaincode) addCandidate(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {
	return nil, errors.New("Incorrect number of arguments. Expecting 2")
	if len(args) !=2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var Id string
	var Name string
	var err error
	var  countByte []byte
	Id=args[0]
	Name=args[1]
	countByte,err=stub.GetState(Id)
	if err==nil {
		jsonResp := "{\"Error\":\"Failed to add candidate for " + Name + "\"}"
		return nil, errors.New(jsonResp)
	}
	if countByte!=nil {
		jsonResp := "{\"Error\":\"Failed to add candidate for " + Name + "\"}"
		return nil, errors.New(jsonResp)
	}
	err=stub.PutState(Id,[]byte(strconv.Itoa(0)))
	if err!=nil {
		return nil,err
	}
	/*var ok bool
	ok,err=stub.InsertRow("Candidate",shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Id}},
			&shim.Column{Value: &shim.Column_String_{String_: Name}},
			},
	})
	if !ok {
		stub.DelState(Id)
		jsonResp := "{\"Error\":\"Failed to add candidate for " + Name + "\"}"
		return nil, errors.New(jsonResp)
	}*/
	return nil,nil
}

func (t *VoteChaincode) vote(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {
	//CId,VId,time
	if len(args)!=3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	var CId string
	var VId string
	var time string
	var CountByte []byte
	var count int
	var err error
	CId=args[0]
	VId=args[1]
	time=args[2]
	CountByte,err=stub.GetState(VId)
	//no vote
	if err!=nil {
		count=0

	} else {
		count,err=strconv.Atoi(string(CountByte))
		if err!=nil {
			return nil,err
		}

	}
	var maxbyte []byte
	maxbyte,err=stub.GetState("max")
	var max int
	max,err=strconv.Atoi(string(maxbyte))
	if count==max {
		return nil, errors.New("Can not vote")
	}
	count=count+1

	var ok bool
	ok,err=stub.InsertRow("",shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: CId}},
			&shim.Column{Value: &shim.Column_String_{String_: VId}},
			&shim.Column{Value: &shim.Column_String_{String_: time}},
		},
	})
	if !ok {
		return nil, err
	}
	err=stub.PutState(VId,[]byte(strconv.Itoa(count)))
	if err!=nil {
		return nil,err
	}

	var CCount int
	var CCountByte []byte
	CCountByte,err=stub.GetState(CId)
	if err!=nil {
		return nil,errors.New("No this Candidate")
	}

	CCount,_=strconv.Atoi(string(CCountByte))
	CCount=CCount+1
	err=stub.PutState(CId,[]byte(strconv.Itoa(CCount)))


	return nil,nil
}



func (t *VoteChaincode) Query(stub shim.ChaincodeStubInterface,function string,args []string) ([]byte,error){
	if function=="candidate" {
		return t.CandidateQuery(stub,args)

	} else if function=="total" {
		return t.TotalQuery(stub,args)
	} else if function=="votenum" {
		return t.VoteNum(stub,args)
	} else if function=="voter" {
		return t.Voter(stub,args)
	}
	return nil, errors.New("Invalid query function name.")
}

func (t *VoteChaincode) CandidateQuery(stub shim.ChaincodeStubInterface,args []string) ([]byte,error){
	if len(args)!=1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var CId string
	var err error
	var countByte []byte
	CId=args[0]
	countByte,err=stub.GetState(CId)
	if err!=nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + CId+ "\"}"
		return nil, errors.New(jsonResp)
	}
	if countByte==nil {
		jsonResp := "{\"Error\":\"Nil amount for " + CId + "\"}"
		return nil, errors.New(jsonResp)
	}
	jsonResp := "{\"Name\":\"" + CId+ "\",\"Amount\":\"" + string(countByte) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return countByte, nil
}
func (t *VoteChaincode) Voter(stub shim.ChaincodeStubInterface,args []string) ([]byte,error){
	if len(args)!=1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	VId:=args[0]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: VId}}
	columns = append(columns, col1)
	rows, err := stub.GetRows("vote", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving asset [%s]", err)
	}
	var resp string
	resp=""
	for row:=range rows {
		fmt.Println(row)
		resp=resp+" "+string(row.Columns[0].GetBytes())
		if len(rows)<=0 {
			break;
		}
	}
	return []byte(resp),nil
}
func (t *VoteChaincode) TotalQuery(stub shim.ChaincodeStubInterface,  args []string) ([]byte,error){
	if len(args)!=1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	CId:=args[0]
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: CId}}
	columns = append(columns, col1)
	rows, err := stub.GetRows("vote", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving asset [%s]", err)
	}
	var resp string
	resp=""
	for row:=range rows {
		fmt.Println(row)
		resp=resp+" "+string(row.Columns[1].GetBytes())
		if len(rows)<=0 {
			break;
		}
	}
	return []byte(resp),nil
}

func (t *VoteChaincode) VoteNum(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args)!=0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	maxByte,_:=stub.GetState("max")

	return maxByte,nil
}

func main() {
	err := shim.Start(new(VoteChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
