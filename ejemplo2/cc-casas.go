/*
 Ejemplo simple de transacciones en blockchain.
 
 Hay tres usuarios registrados y dos casas.
 Cada casa tiene un propietario y un valor.
 Para simplificar el ejemplo, en la chaincode (cc) solo se almacena la casa (key) y el propietario (value).
 El valor de la casa1000 es de 1000 y el de la casa2000 de 2000
 
 Los usuarios se almacenan por su nombre (key) y un saldo (value)
 
 Una casa puede cambiar de dueno si la vende el propietario y el comprador tiene
 el dinero suficiente.
 
 
 Funciones:
 Init - para inicializar valores (transacciones) en la cadena
 Compra - Intento de comprar la casa
 AgregaSaldo - Incrementa el dinero disponible en la cuenta de un usuario
 ConsultaPropietario - Devuelve quien es el propietario de la casa
 ConsultaSaldo - Devuelve el saldo del usuario
 */


package main

import (
	"errors"
	"fmt"
    "strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
   )

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    var err error
    err = stub.PutState("casa1000",[]byte("luis"))
    if err != nil {
        return nil, err
    }
    err = stub.PutState("casa2000",[]byte("maria"))
    if err != nil {
        return nil, err
    }
    err = stub.PutState("luis",[]byte(strconv.Itoa(100)))
    if err != nil {
        return nil, err
    }
    err = stub.PutState("maria",[]byte(strconv.Itoa(200)))
    if err != nil {
        return nil, err
    }
    err = stub.PutState("juan",[]byte(strconv.Itoa(300)))
    if err != nil {
        return nil, err
    }
    
    return nil, nil
}
// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
    } else if function == "agregasaldo" {
        return t.AgregaSaldo(stub, args)
    } else if function == "compra" {
        return t.Compra(stub, args)
    }
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "consultapropietario" {		   //manda llamar la función
        return t.ConsultaPropietario(stub, args)
    } else if function == "consultasaldo" {		   //manda llamar la función
        return t.ConsultaSaldo(stub, args)
    }
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

/*func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running write()")
    
    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
    }
    
    key = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}*/

// Compra recibe tres argumentos: la casa, el vendedor y el comprador.
// Implementa la regla de negocio y devuelve error si no se pudo hacer la compra.
func (t *SimpleChaincode) Compra(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    
    var casa, duenio, comprador string
    //var duenio, vendedor, comprador string
    var valCasa int
    var saldoComprador int
    var err error
    var jsonResp string
    
    fmt.Println("en Compra( )")
    valCasa = 0
    if len(args) != 3 {
        return nil, errors.New("Numero incorrecto de argumentos. Necesito 3: casa, vendedor y comprador")
    }

    casa = args[0]
    if casa == "casa1000" {
       valCasa = 1000
    } else if casa == "casa2000" {
        valCasa = 2000
    } else {
        return nil, errors.New("Error. La casa " + casa + "no existe en la chaincode")
    }
    valAsbytes, err := stub.GetState(casa)
    duenio = string(valAsbytes[:])
    if err != nil {
        jsonResp = "{\"Error\":\"No pude obtener el propietario de " + casa + "\"}"
        return nil, errors.New(jsonResp)
    }
    
    fmt.Println("la casa " + casa + "tiene un valor " + duenio)
    
    if duenio != args[1] {
        return nil, errors.New("Error. El valor recibido no es el propietario de la casa")
    }
    
    comprador = args[2]
    valAsbytes, err = stub.GetState(comprador)
    saldoComprador, _  = strconv.Atoi(string(valAsbytes))
    
    if saldoComprador < valCasa {
        return nil, errors.New("Error. El comprador no tiene saldo suficiente")
    }
    
    err = stub.PutState(casa,[]byte(comprador))
    if err != nil {
        return nil, err
    }
    err = stub.PutState(comprador, []byte(strconv.Itoa(saldoComprador-valCasa)))
    if err != nil {
        return nil, err
    }

    
    /*
     Lee casa como primer argumento y haz GetState para conocer a su propietario
     Asigna costo casa en función del nombre de casa
     
     Lee propietario en segundo argumento y compara con el propietario de la casa
        Si los propietarios son distintos, regresa con un mensaje de que no se puede vender
     Lee comprador en el tercer argumento y toma su saldo
       Si saldo es menor que costo de la casa, regresa con mensaje de que no se puede comprar
     Asigna nuevo propietario a casa y guarda en cc
     Descuenta del saldo el monto de la casa y guarda en cc
     */
    return nil, nil
}

// AgregaSaldo recibe un argumento: El usuario.
func (t *SimpleChaincode) AgregaSaldo(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
        var nombreUsuario  string
        var cantidad int
        var err error
        fmt.Println("en agregasaldo( )")
        
        if len(args) != 2 {
            return nil, errors.New("Numero incorrecto de argumentos. Necesito 2: llave y valor")
        }
        
        // En un codigo mas formal, seria necesario validar que el usuario exista
        nombreUsuario = args[0]
        cantidad, err  = strconv.Atoi(args[1])  // Convierte cantidad en valor entero
   
    if err != nil {
        return nil, errors.New("2o argumento debe ser numeric string")
    }
        
        err = stub.PutState(nombreUsuario, []byte(strconv.Itoa(cantidad)))  // Escribe y cambia el cc state
        if err != nil {
            return nil, err
        }
        return nil, nil
    }

// Consulta propietario recibe un argumento: La casa.
func (t *SimpleChaincode) ConsultaPropietario(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error
    
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
    
    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}

// ConsultaSaldo recibe un argumento: El usuario.
func (t *SimpleChaincode) ConsultaSaldo(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error
    
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
    
    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}

