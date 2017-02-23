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
	UniqueAdspotId     string  `json:"uniqueAdspotId"`
	LotId              int     `json:"lotId"`
	AdspotId           int     `json:"adspotId"`
	InventoryDate      string  `json:"inventoryDate"`
	ProgramName        string  `json:"programName"`
	SeasonEpisode      string  `json:"seasonEpisode"`
	BroadcasterId      string  `json:"broadcasterId"`
	Genre              string  `json:"genre"`
	DayPart            string  `json:"dayPart"`
	TargetGrp          float64 `json:"targetGrp"`
	TargetDemographics string  `json:"targetDemographics"`
	InitialCpm         float64 `json:"initialCpm"`
	Bsrp               float64 `json:"bsrp"`
	OrderDate          string  `json:"orderDate"`
	AdAgencyId         string  `json:"adAgencyId"`
	OrderNumber        int     `json:"orderNumber"`
	AdvertiserId       string  `json:"advertiserId"`
	AdContractId       int     `json:"adContractId"`
	AdAssignedDate     string  `json:"adAssignedDate"`
	CampaignName       string  `json:"campaignName"`
	CampaignId         string  `json:"campaignId"`
	WasAired           string  `json:"wasAired"`
	AiredDate          string  `json:"airedDate"`
	AiredTime          string  `json:"airedTime"`
	ActualGrp          float64 `json:"actualGrp"`
	ActualProgramName  string  `json:"actualProgramName"`
	ActualDemographics string  `json:"actualDemographics"`
	MakupAdspotId      string  `json:"makupAdspotId"`
}

//This is a broadcaster's inventory
type releaseInventory struct {
	LotId              string `json:"lotId"`
	AdspotId           string `json:"adspotId"`
	InventoryDate      string `json:"inventoryDate"`
	ProgramName        string `json:"programName"`
	SeasonEpisode      string `json:"seasonEpisode"`
	BroadcasterId      string `json:"broadcasterId"`
	Genre              string `json:"genre"`
	DayPart            string `json:"dayPart"`
	TargetGrp          string `json:"targetGrp"`
	TargetDemographics string `json:"targetDemographics"`
	InitialCpm         string `json:"initialCpm"`
	Bsrp               string `json:"bsrp"`
	NumberOfSpots      string `json:"numberofSpots"`
	//releaseDate
	//UniqueAdspotID
}

type queryPlaceOrders struct {
	LotId              int     `json:"lotId"`
	AdspotId           int     `json:"adspotId"`
	ProgramName        string  `json:"programName"`
	BroadcasterId      string  `json:"broadcasterId"`
	Genre              string  `json:"genre"`
	DayPart            string  `json:"dayPart"`
	TargetGrp          float64 `json:"targetGrp"`
	TargetDemographics string  `json:"targetDemographics"`
	InitialCpm         float64 `json:"initialCpm"`
	Bsrp               float64 `json:"bsrp"`
	NumberOfSpots      int     `json:"numberofSpots"`
}

//To place an order for the adspots
type placeOrders struct {
	LotId         string `json:"lotId"`
	AdspotId      string `json:"adspotId"`
	OrderNumber   string `json:"orderNumber"`
	ProgramName   string `json:"programName"`
	AdvertiserId  string `json:"advertiserId"`
	AdContractId  string `json:"adContractId"`
	NumberOfSpots string `json:"numberofSpots"`
}

//This is a pointer to allAdspots
type AllAdspots struct {
	UniqueAdspotId []string `json:"uniqueAdspotId"`
}

const noData string = "NA"
const noValue int = -1

//For Debugging
func showArgs(args []string) {

	for i := 0; i < len(args); i++ {
		fmt.Printf("\n %d) : [%s]", i, args[i])
	}
	fmt.Printf("\n")
}

// Init function
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Launching Init Function")

	broadcasterId := "BroadcasterA"
	agencyId := "AgencyA"
	advertiser1Id := "AdvertiserA"
	advertiser2Id := "AdvertiserB"

	//Create array for all adspots in ledger
	var AllAdspotsArray AllAdspots

	t.putAllAdspotPointers(stub, AllAdspotsArray, broadcasterId)
	t.putAllAdspotPointers(stub, AllAdspotsArray, agencyId)
	t.putAllAdspotPointers(stub, AllAdspotsArray, advertiser1Id)
	t.putAllAdspotPointers(stub, AllAdspotsArray, advertiser2Id)

	fmt.Println("Init Function Complete")
	return nil, nil
}

//STEP 1 Function - Replease Broadcaster's Inventory
func (t *SimpleChaincode) releaseInventory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Running releaseInventory")

	var broadcasterID = args[0]
	var increment = 1

	fmt.Println(broadcasterID)

	allAdspotPointers, _ := t.getAllAdspotPointers(stub, broadcasterID)

	//Outer Loop
	for i := 1; i < len(args); i++ {

		var in = args[i]

		bytes := []byte(in)
		var releaseInventoryObj releaseInventory
		err := json.Unmarshal(bytes, &releaseInventoryObj)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v", releaseInventoryObj)
		fmt.Printf("\n program name: %v \n", releaseInventoryObj.ProgramName)

		NumberOfSpots, _ := strconv.Atoi(releaseInventoryObj.NumberOfSpots)

		for x := 0; x < NumberOfSpots; x++ {
			var ThisAdspot adspot

			ThisAdspot.UniqueAdspotId = ("1000_" + strconv.Itoa(increment))
			ThisAdspot.LotId, _ = strconv.Atoi(releaseInventoryObj.LotId)
			ThisAdspot.AdspotId, _ = strconv.Atoi(releaseInventoryObj.AdspotId)
			ThisAdspot.InventoryDate = releaseInventoryObj.InventoryDate
			ThisAdspot.ProgramName = releaseInventoryObj.ProgramName
			ThisAdspot.SeasonEpisode = releaseInventoryObj.SeasonEpisode
			ThisAdspot.BroadcasterId = broadcasterID
			ThisAdspot.Genre = releaseInventoryObj.Genre
			ThisAdspot.DayPart = releaseInventoryObj.DayPart

			ThisAdspot.TargetGrp, _ = strconv.ParseFloat(releaseInventoryObj.TargetGrp, 64)
			ThisAdspot.TargetDemographics = releaseInventoryObj.TargetDemographics
			ThisAdspot.InitialCpm, _ = strconv.ParseFloat(releaseInventoryObj.InitialCpm, 64)
			ThisAdspot.Bsrp, _ = strconv.ParseFloat(releaseInventoryObj.Bsrp, 64)
			ThisAdspot.OrderDate = noData
			ThisAdspot.AdAgencyId = noData
			ThisAdspot.OrderNumber = noValue
			ThisAdspot.AdvertiserId = noData
			ThisAdspot.AdContractId = noValue
			ThisAdspot.AdAssignedDate = noData
			ThisAdspot.CampaignName = noData
			ThisAdspot.CampaignId = noData
			ThisAdspot.WasAired = noData
			ThisAdspot.AiredDate = noData
			ThisAdspot.AiredTime = noData
			ThisAdspot.ActualGrp = float64(noValue)
			ThisAdspot.ActualProgramName = noData
			ThisAdspot.ActualDemographics = noData
			ThisAdspot.MakupAdspotId = noData

			increment++
			fmt.Printf("ThisAdspot: %+v ", ThisAdspot)
			fmt.Printf("\n")
			allAdspotPointers.UniqueAdspotId = append(allAdspotPointers.UniqueAdspotId, ThisAdspot.UniqueAdspotId)

			t.putAdspot(stub, ThisAdspot)
		}

	}

	t.putAllAdspotPointers(stub, allAdspotPointers, broadcasterID)
	return nil, nil
}

//STEP 2 Function - Place Orders for ad spots
func (t *SimpleChaincode) placeOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Running placeOrders")
	showArgs(args)

	agencyId := args[0]
	broadcasterId := args[1]

	broadcasterAllAdspotsPointers, _ := t.getAllAdspotPointers(stub, broadcasterId)
	agencyAllAdsportsPointers, _ := t.getAllAdspotPointers(stub, agencyId)

	// loop through all entries
	for i := 2; i < len(args); i++ {

		// loop through the ad contracts
		in := args[i]
		bytes := []byte(in)
		var placeOrdersObj placeOrders
		err := json.Unmarshal(bytes, &placeOrdersObj)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Place Orders Object: %+v", placeOrdersObj)

		// now look through the inventory of ad spots, find match
		for j := 0; j < len(broadcasterAllAdspotsPointers.UniqueAdspotId); j++ {
			uniqueAdspotKey := broadcasterAllAdspotsPointers.UniqueAdspotId[j]
			AdSpotObj, _ := t.getAdspot(stub, uniqueAdspotKey)

			spotid, _ := strconv.Atoi(placeOrdersObj.AdspotId)

			advertiserAllAdsportsPointers, _ := t.getAllAdspotPointers(stub, placeOrdersObj.AdvertiserId)

			if AdSpotObj.AdspotId == spotid && AdSpotObj.AdContractId == noValue { // found adspot and  ad spot not taken
				numberOfSpotsToPurchase, _ := strconv.Atoi(placeOrdersObj.NumberOfSpots)

				fmt.Printf("Inside if AdSpotObj.AdspotId == spotid && AdSpotObj.AdContractId == noValue \n")
				fmt.Printf("AdsPotObj.AdspotId: %v \n", AdSpotObj.AdspotId)
				fmt.Printf("spotid: %v \n", spotid)
				fmt.Printf("AdSpotObj.AdContractId: %v \n", AdSpotObj.AdContractId)
				fmt.Println("END IF")

				for k := 0; k < numberOfSpotsToPurchase; k++ {
					if k > 0 { // get the correct ad spot if needed
						uniqueAdspotKey = broadcasterAllAdspotsPointers.UniqueAdspotId[j+k]
						AdSpotObj, _ = t.getAdspot(stub, uniqueAdspotKey)
					}
					AdSpotObj.AdAgencyId = agencyId
					AdSpotObj.AdvertiserId = placeOrdersObj.AdvertiserId
					AdSpotObj.AdContractId, _ = strconv.Atoi(placeOrdersObj.AdContractId)
					AdSpotObj.OrderNumber, _ = strconv.Atoi(placeOrdersObj.OrderNumber)

					t.putAdspot(stub, AdSpotObj)

					// save all pointers for appropriate ad agency
					agencyAllAdsportsPointers.UniqueAdspotId = append(agencyAllAdsportsPointers.UniqueAdspotId, AdSpotObj.UniqueAdspotId)

					// save all pointers for appropriate advertiser
					advertiserAllAdsportsPointers.UniqueAdspotId = append(advertiserAllAdsportsPointers.UniqueAdspotId, AdSpotObj.UniqueAdspotId)
				}
				t.putAllAdspotPointers(stub, advertiserAllAdsportsPointers, placeOrdersObj.AdvertiserId)
				break // break out of for j loop
			} else {
				fmt.Printf("NOHIT -  if AdSpotObj.AdspotId == spotid && AdSpotObj.AdContractId == noValue \n")
				fmt.Printf("AdsPotObj.AdspotId: %v \n", AdSpotObj.AdspotId)
				fmt.Printf("spotid: %v \n", spotid)
				fmt.Printf("AdSpotObj.AdContractId: %v \n", AdSpotObj.AdContractId)
				fmt.Println("END NO HIT on IF")
			}
		}
	}

	t.putAllAdspotPointers(stub, agencyAllAdsportsPointers, agencyId)

	return nil, nil
}

func (t *SimpleChaincode) putAdspot(stub shim.ChaincodeStubInterface, adspotObj adspot) ([]byte, error) {
	//marshalling
	fmt.Println("Launching putAdspot helper function")
	fmt.Printf("putAdspot obj: %+v ", adspotObj)
	fmt.Printf("\n")

	bytes, _ := json.Marshal(adspotObj)
	err := stub.PutState(adspotObj.UniqueAdspotId, bytes)
	if err != nil {
		fmt.Println("Error - could not Marshall in putAdspot")
		//return nil, err
	} else {
		fmt.Println("Success - putAdspot putState works")
	}
	fmt.Println("putAdspot Function Complete")
	return nil, nil
}

func (t *SimpleChaincode) getAdspot(stub shim.ChaincodeStubInterface, uniqueAdspotId string) (adspot, error) {
	//unmarshalling
	fmt.Println("Launching getAdspot helper function")
	bytes, err := stub.GetState(uniqueAdspotId)
	if err != nil {
		fmt.Println("Error - Could not get Unique Adspot ID: %s", uniqueAdspotId)
		//return nil, err
	} else {
		fmt.Println("Success - getAdspot getState worked with Unique Adspot ID %s", uniqueAdspotId)
	}

	var adspotObj adspot
	err = json.Unmarshal(bytes, &adspotObj)
	if err != nil {
		fmt.Println("Error - could not Unmarshall in getAdspot - uniqueAdspotID %s", uniqueAdspotId)
	} else {
		fmt.Println("Success - Unmarshall in getAdspot good - uniqueAdspotID %s", uniqueAdspotId)
	}

	fmt.Printf("getAdspot: %+v ", adspotObj)
	fmt.Printf("\n")
	fmt.Println("getAdspot Function Complete")
	return adspotObj, err
}

func (t *SimpleChaincode) getAllAdspotPointers(stub shim.ChaincodeStubInterface, userId string) (AllAdspots, error) {
	//unmarshalling
	fmt.Println("Launching getAllAdspotPointers helper function  -  userid: ", userId)
	bytes, err := stub.GetState(userId)
	if err != nil {
		fmt.Println("Error - Could not get Broadcaster ID")
		//return nil, err
	} else {
		fmt.Println("Success - got Broadcaster ID")
	}

	var allAdspotPointers AllAdspots
	err = json.Unmarshal(bytes, &allAdspotPointers)
	if err != nil {
		fmt.Println("Error - could not Unmarshall within getAllAdspotPointers")
	} else {
		fmt.Println("Success - Unmarshall within getAllAdspotPointers")
	}

	fmt.Printf("allAdspotsObj: %+v ", allAdspotPointers)
	fmt.Printf("\n")

	fmt.Println("getAllAdspotPointers Function Complete - userid: ", userId)

	return allAdspotPointers, err
}

func (t *SimpleChaincode) putAllAdspotPointers(stub shim.ChaincodeStubInterface, allAdspotsObj AllAdspots, userId string) ([]byte, error) {
	//marshalling
	fmt.Println("Launching putAllAdspotPointers helper function userid: ", userId)
	fmt.Printf("putAllAdspotPointers: %+v ", allAdspotsObj)
	fmt.Printf("\n")
	bytes, _ := json.Marshal(allAdspotsObj)
	err := stub.PutState(userId, bytes)
	if err != nil {
		fmt.Println("Error - could not Marshall in putAllAdspotPointers")
		//return nil, err
	} else {
		fmt.Println("Success - Marshall in putAllAdspotPointers")
	}
	fmt.Println("putAllAdspotPointers Function Complete - userid: ", userId)
	return nil, nil
}

// -----------------------------------------------------------------------------------------------------
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
	if function == "releaseInventory" {
		fmt.Printf("Function is releaseInventory")
		return t.releaseInventory(stub, args)
	} else if function == "placeOrders" {
		fmt.Printf("Function is placeOrders")
		return t.placeOrders(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// args: 0=BroadcasterID, 1=lotID, 2=Spots, 3=Ratings, 4=Demographics, 5=InitialPricePerSpot

/*
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
*/

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
