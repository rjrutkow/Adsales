package main

//Packages to import followed by a Pointer to your hyperledger installation
import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//CONSTANTS --------------------------------------------------------------------------------------------------------------------------
//These are defined Constants for use throughout the gocode
const noData string = "NA" //Defalt for empty string values
const noValue int = -1     //Default for empty numerical values

//STRUCTURES --------------------------------------------------------------------------------------------------------------------------
// SimpleChaincode required structure
type SimpleChaincode struct {
}

// This is our primary structure for Adspots, based on columns defined within the ledger template.
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

//This is a helper structure for releasing Adspots (STEP 1)
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
	//releaseDate - do we need this?
	//UniqueAdspotID - to we need this?
}

// This is a helper structure for querying placed orders (STEP 2)
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

//This is a helper structure to place an order for the adspots (STEP 2)
type placeOrders struct {
	LotId         string `json:"lotId"`
	AdspotId      string `json:"adspotId"`
	OrderNumber   string `json:"orderNumber"`
	ProgramName   string `json:"programName"`
	AdvertiserId  string `json:"advertiserId"`
	AdContractId  string `json:"adContractId"`
	NumberOfSpots string `json:"numberofSpots"`
}

//This is a helper structure to point to allAdspots
type AllAdspots struct {
	UniqueAdspotId []string `json:"uniqueAdspotId"`
}

// This is a helper structure for querying placed orders (STEP 2)
type queryPlaceOrdersStruc struct {
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
}

// This is a helper structure for querying placed orders (STEP 2)
type queryPlaceOrdersArray struct {
	PlacedOrderData []queryPlaceOrdersStruc `json:"placedOrderData"`
}

//This is a helper structure for querying adspots before mapping(STEP 3)
type queryAdspotsToMapArray struct {
	AdspotsToMapData []queryAdspotsToMapStruct `json:"adspotsToMapData"`
}

//This is a helper structure for mapping adspots (STEP 3)
type queryAdspotsToMapStruct struct {
	UniqueAdspotId     string  `json:"uniqueAdspotId"`
	BroadcasterId      string  `json:"broadcasterId"`
	AdContractId       int     `json:"adContractId"`
	CampaignName       string  `json:"campaignName"`
	AdvertiserId       string  `json:"advertiserId"`
	TargetGrp          float64 `json:"targetGrp"`
	TargetDemographics string  `json:"targetDemographics"`
	InitialCpm         float64 `json:"initialCpm"`
}

//For Debugging
func showArgs(args []string) {

	for i := 0; i < len(args); i++ {
		fmt.Printf("\n %d) : [%s]", i, args[i])
	}
	fmt.Printf("\n")
}

//ADSALES USE CASE FUNCTIONS --------------------------------------------------------------------------------------------------------------------------

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

			ThisAdspot.UniqueAdspotId = (releaseInventoryObj.LotId + "_" + strconv.Itoa(increment))
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
// NEEDS TESTING
func (t *SimpleChaincode) placeOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Running placeOrders")
	showArgs(args)

	agencyId := args[0]
	broadcasterId := args[1]

	broadcasterAllAdspotsPointers, _ := t.getAllAdspotPointers(stub, broadcasterId)
	agencyAllAdspotsPointers, _ := t.getAllAdspotPointers(stub, agencyId)

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

			if (AdSpotObj.AdspotId == spotid) && (AdSpotObj.AdContractId == noValue) { // found adspot and  adspot not already taken
				numberOfSpotsToPurchase, _ := strconv.Atoi(placeOrdersObj.NumberOfSpots)

				fmt.Printf("Inside if AdSpotObj.AdspotId == spotid && AdSpotObj.AdContractId == noValue \n")
				fmt.Printf("AdsPotObj.AdspotId: %v \n", AdSpotObj.AdspotId)
				fmt.Printf("spotid: %v \n", spotid)
				fmt.Printf("AdSpotObj.AdContractId: %v \n", AdSpotObj.AdContractId)
				fmt.Println("END IF")

				//Loop on number of spots to purchase
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
					agencyAllAdspotsPointers.UniqueAdspotId = append(agencyAllAdspotsPointers.UniqueAdspotId, AdSpotObj.UniqueAdspotId)

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

	t.putAllAdspotPointers(stub, agencyAllAdspotsPointers, agencyId)

	return nil, nil
}

func (t *SimpleChaincode) queryPlaceOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//agencyId := args[0]
	fmt.Println("Launching queryPlaceOrders Function")
	broadcasterId := args[1]
	var queryPlaceOrdersArrayObj queryPlaceOrdersArray

	broadcasterAllAdspotsPointers, _ := t.getAllAdspotPointers(stub, broadcasterId)

	for i := 0; i < len(broadcasterAllAdspotsPointers.UniqueAdspotId); i++ {
		var queryPlaceOrdersStrucObj queryPlaceOrdersStruc
		ThisAdspot, _ := t.getAdspot(stub, broadcasterAllAdspotsPointers.UniqueAdspotId[i])

		//Code Fix: If statement needs to be checked
		if ThisAdspot.AdContractId == noValue {
			queryPlaceOrdersStrucObj.AdspotId = ThisAdspot.AdspotId
			queryPlaceOrdersStrucObj.BroadcasterId = ThisAdspot.BroadcasterId
			queryPlaceOrdersStrucObj.Bsrp = ThisAdspot.Bsrp
			queryPlaceOrdersStrucObj.DayPart = ThisAdspot.DayPart
			queryPlaceOrdersStrucObj.Genre = ThisAdspot.Genre
			queryPlaceOrdersStrucObj.InitialCpm = ThisAdspot.InitialCpm
			queryPlaceOrdersStrucObj.LotId = ThisAdspot.LotId
			queryPlaceOrdersStrucObj.ProgramName = ThisAdspot.ProgramName
			queryPlaceOrdersStrucObj.TargetDemographics = ThisAdspot.TargetDemographics
			queryPlaceOrdersStrucObj.TargetGrp = ThisAdspot.TargetGrp
			queryPlaceOrdersArrayObj.PlacedOrderData = append(queryPlaceOrdersArrayObj.PlacedOrderData, queryPlaceOrdersStrucObj)
		}
	}

	jsonAsBytes, err := json.Marshal(queryPlaceOrdersArrayObj)
	if err != nil {
		fmt.Println("Error returning json output for queryPlaceOrders ")
		return nil, err
	}

	fmt.Println("queryPlaceOrders Function Complete")
	fmt.Printf("queryPlaceOrdersArrayObj: %+v ", queryPlaceOrdersArrayObj)
	return jsonAsBytes, nil
}

func (t *SimpleChaincode) queryAdspotsToMap(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Launching queryAdspotsToMap Function")

	agencyId := args[0]
	var queryAdspotsToMapArrayObj queryAdspotsToMapArray

	agencyAllAdspotsPointers, _ := t.getAllAdspotPointers(stub, agencyId)

	for i := 0; i < len(agencyAllAdspotsPointers.UniqueAdspotId); i++ {
		var queryAdspotsToMapStructObj queryAdspotsToMapStruct
		ThisAdspot, _ := t.getAdspot(stub, agencyAllAdspotsPointers.UniqueAdspotId[i])

		queryAdspotsToMapStructObj.UniqueAdspotId = ThisAdspot.UniqueAdspotId
		queryAdspotsToMapStructObj.BroadcasterId = ThisAdspot.BroadcasterId
		queryAdspotsToMapStructObj.AdContractId = ThisAdspot.AdContractId
		queryAdspotsToMapStructObj.CampaignName = ThisAdspot.CampaignName
		queryAdspotsToMapStructObj.AdvertiserId = ThisAdspot.AdvertiserId
		queryAdspotsToMapStructObj.TargetGrp = ThisAdspot.TargetGrp
		queryAdspotsToMapStructObj.TargetDemographics = ThisAdspot.TargetDemographics
		queryAdspotsToMapStructObj.InitialCpm = ThisAdspot.InitialCpm
		queryAdspotsToMapArrayObj.AdspotsToMapData = append(queryAdspotsToMapArrayObj.AdspotsToMapData, queryAdspotsToMapStructObj)
	}

	jsonAsBytes, err := json.Marshal(queryAdspotsToMapArrayObj)
	if err != nil {
		fmt.Println("Error returning json output for queryAdspotsToMap")
		return nil, err
	}

	fmt.Println("queryAdspotsToMap Function Complete")
	fmt.Printf("queryAdspotsToMapArrayObj: %+v ", queryAdspotsToMapArrayObj)
	return jsonAsBytes, nil
}

//HELPER FUNCTIONS --------------------------------------------------------------------------------------------------------------------------
//putAdspot: To put data back to the ledger
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

//getAdspot: To get data back from the ledger
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

//getAllAdspotPointers: To get an array containing pointers to all blocks for a particular user(or peer) from the ledger
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

//getAllAdspotPointers: To put an array containing pointers to all blocks for a particular user(or peer) on the ledger
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

//REQUIRED FUNCTIONS --------------------------------------------------------------------------------------------------------------------------
// INIT FUNCTION
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Launching Init Function")

	//Peers hard coded here
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

//INVOKE FUNCTION
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

//QUERY FUNCTION
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("======== Query called, determining function")

	showArgs(args)

	if function == "queryPlaceOrders" {
		fmt.Printf("Function is queryPlaceOrders")
		return t.queryPlaceOrders(stub, args)
	} else if function == "queryAdspotsToMap" {
		fmt.Printf("Function is queryAdspotsToMap")
		return t.queryAdspotsToMap(stub, args)
	} else {
		fmt.Printf("Invalid Function!")
	}

	return nil, nil
}

//MAIN FUNCTION
func main() {
	err := shim.Start(new(SimpleChaincode))

	fmt.Printf("IN MAIN of AdsalesChaincode")
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
