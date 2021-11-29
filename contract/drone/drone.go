package main



import (

	"bytes"

	"encoding/json"

	"fmt"

	"strconv"



	"github.com/hyperledger/fabric/core/chaincode/shim"

	sc "github.com/hyperledger/fabric/protos/peer"

)



type SmartContract struct {

}



type Drone struct{

	DroneID string `json:"dronid"`

	Make string `json:"make"`

	Model string `json:"model"`

	Owner string `json:"owner"`

	Datas []Data `json:"datas"`

}

type Data struct{

	DroneID string  `json:"droneid"`

	Altitude string `json:"고도"`

	Angle string `json:"각도"`

	Height string `json:"무게"`

	Speed string `json:"속도"`
	
	Latitude string `json:"위도"`

	Longitude string `json:"경도"`

}



func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	return shim.Success(nil)

}



func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {



	function, args := APIstub.GetFunctionAndParameters()



	if function == "queryDrone" {

		return s.queryDrone(APIstub, args)

	} else if function == "initLedger" {

		return s.initLedger(APIstub)

	} else if function == "createDrone" {

		return s.createDrone(APIstub, args)

	} else if function == "queryAllDrones" {

		return s.queryAllDrones(APIstub)

	} else if function == "changeDroneOwner" {

		return s.changeDroneOwner(APIstub, args)

	} else if function == "addRating" {

		return s.addRating(APIstub, args)

	}

	return shim.Error("Invalid Smart Contract function name.")

}





func (s *SmartContract) queryDrone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {



	if len(args) != 1 {

		return shim.Error("Incorrect number of arguments. Expecting 1")

	}

	droneAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(droneAsBytes)

}



func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	drones := []Drone{

		Drone{DroneID: "DRONE000", Make: "Airbus", Model: "NextGen", Owner: "Ajou Univ"},

		Drone{DroneID: "DRONE001", Make: "aeroG", Model: "aeroG Aviation aG-4 Liberty", Owner: "Samsung Co"},

		Drone{DroneID: "DRONE002", Make: "Hyundai", Model: "Cargo UAS2", Owner: "Hyundai Co"},

		Drone{DroneID: "DRONE003", Make: "Uber", Model: "eCRM-002", Owner: "Uber TAXI Co"},

		Drone{DroneID: "DRONE004", Make: "Bell", Model: "Nexus 4EX", Owner: "Ajou Univ"},

		Drone{DroneID: "DRONE005", Make: "Hyundai", Model: "UAM SA1", Owner: "Hyundai Co"},

		Drone{DroneID: "DRONE006", Make: "Bell", Model: "Nexus 6EX", Owner: "Ajou Univ"},

		Drone{DroneID: "DRONE007", Make: "Hyundai", Model: "UAM SA2", Owner: "GOSUOHWAN"},

		Drone{DroneID: "DRONE008", Make: "Airbus", Model: "CityAirbus", Owner: "Ajou Univ"},

		Drone{DroneID: "DRONE009", Make: "aeroG", Model: "aeroG Aviation aG-4 Liberty", Owner: "Samsung Co"},

		Drone{DroneID: "DRONE0010", Make: "HyundaiMobis", Model: "Cargo UAS", Owner: "SooYeonida"},

		Drone{DroneID: "DRONE0011", Make: "Hyundai", Model: "UAM SA3", Owner: "Hyundai Co"},

		Drone{DroneID: "DRONE0012", Make: "Dusan Stan", Model: "Aliptera APV-1", Owner: "Subin"},

		Drone{DroneID: "DRONE0013", Make: "Porsche Boeing", Model: "Porsche", Owner: "Gunyeong"},

		}

		i := 0

	for i < len(drones) {

		fmt.Println("i is ", i)

		droneAsBytes, _ := json.Marshal(drones[i])

		APIstub.PutState("DRONE"+strconv.Itoa(i), droneAsBytes)

		fmt.Println("Added", drones[i])

		i = i + 1

	}



	return shim.Success(nil)

}



func (s *SmartContract) createDrone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {



	if len(args) != 4 {

		return shim.Error("Incorrect number of arguments. Expecting 4")

	}

	var drone = Drone{DroneID: args[0], Make: args[1], Model: args[2], Owner: args[3]}



	droneAsBytes, _ := json.Marshal(drone)

	APIstub.PutState(args[0], droneAsBytes)



	return shim.Success(nil)

}



func (s *SmartContract) queryAllDrones(APIstub shim.ChaincodeStubInterface) sc.Response {



	startKey := ""

	endKey := ""



	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)



	if err != nil {

		return shim.Error(err.Error())

	}



	defer resultsIterator.Close()


	var buffer bytes.Buffer

	buffer.WriteString("[")



	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()

		if err != nil {

			return shim.Error(err.Error())

		}


		if bArrayMemberAlreadyWritten == true {

			buffer.WriteString(",")

		}

		buffer.WriteString("{\"Key\":")

		buffer.WriteString("\"")

		buffer.WriteString(queryResponse.Key)

		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")



	fmt.Printf("- queryAllDrones:\n%s\n", buffer.String())



	return shim.Success(buffer.Bytes())

}



func (s *SmartContract) changeDroneOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {



	if len(args) != 2 {

		return shim.Error("Incorrect number of arguments. Expecting 2")

	}



	droneAsBytes, _ := APIstub.GetState(args[0])

	drone := Drone{}



	json.Unmarshal(droneAsBytes, &drone)

	drone.Owner = args[1]



	droneAsBytes, _ = json.Marshal(drone)

	APIstub.PutState(args[0], droneAsBytes)



	return shim.Success(nil)

}



func (s *SmartContract) addRating(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {

		return shim.Error("Incorrect number of arguments. Expecting 7")

	}

	droneAsBytes, err := APIstub.GetState(args[0])

	if err != nil{

		jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"

		return shim.Error(jsonResp)

	} else if droneAsBytes == nil{ 

		jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"

		return shim.Error(jsonResp)

	}


	drone := Drone{}

	err = json.Unmarshal(droneAsBytes, &drone)

	if err != nil {

		return shim.Error(err.Error())

	}

	

	var Data = Data{DroneID: args[0] , Altitude: args[1], Angle: args[2], Height: args[3], Speed: args[4], Latitude: args[5],
	Longitude: args[6], }

	drone.Datas=append(drone.Datas,Data)
	
	droneAsBytes, err = json.Marshal(drone);

	APIstub.PutState(args[0], droneAsBytes)

	return shim.Success([]byte("data is updated"))
}

func main() {

	err := shim.Start(new(SmartContract))

	if err != nil {

		fmt.Printf("Error creating new Smart Contract: %s", err)

	}

}
