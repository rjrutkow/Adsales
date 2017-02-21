/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License .
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// This is our structure for the broadcaster creating bulk inventory
type adspot struct {
	uniqueAdspotId     string  `json:"uniqueAdspotId"`
	lotId              int     `json:"lotId"`
	adspotId           int     `json:"adspotId"`
	inventoryDate      string  `json:"inventoryDate"`
	programName        string  `json:"programName"`
	seasonEpisode      string  `json:"seasonEpisode"`
	broadcasterId      string  `json:"broadcasterId"`
	genre              string  `json:"genre"`
	dayPart            string  `json:"dayPart"`
	targetGrp          float64 `json:"targetGrp"`
	targetDemographics string  `json:"targetDemographics"`
	initialCpm         float64 `json:"initialCpm"`
	bsrp               float64 `json:"bsrp"`
	orderDate          string  `json:"orderDate"`
	adAgencyId         string  `json:"adAgencyId"`
	orderNumber        int     `json:"orderNumber"`
	advertiserId       string  `json:"advertiserId"`
	adContractId       int     `json:"adContractId"`
	adAssignedDate     string  `json:"adAssignedDate"`
	campaignName       string  `json:"campaignName"`
	campaignId         string  `json:"campaignId"`
	wasAired           string  `json:"wasAired"`
	airedDate          string  `json:"airedDate"`
	airedTime          string  `json:"airedTime"`
	actualGrp          float64 `json:"actualGrp"`
	actualProgramName  string  `json:"actualProgramName"`
	actualDemographics string  `json:"actualDemographics"`
	makupAdspotId      string  `json:"makupAdspotId"`
}

//This is a broadcaster's inventory
type releaseInventory struct {
	lotId              int     `json:"lotId"`
	adspotId           int     `json:"adspotId"`
	inventoryDate      string  `json:"inventoryDate"`
	programName        string  `json:"programName"`
	seasonEpisode      string  `json:"seasonEpisode"`
	broadcasterId      string  `json:"broadcasterId"`
	genre              string  `json:"genre"`
	dayPart            string  `json:"dayPart"`
	targetGrp          float64 `json:"targetGrp"`
	targetDemographics string  `json:"targetDemographics"`
	initialCpm         float64 `json:"initialCpm"`
	bsrp               float64 `json:"bsrp"`
	numberOfSpots      int     `json:"numberofSpots"`
	//releaseDate
	//UniqueAdspotID
}

//This is a pointer to allAdspots
type AllAdspots struct {
	uniqueAdspotId []string `json:"uniqueAdspotId"`
}

//For Debugging
func showArgs(args []string) {

	for i := 0; i < len(args); i++ {
		fmt.Printf("\n %d) : [%s]", i, args[i])
	}
	fmt.Printf("\n")
}

// Init function
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	fmt.Println("Launching Init Function")

	//Create array for all adspots in ledger
	var AllAdspotsArray AllAdspots

	jsonAsBytes, _ := json.Marshal(AllAdspotsArray)
	err = stub.PutState("BroadcasterA", jsonAsBytes)
	if err != nil {
		fmt.Println("Error Creating AllAdspotsArray")
		return nil, err
	}

	fmt.Println("Init Function Complete")
	return nil, nil
}

//STEP 1 Function - Replease Broadcaster's Inventory
func (t *SimpleChaincode) releaseInventory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Running releaseInventory")

	showArgs(args)

	//var broadcasterID string = args[0]

	//Outer Loop
	for i := 1; i < len(args); i++ {
		var releaseInventoryObj releaseInventory

		//Unmarshall each argument, assign to releaseInventoryObj
		b := []byte(args[i])
		err := json.Unmarshal(b, &releaseInventoryObj)

		if err != nil {
			fmt.Println("Error Unmarshalling arguments")
			return nil, err
		} else {
			fmt.Println(releaseInventoryObj.programName)
		}

	}

	//Inner Loop
	//for x := 1; x < releaseInventoryObj.numberOfSpots; x++ {
	//TO DO add args to addSpotStructure object
	//}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")

	showArgs(args)

	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	showArgs(args)

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")

	showArgs(args)

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	} else if function == "createBroadcasterLot" {
		fmt.Printf("Function is createBroadcasterLot")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// args: 0=BroadcasterID, 1=lotID, 2=Spots, 3=Ratings, 4=Demographics, 5=InitialPricePerSpot

func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")

	showArgs(args)

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from A to B
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("======== Query called, determining function")

	showArgs(args)

	if function != "Query" {
		fmt.Printf("Function is Query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

/*
func (t *SimpleChaincode) transferSpots(stub shim.ChaincodeStubInterface function string,, args []string) ([]byte, error) {

	fmt.Println("=============== entering transferSpots")

	var transferSpotsBroadcasterID string
	var spotsToTransfer int

	transferSpotsBroadcasterID = args[0]
	spotsToXfer, _ := strconv.Atoi(args[1])
	spotsToTransfer = spotsToXfer

	GetBroadCasterForTransfer, err := stub.GetState(transferSpotsBroadcasterID)
	if err != nil {
		return nil, errors.New("invalid broadcaster ID\n")
	}
	var BroadCasterForTransfer BulkAdSpot
	fmt.Println("Unmarshalling GetBroadCasterForTransfer")

	err = json.Unmarshal(GetBroadCasterForTransfer, &BroadCasterForTransfer)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("UnMarshalling Failed")
	}
	fmt.Println("Broadcaster ID is:%s", BroadCasterForTransfer.BroadcasterID)
	fmt.Println("Broadcaster Spots is:%s", BroadCasterForTransfer.Spots)

	BroadCasterForTransfer.Spots = BroadCasterForTransfer.Spots - spotsToTransfer

	jsonAsBytes, _ := json.Marshal(BroadCasterForTransfer)
	err = stub.PutState(BroadCasterForTransfer.BroadcasterID, jsonAsBytes)

	fmt.Println("Broadcaster ID is:%s", BroadCasterForTransfer.BroadcasterID)
	fmt.Println("Broadcaster New Spots balance is:%s", BroadCasterForTransfer.Spots)

	var AgencyForTransfer AdAgencyBulkAdSpots
	AgencyForTransfer.AgencyID = args[2]
	AgencyForTransfer.Demographics = BroadCasterForTransfer.Demographics
	AgencyForTransfer.InitialPricePerSpot = BroadCasterForTransfer.InitialPricePerSpot
	AgencyForTransfer.LotID = BroadCasterForTransfer.LotID
	AgencyForTransfer.OwnerBroadcasterID = BroadCasterForTransfer.BroadcasterID
	AgencyForTransfer.Ratings = BroadCasterForTransfer.Ratings
	AgencyForTransfer.Spots = spotsToTransfer

	jsonAsBytes2, _ := json.Marshal(AgencyForTransfer)
	err = stub.PutState(AgencyForTransfer.AgencyID, jsonAsBytes2)

	fmt.Println("New Agency ID is:%s", AgencyForTransfer.AgencyID)
	fmt.Println("Agency New Spots balance is:%s", AgencyForTransfer.Spots)

		// Get Receiver account from BC and update point balance
			rfidBytes, err := stub.GetState(tx.To)
			if err != nil {
				return nil, errors.New("transferPoints Failed to get Receiver from BC")
			}
			var receiver User
			fmt.Println("transferPoints Unmarshalling User Struct")
			err = json.Unmarshal(rfidBytes, &receiver)
			receiver.Balance = receiver.Balance + tx.Amount
			receiver.Modified = currentDateStr
			receiver.NumTxs = receiver.NumTxs + 1
			tx.ToName = receiver.Name


	return nil, nil
}
*/
func main() {
	err := shim.Start(new(SimpleChaincode))

	fmt.Printf("IN MAIN of example02 chaincode")
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

/*

Init time:


create key: allAdSpots_BroadcasterA, allAdSpots_AgencyA, allAdSpots_AdvertiserA, allAdSpots_AdvertiserC
       value: [UniqueAdspotID, UniqueAdSpotID...]


 type AllAdSpots struct{
 	Transactions []Transaction `json:"transactions"`
}
        To support Trace Ad Spot, we will need to also create the same for allAdSpots_AgencyA & allAdSpots_AdvertiserA & allAdSpots_AdvertiserC!

=================================

Release Inventory:

invoke function: releaseInventory

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "releaseInventory",
      "args": [
        "BroadcasterA",
        "{"UniqueAdspotID", "LotID", "AdspotID", "InventoryDate","ProgramName","SeasonEpisode","Genre","DayPart","TargetGRP","TargetDemographics","InitialCPMSpot","BSRP"}",
        "{"UniqueAdspotID", "LotID", "AdspotID", "InventoryDate","ProgramName","SeasonEpisode","Genre","DayPart","TargetGRP","TargetDemographics","InitialCPMSpot","BSRP"}",
        ...
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

chaincode:

type adSpot struct {					// this structue will store what is in Ron's spreadsheet
	UserId		string   `json:"UserId"`
	Name   		string   `json:"Name"`
	Balance 	float64  `json:"Balance"`
	NumTxs 	    int      `json:"NumberOfTransactions"`
	Status      string 	 `json:"Status"`
	Expiration  string   `json:"ExpirationDate"`
	Join		string   `json:"JoinDate"`
	Modified	string   `json:"LastModifiedDate"`
}

var UniqueAdSpotCounter = 1;
var BroadcasterID = args[0]
AllAdSpots = getState("allAdSpots_" + BroadcasterID)
for i = 1; i < args.length; i++ {
	var AdSpotData = args[i]
	for j = 0 j < AdSpotData.numberOFSpots; j++ {

		var adSpotStructure
		fill in adSpotStructure with AdSpotData

		var key = BroadcasterID + "_" + AdSpotData.lotID + "_" + UniqueAdSpotCounter
		putState(key, addSpotStructure)
		AllAdSpots.append(UniqueAdSpotID) 		// add this ad spot to the broadcaster array
		UniqueAdSpotCounter ++				// increment counter

	}
}

putState("allAdSpots_" + BroadcasterID,AllAdSpots)		// store the array with pointers to all the add spots

====================================

Place Orders

query function:  agencyViewLedger

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "agencyViewLedger",
      "args": [
        "AgencyId"
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

Result: { [ {"LotID","AdspotID","ProgramName","BroadcasterID","Genre","DayPart","TargetGRP","InitialCPMSpot","BSRP","NumberOfSpotsAvail"},
	    {"LotID","AdspotID","ProgramName","BroadcasterID","Genre","DayPart","TargetGRP","InitialCPMSpot","BSRP","NumberOfSpotsAvail"},
	    ...
	  ]
	}


chaincode:

var allAdspot = getState (AllAdSpots_BroadcasterA)
var returnData[]

for i = 0; i < allAdspots.length; i ++ {
   var adSpotData = getState(allAdspot[i].UniqueAdSpotID)

   allocate a structure to get subset of the entire adSpot Data
   fill in the structure with the required data

   add structure to returnData[]
}

return returnData[]



*** Assume line items are returned, the UI needs to aggragate

The UI will display the above fields addtional entry fields to be filled in by the user: AdvertiserID, OrderNumber, AdContractID, NumberOfSpotsToOrder

invoke function: placeOrders

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "placeOrders",
      "args": [

        "AgencyId",
        "BroadcasterID",
        "{"LotID", "ProgramName and/or AdspotID", "OrderNumber", "AdvertiserID", "AdContractID", "NumberOfSpotsToOrder"}",
        "{"LotID", "ProgramName and/or AdspotID", "OrderNumber", "AdvertiserID", "AdContractID", "NumberOfSpotsToOrder"}",
        ...
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

chaincode:

**** To Support Trace Ad Spot, we will need to also update global entry: "allAdSpots_" + agencyID and "allAd_Spots_" + advertiserID


var allAdSpotsBroadcaster = getState ("allAdSpots_" + broadcasterID)
var allAdSpotsAgency[] = getState("allAdSpots_" + agencyID)

for i = 2; i < args.length; i++ {		// loop through all the input data

	for j = 0; j < allAdSpotsBroadcaster.length; j++ {	// loop through all ad spots to find the first spot that is open
		var adSpotInfo = getState(allAdSpotsBroadcaster[j].UniqueAdSpotID)  	// retreive spot informationn
		if (adSpotInfo.AdSpotID == args[i].AdSpotID) && (adSpotInfo.AdContractID == "NA") {  // if ad spot is open

			for k = 0; k < args[i].NumberOfSpotsToOrder; k++ {	// loop until the number of spots is fullfilled
				adSpotInfo = getState(allAdSpots[j + k].UniqueAdSpotID)
				fill in add spot Info with purchase data
				putState (allAdSpots[j + k].UniqueAdSpotID, adspot with the addtions)  // save the addional data

				allAdSpotAgency[].append(allAdSpots[j + k].UniqueAdSpotID)	// add this ad spot to the agency array

				var allAdSpotsAdvertiser[] = getState("allAdSpots_" + args[i].AdvertiserID)  // add this ad spot to the advertiser array
				allAdSpotsAdvertiser.append(allAdSpots[j + k].UniqueAdSpotID)
				putState("allAdSpots_" + args[i].AdvertiserID, allAdSpotAdvertiser)
			}

			break 		// break of out of "j" loop
		}

	}

}
putState("allAdSpots_" + agencyID);

=====================================

MapAdsSpots

query function: getAdsSpotsToMap

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "getAdsSpotsToMap",
      "args": [
        "AgencyId" , "AdvertiserID" ***(optional)
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

Result: { [ {"UniqueAdspotID","BroadcasterID","AdContractID","AdvertiserID","TargetGRP","TargetDemographics","InitialCPMSpot"},
	    {"UniqueAdspotID","BroadcasterID","AdContractID","AdvertiserID","TargetGRP","TargetDemographics","InitialCPMSpot"},
	    ...
	  ]
	}

The UI will show the above filed and an entry field for Campaign

invoke function: mapAdsSpots

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "mapAdsSpots",
      "args": [
        "{"UniqueAdspotID", "AgencyId", "BroadcasterID","Campaign"}",
        "{"UniqueAdspotID", "AgencyId", "BroadcasterID","Campaign"}",
        ...
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

chaincode:
	for i = 0 ; i < args.length; i++ {
		ad spot info = getState(UniqueAdspotID)
		fill in ad spot info with passed data
		putState(UniqueAdspotID, ad spot info);
	}

======================================

 ***** NOT NEEDED *****

AdvertiserViewsLedger

query function: advertiserViewsLedger

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "advertiserViewsLedger",
      "args": [
        "AdveriserId"
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

Result: { [ {"LotID","ProgramName","BroadcasterID","Genre","DayPart","TargetGRP","TargetDemographics","InitialCPMSpot","BSRP","AdContractID","NumberOfSpotsAvail"},
	    {"LotID","ProgramName","BroadcasterID","Genre","DayPart","TargetGRP","TargetDemographics","InitialCPMSpot","BSRP","AdContractID","NumberOfSpotsAvail"},
	    ...
	  ]
	}

=======================================

AsRun

query function: prepForAsRun

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "prepForAsRun",
      "args": [
        "BroadcasterID"
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

Result: { [ {"UniqueAdspotID","AdContractID","Campaign","TargetGRP","TargetDemographics","ProgramName"},
 	    {"UniqueAdspotID","AdContractID","Campaign","TargetGRP","TargetDemographics","ProgramName"},
	    ...
	  ]
	}

chaincode:

AllAdsSpot = getState("allAdSpot");

var returnData[];

for i = 0; i < AllAdsSports.length; i++ {
	data = getState AllAdsSpots[i].UniqueAdspotId;
	retunData.append(data); // pass only what is needed
}

return retunData;

UI Shows the above fields and add the following entry fields: WasAired, AiredDate, AiredTime, ActualGRP, ActualDemographics, MakeUpAdSpotID

invoke function: asRun

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "asRun",
      "args": [
        "{"UniqueAdspotID", "WasAired", "AiredDate","AiredTime","ActualGRP","ActualDemographics", "MakeUpAdSpotID"}",
        "{"UniqueAdspotID", "WasAired", "AiredDate","AiredTime","ActualGRP","ActualDemographics", "MakeUpAdSpotID"}",
        ...
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

chaincode:

for i = 0; i < args.length; i++ {
	var AdSpotData = getState(args[i].UniqueAdspotID)
	fill in AdSpot data with data passed in
	putState(args[i].UniqueAdspotID, AdSpotData);

	if args[i].MapAdSpotID is pointing to another UniqueAdSpotID {

		var unAiredAdSpotData = getState(MakeUpAdSpotID)  // get the data for the unaired spot

		var allSpotsAgency[] = getState(unAiredAdSpotData.AgencyID)		// add to ad agency all transaction array
		allSpotsAgency.append(args[i].uniqueAdspotID)
		putState(unAiredAdSpotData.AgencyID, allSpotsAgency[])

		var allSpotsAdvertiser[] =  getState(unAiredAdSpotData.AdvertiserID)	// add to advertiser all transaction array
		allSpotsAdvertiser.apend(args[i].uniqueAdspotID)
		putState(unairedAdSpotData.AdvertiserID, allSpotsAdvertiser[])

	}
}


================================

Trace Ad Spot

query function: traceAdSpot

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "string"
    },
    "ctorMsg": {
      "function": "traceAdSpot",
      "args": [
        "AdvertiserID" or "BroadcasterID" or "AgencyID"
      ]
    },
    "secureContext": "string"
  },
  "id": 0
}

Result: { [ {"data for the UniqueAdspot...."} ],
	  [ {"data for the UniqueAdspot...."} ],
	  [ {"data for the UniqueAdspot...."},
 	    {"data for the makeup spot(s)"},
 	    {"data for the makeup spot(s)"},
	    ...
	  ]
	}

chaincode:

var AllAdSpots[] = getState("allAdSpots_" + PassedInID)

var returnAllData[]

for i = 0; i < AllAdSpots.length; i++ {
	var AdSpotData = getState(AllAdSpots[i])
	var returnSpotHistoryData[]
	returnEntry.append(AdSpotData)
	for j = 0; j < AllAdSpots.length; j++ {				// add makeup spots, if any
		var makeupAdSpotData = getState(AllAdSpots[j])
		if makeupAdSpotData.MakeUpAdSpotID == AdSpotData.UniqueAdSpotID {
			returnSpotHistoryData.append(makeupAdSpotData)
	        }
	}
	returnAllData.append(returnSpotHistory)
}

return returnAllData


*/
