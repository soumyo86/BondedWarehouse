
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)



 
// CBDoc is a high level smart contract that POCs together business artifact based smart contracts
type CBDoc struct {

}

// UserDetails is for storing User Details

type Document struct{	
	DocumentId string `json:"documentId"`
	Source string `json:"source"`
	Destination string `json:"destination"`
	Status string `json:"status"`
	DeliveryDate string `json:"deliveryDate"`
	DeliveryNo string `json:"deliveryNo"`
	CustPO string `json:"custPO"`
	Documenthash string `json:"documenthash"`
	ReasonCode string `json:"reasonCode"`
	TruckId string `json:"truckId"`
	CustomLocation string `json:"customLocation"`
	SourceLatLong string `json:"sourceLatLong"`
	DestnationLatLong string `json:"destnationLatLong"`
	CustomLatLong string `json:"customLatLong"`
}

// TrxnHistory is for storing document status change history
type TrxnHistory struct{	
	TrxId string `json:"trxId"`
	TimeStamp string `json:"timeStamp"`
	DocumentId string `json:"documentId"`
	UpdatedBy string `json:"updatedBy"`
	Status string `json:"status"`	
}

// ItemDetails is for storing document status change history
type ItemDetails struct{	
	ItemId string `json:"itemId"`
	DocumentId string `json:"documentId"`
	Name string `json:"name"`
	Quantity string `json:"quantity"`
	Description string `json:"description"`
	Weightvolume string `json:"weightvolume"`	
}

// ItemTracker is for storing tracking details of the items
type ItemTracker struct{	
	DocumentId string `json:"documentId"`
	TruckId string `json:"truckId"`
	LeftFactory string `json:"leftFactory"`
	ArrivingCustom string `json:"arrivingCustom"`
	ReachedCustom string `json:"reachedCustom"`
	LeavingCustom string `json:"leavingCustom"`
	ReceivedWarehouse string `json:"receivedWarehouse"`
}

type Counter struct {
	count int
}

func (self Counter) currentValue() int {
	return self.count
}

func (self *Counter) increment() {
	self.count++
}


// Init initializes the smart contracts
func (t *CBDoc) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check if table already exists
	_, err := stub.GetTable("Document")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("Document", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "documentId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "source", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "destination", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "deliveryDate", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "deliveryNo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "custPO", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "documenthash", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "reasonCode", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "truckId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "customLocation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "sourceLatLong", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "destnationLatLong", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "customLatLong", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Document table.")
	}

	// Check if table already exists
	_, err = stub.GetTable("TrxnHistory")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("TrxnHistory", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "trxId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "timeStamp", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "documentId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "updatedBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "status", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating TrxnHistory table.")
	}	

	// Check if table already exists
	_, err = stub.GetTable("ItemDetails")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("ItemDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "itemId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "documentId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "name", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "quantity", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "description", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "weightvolume", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating ItemDetails table.")
	}
	
	// Check if table already exists
	_, err = stub.GetTable("ItemTracker")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("ItemTracker", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "documentId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "truckId", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "leftFactory", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "arrivingCustom", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "reachedCustom", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "leavingCustom", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "receivedWarehouse", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating ItemTracker table.")
	}
	
	// setting up the increment for Transaction Id
	stub.PutState("Trx_increment", []byte("1"))
	stub.PutState("Item_increment", []byte("1"))
	
	return nil, nil
}

// generate booking number for shipping item
func (t *CBDoc) createDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 14 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 14. Got: %d.", len(args))
		}
		
		documentId:=args[0]
		source:=args[1]
		destination:=args[2]
		timeStamp:=args[3]
		updatedBy:=args[4]
		documenthash:=args[5] 
		status:="Initial"
		deliveryDate:=args[6]
		deliveryNo:=args[7]
		custPO:=args[8]
		reasonCode:=""
		truckId:=args[9]
		customLocation:=args[10]
		sourceLatLong:=args[11]
		destnationLatLong:=args[12]
		customLatLong:=args[13]
		leftFactory:=""
		arrivingCustom:=""
		reachedCustom:=""
		leavingCustom:=""
		receivedWarehouse:=""
		
		//getting TrxId incrementer
		Avalbytes, err := stub.GetState("Trx_increment")
		Aval, _ := strconv.ParseInt(string(Avalbytes), 10, 0)
		newAval:=int(Aval) + 1
		newTrx_increment:= strconv.Itoa(newAval)
		stub.PutState("Trx_increment", []byte(newTrx_increment))
		
		trxId:=string(Avalbytes)

		// Insert a row
		ok, err := stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: source}},
				&shim.Column{Value: &shim.Column_String_{String_: destination}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryNo}},
				&shim.Column{Value: &shim.Column_String_{String_: custPO}},
				&shim.Column{Value: &shim.Column_String_{String_: documenthash}},
				&shim.Column{Value: &shim.Column_String_{String_: reasonCode}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: customLocation}},
				&shim.Column{Value: &shim.Column_String_{String_: sourceLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: destnationLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: customLatLong}},
			}})

		if err != nil {
			return nil, err
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		// Insert details in TrxnHistory table
		ok, err = stub.InsertRow("TrxnHistory", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: trxId}},
				&shim.Column{Value: &shim.Column_String_{String_: timeStamp}},
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: updatedBy}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},			
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
		
		// Insert details in ItemTracker table
		ok, err = stub.InsertRow("ItemTracker", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: leftFactory}},
				&shim.Column{Value: &shim.Column_String_{String_: arrivingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: reachedCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: leavingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: receivedWarehouse}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
			
		return nil, nil

}	

// insert line item details
func (t *CBDoc) insertItemList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 5 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 5. Got: %d.", len(args))
		}
		
		documentId:=args[0]
		name:=args[1]
		quantity:=args[2]
		description:=args[3]
		weightvolume:=args[4]
				
		//getting ItemId incrementer
		Avalbytes, err := stub.GetState("Item_increment")
		Aval, _ := strconv.ParseInt(string(Avalbytes), 10, 0)
		newAval:=int(Aval) + 1
		newItem_increment:= strconv.Itoa(newAval)
		stub.PutState("Item_increment", []byte(newItem_increment))
		
		itemId:=string(Avalbytes)

		// Insert a row
		ok, err := stub.InsertRow("ItemDetails", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: itemId}},
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: name}},
				&shim.Column{Value: &shim.Column_String_{String_: quantity}},
				&shim.Column{Value: &shim.Column_String_{String_: description}},
				&shim.Column{Value: &shim.Column_String_{String_: weightvolume}},			
			}})

		if err != nil {
			return nil, err
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
			
		return nil, nil

}

//get all item list for specified documentId
func (t *CBDoc) viewItemListByDocumentId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document id to query")
	}

	documentId := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("ItemDetails", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*ItemDetails{}	
	
	for row := range rows {		
		newApp:= new(ItemDetails)
		newApp.ItemId = row.Columns[0].GetString_()
		newApp.DocumentId = row.Columns[1].GetString_()
		newApp.Name = row.Columns[2].GetString_()
		newApp.Quantity = row.Columns[3].GetString_()
		newApp.Description = row.Columns[4].GetString_()
		newApp.Weightvolume = row.Columns[5].GetString_()

		if newApp.DocumentId == documentId{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}	

//get all booking details for specified document status
func (t *CBDoc) viewDocumentsByStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document status to query")
	}

	status := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()

		if newApp.Status == status{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//get all truckIds for which status is not Received
func (t *CBDoc) viewAllTruckByDocStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document status to query")
	}

	status := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()

		if newApp.Status != status{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//get all booking details for specified document status
func (t *CBDoc) viewDetailsByDocId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document status to query")
	}

	docId := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()

		if newApp.DocumentId == docId{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//get all booking details for specified document status
func (t *CBDoc) viewDocumentTransactionHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document status to query")
	}

	documentId := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("TrxnHistory", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*TrxnHistory{}	
	
	for row := range rows {		
		newApp:= new(TrxnHistory)
		newApp.TrxId = row.Columns[0].GetString_()
		newApp.TimeStamp = row.Columns[1].GetString_()
		newApp.DocumentId = row.Columns[2].GetString_()
		newApp.UpdatedBy = row.Columns[3].GetString_()
		newApp.Status = row.Columns[4].GetString_()

		if newApp.DocumentId == documentId{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

 //get all booking details for specified document source
func (t *CBDoc) viewDocumentsBySource(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document source to query")
	}

	source := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()
		
		if newApp.Source == source{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//get all booking details for specified document destination
func (t *CBDoc) viewDocumentsByDestination(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document destination to query")
	}

	destination := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()
		
		if newApp.Destination == destination{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//get all booking details for specified document destination
func (t *CBDoc) countDocumentsByStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document destination to query")
	}

	status := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("Document", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	res2E:= []*Document{}	
	counter := Counter{0}
		for row := range rows {		
		newApp:= new(Document)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.Source = row.Columns[1].GetString_()
		newApp.Destination = row.Columns[2].GetString_()
		newApp.Status = row.Columns[3].GetString_()
		newApp.DeliveryDate = row.Columns[4].GetString_()
		newApp.DeliveryNo = row.Columns[5].GetString_()
		newApp.CustPO = row.Columns[6].GetString_()
		newApp.Documenthash = row.Columns[7].GetString_()
		newApp.ReasonCode = row.Columns[8].GetString_()
		newApp.TruckId = row.Columns[9].GetString_()
		newApp.CustomLocation = row.Columns[10].GetString_()
		newApp.SourceLatLong = row.Columns[11].GetString_()
		newApp.DestnationLatLong = row.Columns[12].GetString_()
		newApp.CustomLatLong = row.Columns[13].GetString_()

		if newApp.Status == status{
		res2E=append(res2E,newApp)
		counter.increment()
		}				
	}
	 mapB, _ := json.Marshal(strconv.Itoa(counter.currentValue()))
		
      return mapB, nil

}

//get truckId and status for a document
func (t *CBDoc) viewTruckDetailsByDocId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting zero argument")
	}
	documentId := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("ItemTracker", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*ItemTracker{}	
	
	for row := range rows {		
		newApp:= new(ItemTracker)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.TruckId = row.Columns[1].GetString_()
		newApp.LeftFactory = row.Columns[2].GetString_()
		newApp.ArrivingCustom = row.Columns[3].GetString_()
		newApp.ReachedCustom = row.Columns[4].GetString_()
		newApp.LeavingCustom = row.Columns[5].GetString_()
		newApp.ReceivedWarehouse = row.Columns[6].GetString_()
		
		if newApp.DocumentId == documentId{
		res2E=append(res2E,newApp)		
		}
	}
	    
	mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

//update document status by document id
func (t *CBDoc) updateDocumentStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5.")
	}
	documentId := args[0]
	newStatus :=  args[1]
	timeStamp :=  args[2]
	updatedBy :=  args[3]
	reasonCode := args[4]
	
	//getting TrxId incrementer
	Avalbytes, err := stub.GetState("Trx_increment")
	Aval, _ := strconv.ParseInt(string(Avalbytes), 10, 0)
	newAval:=int(Aval) + 1
	newTrx_increment:= strconv.Itoa(newAval)
	stub.PutState("Trx_increment", []byte(newTrx_increment))
		
	trxId:=string(Avalbytes)

	// Get the row pertaining to this asnNumber
		var columns1 []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: documentId}}
		columns1 = append(columns1, col1)

		row, err := stub.GetRow("Document", columns1)
		if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document with document id %s. Error %s", documentId, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
			return nil, nil
		}
		// Delete the row pertaining to this applicationId
		err = stub.DeleteRow(
			"Document",
			columns1,
		)
		if err != nil {
			return nil, errors.New("Failed deleting row.")
		}
		documentId = row.Columns[0].GetString_()
		source := row.Columns[1].GetString_()
		destination := row.Columns[2].GetString_()
		status := newStatus
		deliveryDate := row.Columns[4].GetString_()
		deliveryNo := row.Columns[5].GetString_()
		custPO := row.Columns[6].GetString_()
		documenthash := row.Columns[7].GetString_()
		truckId := row.Columns[9].GetString_()
		customLocation := row.Columns[10].GetString_()
		sourceLatLong := row.Columns[11].GetString_()
		destnationLatLong := row.Columns[12].GetString_()
		customLatLong := row.Columns[13].GetString_()
		
		// Inserting document details
		ok, err := stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: source}},
				&shim.Column{Value: &shim.Column_String_{String_: destination}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryNo}},
				&shim.Column{Value: &shim.Column_String_{String_: custPO}},
				&shim.Column{Value: &shim.Column_String_{String_: documenthash}},
				&shim.Column{Value: &shim.Column_String_{String_: reasonCode}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: customLocation}},
				&shim.Column{Value: &shim.Column_String_{String_: sourceLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: destnationLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: customLatLong}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		// Insert details in TrxnHistory table
		ok, err = stub.InsertRow("TrxnHistory", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: trxId}},
				&shim.Column{Value: &shim.Column_String_{String_: timeStamp}},
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: updatedBy}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},			
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

	return nil, nil	
}

func (t *CBDoc) updateRejectedDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	if len(args) != 14 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9.")
	}
	documentId := args[0]
	source := args[1]
	destination := args[2]
	status :=  "Initial"
	deliveryDate := args[3]
	deliveryNo := args[4]
	custPO := args[5]
	documenthash := args[6]
	timeStamp :=  args[7]
	updatedBy :=  args[8]
	reasonCode := ""
	truckId:=args[9]
	customLocation:=args[10]
	sourceLatLong:=args[11]
	destnationLatLong:=args[12]
	customLatLong:=args[13]
	//getting TrxId incrementer
	Avalbytes, err := stub.GetState("Trx_increment")
	Aval, _ := strconv.ParseInt(string(Avalbytes), 10, 0)
	newAval:=int(Aval) + 1
	newTrx_increment:= strconv.Itoa(newAval)
	stub.PutState("Trx_increment", []byte(newTrx_increment))
		
	trxId:=string(Avalbytes)

	// Get the row pertaining to this asnNumber
		var columns1 []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: documentId}}
		columns1 = append(columns1, col1)

		row, err := stub.GetRow("Document", columns1)
		if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document with document id %s. Error %s", documentId, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
			return nil, nil
		}
		// Delete the row pertaining to this applicationId
		err = stub.DeleteRow(
			"Document",
			columns1,
		)
		if err != nil {
			return nil, errors.New("Failed deleting row.")
		}
		
		// Inserting document details
		ok, err := stub.InsertRow("Document", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: source}},
				&shim.Column{Value: &shim.Column_String_{String_: destination}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryDate}},
				&shim.Column{Value: &shim.Column_String_{String_: deliveryNo}},
				&shim.Column{Value: &shim.Column_String_{String_: custPO}},
				&shim.Column{Value: &shim.Column_String_{String_: documenthash}},
				&shim.Column{Value: &shim.Column_String_{String_: reasonCode}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: customLocation}},
				&shim.Column{Value: &shim.Column_String_{String_: sourceLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: destnationLatLong}},
				&shim.Column{Value: &shim.Column_String_{String_: customLatLong}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		// Insert details in TrxnHistory table
		ok, err = stub.InsertRow("TrxnHistory", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: trxId}},
				&shim.Column{Value: &shim.Column_String_{String_: timeStamp}},
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: updatedBy}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},			
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
		
		// Get the row pertaining to this documentId
		var columns2 []shim.Column
		col1 = shim.Column{Value: &shim.Column_String_{String_: documentId}}
		columns2 = append(columns2, col1)

		row, err = stub.GetRow("ItemTracker", columns2)
		if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving document with document id %s. Error %s", documentId, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
			return nil, nil
		}
		// Delete the row pertaining to this applicationId
		err = stub.DeleteRow(
			"ItemTracker",
			columns2,
		)
		if err != nil {
			return nil, errors.New("Failed deleting row.")
		}
		leftFactory := row.Columns[2].GetString_()
		arrivingCustom := row.Columns[3].GetString_()
		reachedCustom := row.Columns[4].GetString_()
		leavingCustom := row.Columns[5].GetString_()
		receivedWarehouse := row.Columns[6].GetString_()
		
		// Insert details in ItemTracker table
		ok, err = stub.InsertRow("ItemTracker", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: leftFactory}},
				&shim.Column{Value: &shim.Column_String_{String_: arrivingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: reachedCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: leavingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: receivedWarehouse}},
			}})
	return nil, nil
}

//update ItemTracker status by truckId
func (t *CBDoc) updateTrackerStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5.")
	}
	documentIdIp := args[0]
	truckId := args[1]
	leftFactoryIp :=  args[2]
	arrivingCustomIp :=  args[3]
	reachedCustomIp :=  args[4]
	leavingCustomIp := args[5]
	receivedWarehouseIp := args[6]
	leftFactory :=  ""
	arrivingCustom :=  ""
	reachedCustom :=  ""
	leavingCustom := ""
	receivedWarehouse := ""
	
	// Get the row pertaining to this truckId
		var columns1 []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: documentIdIp}}
		columns1 = append(columns1, col1)

		row, err := stub.GetRow("ItemTracker", columns1)
		if err != nil {
			return nil, fmt.Errorf("Error: Failed retrieving records with truck id %s. Error %s", documentIdIp, err.Error())
		}

		// GetRows returns empty message if key does not exist
		if len(row.Columns) == 0 {
			return nil, nil
		}
		// Delete the row pertaining to this truck
		err = stub.DeleteRow(
			"ItemTracker",
			columns1,
		)
		if err != nil {
			return nil, errors.New("Failed deleting row.")
		}
		documentId := row.Columns[0].GetString_()
		truckId = row.Columns[1].GetString_()
		if leftFactoryIp == ""{
			leftFactory = row.Columns[2].GetString_()
		} else {
			leftFactory = leftFactoryIp
		}
		if arrivingCustomIp == ""{
			arrivingCustom = row.Columns[3].GetString_()
		} else {
			arrivingCustom = arrivingCustomIp
		}
		if reachedCustomIp == ""{
			reachedCustom = row.Columns[4].GetString_()
		} else {
			reachedCustom = reachedCustomIp
		}
		if leavingCustomIp == ""{
			leavingCustom = row.Columns[5].GetString_()
		} else {
			leavingCustom = leavingCustomIp
		}
		if receivedWarehouseIp == ""{
			receivedWarehouse = row.Columns[6].GetString_()
		} else {
			receivedWarehouse = receivedWarehouseIp
		}
		// Inserting ItemTracker details
		ok, err := stub.InsertRow("ItemTracker", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: documentId}},
				&shim.Column{Value: &shim.Column_String_{String_: truckId}},
				&shim.Column{Value: &shim.Column_String_{String_: leftFactory}},
				&shim.Column{Value: &shim.Column_String_{String_: arrivingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: reachedCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: leavingCustom}},
				&shim.Column{Value: &shim.Column_String_{String_: receivedWarehouse}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

	return nil, nil	
}

func (t *CBDoc) viewAllHighPriorityDocuments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting document destination to query")
	}

	arrivingCustom := args[0]
	
	var columns []shim.Column

	rows, err := stub.GetRows("ItemTracker", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
			
	res2E:= []*Document{}	
	
	for row := range rows {		
		newApp:= new(ItemTracker)
		newApp.DocumentId = row.Columns[0].GetString_()
		newApp.TruckId = row.Columns[1].GetString_()
		newApp.LeftFactory = row.Columns[2].GetString_()
		newApp.ArrivingCustom = row.Columns[3].GetString_()
		newApp.ReachedCustom = row.Columns[4].GetString_()
		newApp.LeavingCustom = row.Columns[5].GetString_()
		newApp.ReceivedWarehouse = row.Columns[6].GetString_()
		
		if newApp.ArrivingCustom == arrivingCustom{
			
			var columns1 []shim.Column

			rows1, err1 := stub.GetRows("Document", columns1)
			if err1 != nil {
				return nil, fmt.Errorf("Failed to retrieve row")
			}
			for row1 := range rows1 {
				newApp1:= new(Document)
				newApp1.DocumentId = row1.Columns[0].GetString_()
				newApp1.Source = row1.Columns[1].GetString_()
				newApp1.Destination = row1.Columns[2].GetString_()
				newApp1.Status = row1.Columns[3].GetString_()
				newApp1.DeliveryDate = row1.Columns[4].GetString_()
				newApp1.DeliveryNo = row1.Columns[5].GetString_()
				newApp1.CustPO = row1.Columns[6].GetString_()
				newApp1.Documenthash = row1.Columns[7].GetString_()
				newApp1.ReasonCode = row1.Columns[8].GetString_()
				newApp1.TruckId = row1.Columns[9].GetString_()
				newApp1.CustomLocation = row1.Columns[10].GetString_()
				newApp1.SourceLatLong = row1.Columns[11].GetString_()
				newApp1.DestnationLatLong = row1.Columns[12].GetString_()
				newApp1.CustomLatLong = row1.Columns[13].GetString_()
				
				if newApp1.DocumentId == newApp.DocumentId && newApp1.Status != "Received"{
					res2E=append(res2E,newApp1)		
				}
			}
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}

// Invoke invokes the chaincode
func (t *CBDoc) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createDocument" {
		t := CBDoc{}
		return t.createDocument(stub, args)	
	} 
	
	if function == "updateDocumentStatus" {
		t := CBDoc{}
		return t.updateDocumentStatus(stub, args)	
	}
	
	if function == "insertItemList" {
		t := CBDoc{}
		return t.insertItemList(stub, args)	
	}
	
	if function == "updateRejectedDocument" {
		t := CBDoc{}
		return t.updateRejectedDocument(stub, args)	
	}
	
	if function == "updateTrackerStatus" {
		t := CBDoc{}
		return t.updateTrackerStatus(stub, args)	
	}

	return nil, errors.New("Invalid invoke function name.")
	
}

// query queries the chaincode
func (t *CBDoc) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "viewDocumentsByStatus" {
		t := CBDoc{}
		return t.viewDocumentsByStatus(stub, args)		
	} else if function == "viewDocumentsBySource" { 
		t := CBDoc{}
		return t.viewDocumentsBySource(stub, args)
	}else if function == "viewDocumentsByDestination" { 
		t := CBDoc{}
		return t.viewDocumentsByDestination(stub, args)
	}else if function == "countDocumentsByStatus" { 
		t := CBDoc{}
		return t.countDocumentsByStatus(stub, args)
	}else if function == "viewDocumentTransactionHistory" { 
		t := CBDoc{}
		return t.viewDocumentTransactionHistory(stub, args)
	}else if function == "viewDetailsByDocId" { 
		t := CBDoc{}
		return t.viewDetailsByDocId(stub, args)
	}else if function == "viewItemListByDocumentId" { 
		t := CBDoc{}
		return t.viewItemListByDocumentId(stub, args)
	}else if function == "viewTruckDetailsByDocId" { 
		t := CBDoc{}
		return t.viewTruckDetailsByDocId(stub, args)
	}else if function == "viewAllTruckByDocStatus" { 
		t := CBDoc{}
		return t.viewAllTruckByDocStatus(stub, args)
	}else if function == "viewAllHighPriorityDocuments" { 
		t := CBDoc{}
		return t.viewAllHighPriorityDocuments(stub, args)
	}	
		
	return nil, errors.New("Invalid query function name.")
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CBDoc))
	if err != nil {
		fmt.Printf("Error starting CBDoc: %s", err)
	}
} 