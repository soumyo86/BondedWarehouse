
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
	// setting up the increment for Transaction Id
	stub.PutState("Trx_increment", []byte("1"))
	stub.PutState("Item_increment", []byte("1"))
	
	return nil, nil
}

// generate booking number for shipping item
func (t *CBDoc) createDocument(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 9 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 6. Got: %d.", len(args))
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

		if newApp.Status == status{
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

		if newApp.Status == status{
		res2E=append(res2E,newApp)
		counter.increment()
		}				
	}
	 mapB, _ := json.Marshal(strconv.Itoa(counter.currentValue()))
		
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