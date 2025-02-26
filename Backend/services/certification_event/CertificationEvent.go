// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package certification

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CertificationEventCertEvent is an auto generated low-level Go binding around an user-defined struct.
type CertificationEventCertEvent struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	IsActive         bool
}

// CertificationMetaData contains all meta data concerning the Certification contract.
var CertificationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistryAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"CertificationEventDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"name\":\"CertificationEventStored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"certificationEvents\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"}],\"name\":\"storeCertificationEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"}],\"name\":\"deactivateCertificationEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"}],\"name\":\"getActiveCertificationsForEntity\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"internalType\":\"structCertificationEvent.CertEvent[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"getAllCertifications\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"internalType\":\"structCertificationEvent.CertEvent[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"}],\"name\":\"isCertificationCIDExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
}

// CertificationABI is the input ABI used to generate the binding from.
// Deprecated: Use CertificationMetaData.ABI instead.
var CertificationABI = CertificationMetaData.ABI

// Certification is an auto generated Go binding around an Ethereum contract.
type Certification struct {
	CertificationCaller     // Read-only binding to the contract
	CertificationTransactor // Write-only binding to the contract
	CertificationFilterer   // Log filterer for contract events
}

// CertificationCaller is an auto generated read-only Go binding around an Ethereum contract.
type CertificationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CertificationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CertificationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CertificationSession struct {
	Contract     *Certification    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CertificationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CertificationCallerSession struct {
	Contract *CertificationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// CertificationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CertificationTransactorSession struct {
	Contract     *CertificationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CertificationRaw is an auto generated low-level Go binding around an Ethereum contract.
type CertificationRaw struct {
	Contract *Certification // Generic contract binding to access the raw methods on
}

// CertificationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CertificationCallerRaw struct {
	Contract *CertificationCaller // Generic read-only contract binding to access the raw methods on
}

// CertificationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CertificationTransactorRaw struct {
	Contract *CertificationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCertification creates a new instance of Certification, bound to a specific deployed contract.
func NewCertification(address common.Address, backend bind.ContractBackend) (*Certification, error) {
	contract, err := bindCertification(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Certification{CertificationCaller: CertificationCaller{contract: contract}, CertificationTransactor: CertificationTransactor{contract: contract}, CertificationFilterer: CertificationFilterer{contract: contract}}, nil
}

// NewCertificationCaller creates a new read-only instance of Certification, bound to a specific deployed contract.
func NewCertificationCaller(address common.Address, caller bind.ContractCaller) (*CertificationCaller, error) {
	contract, err := bindCertification(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CertificationCaller{contract: contract}, nil
}

// NewCertificationTransactor creates a new write-only instance of Certification, bound to a specific deployed contract.
func NewCertificationTransactor(address common.Address, transactor bind.ContractTransactor) (*CertificationTransactor, error) {
	contract, err := bindCertification(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CertificationTransactor{contract: contract}, nil
}

// NewCertificationFilterer creates a new log filterer instance of Certification, bound to a specific deployed contract.
func NewCertificationFilterer(address common.Address, filterer bind.ContractFilterer) (*CertificationFilterer, error) {
	contract, err := bindCertification(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CertificationFilterer{contract: contract}, nil
}

// bindCertification binds a generic wrapper to an already deployed contract.
func bindCertification(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CertificationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Certification *CertificationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Certification.Contract.CertificationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Certification *CertificationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Certification.Contract.CertificationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Certification *CertificationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Certification.Contract.CertificationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Certification *CertificationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Certification.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Certification *CertificationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Certification.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Certification *CertificationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Certification.Contract.contract.Transact(opts, method, params...)
}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationCaller) CertificationEvents(opts *bind.CallOpts, arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	IsActive         bool
}, error) {
	var out []interface{}
	err := _Certification.contract.Call(opts, &out, "certificationEvents", arg0)

	outstruct := new(struct {
		EventID          string
		EntityType       string
		EntityID         string
		CertificationCID string
		IssuedDate       *big.Int
		ExpiryDate       *big.Int
		CreatedOn        *big.Int
		IsActive         bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EventID = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.EntityType = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.EntityID = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.CertificationCID = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.IssuedDate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.ExpiryDate = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.CreatedOn = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.IsActive = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationSession) CertificationEvents(arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	IsActive         bool
}, error) {
	return _Certification.Contract.CertificationEvents(&_Certification.CallOpts, arg0)
}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationCallerSession) CertificationEvents(arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	IsActive         bool
}, error) {
	return _Certification.Contract.CertificationEvents(&_Certification.CallOpts, arg0)
}

// GetActiveCertificationsForEntity is a free data retrieval call binding the contract method 0xf42e4e56.
//
// Solidity: function getActiveCertificationsForEntity(string entityID) view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationCaller) GetActiveCertificationsForEntity(opts *bind.CallOpts, entityID string) ([]CertificationEventCertEvent, error) {
	var out []interface{}
	err := _Certification.contract.Call(opts, &out, "getActiveCertificationsForEntity", entityID)

	if err != nil {
		return *new([]CertificationEventCertEvent), err
	}

	out0 := *abi.ConvertType(out[0], new([]CertificationEventCertEvent)).(*[]CertificationEventCertEvent)

	return out0, err

}

// GetActiveCertificationsForEntity is a free data retrieval call binding the contract method 0xf42e4e56.
//
// Solidity: function getActiveCertificationsForEntity(string entityID) view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationSession) GetActiveCertificationsForEntity(entityID string) ([]CertificationEventCertEvent, error) {
	return _Certification.Contract.GetActiveCertificationsForEntity(&_Certification.CallOpts, entityID)
}

// GetActiveCertificationsForEntity is a free data retrieval call binding the contract method 0xf42e4e56.
//
// Solidity: function getActiveCertificationsForEntity(string entityID) view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationCallerSession) GetActiveCertificationsForEntity(entityID string) ([]CertificationEventCertEvent, error) {
	return _Certification.Contract.GetActiveCertificationsForEntity(&_Certification.CallOpts, entityID)
}

// GetAllCertifications is a free data retrieval call binding the contract method 0xb46dc5cb.
//
// Solidity: function getAllCertifications() view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationCaller) GetAllCertifications(opts *bind.CallOpts) ([]CertificationEventCertEvent, error) {
	var out []interface{}
	err := _Certification.contract.Call(opts, &out, "getAllCertifications")

	if err != nil {
		return *new([]CertificationEventCertEvent), err
	}

	out0 := *abi.ConvertType(out[0], new([]CertificationEventCertEvent)).(*[]CertificationEventCertEvent)

	return out0, err

}

// GetAllCertifications is a free data retrieval call binding the contract method 0xb46dc5cb.
//
// Solidity: function getAllCertifications() view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationSession) GetAllCertifications() ([]CertificationEventCertEvent, error) {
	return _Certification.Contract.GetAllCertifications(&_Certification.CallOpts)
}

// GetAllCertifications is a free data retrieval call binding the contract method 0xb46dc5cb.
//
// Solidity: function getAllCertifications() view returns((string,string,string,string,uint256,uint256,uint256,bool)[])
func (_Certification *CertificationCallerSession) GetAllCertifications() ([]CertificationEventCertEvent, error) {
	return _Certification.Contract.GetAllCertifications(&_Certification.CallOpts)
}

// IsCertificationCIDExists is a free data retrieval call binding the contract method 0xe8dab040.
//
// Solidity: function isCertificationCIDExists(string certificationCID) view returns(bool)
func (_Certification *CertificationCaller) IsCertificationCIDExists(opts *bind.CallOpts, certificationCID string) (bool, error) {
	var out []interface{}
	err := _Certification.contract.Call(opts, &out, "isCertificationCIDExists", certificationCID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsCertificationCIDExists is a free data retrieval call binding the contract method 0xe8dab040.
//
// Solidity: function isCertificationCIDExists(string certificationCID) view returns(bool)
func (_Certification *CertificationSession) IsCertificationCIDExists(certificationCID string) (bool, error) {
	return _Certification.Contract.IsCertificationCIDExists(&_Certification.CallOpts, certificationCID)
}

// IsCertificationCIDExists is a free data retrieval call binding the contract method 0xe8dab040.
//
// Solidity: function isCertificationCIDExists(string certificationCID) view returns(bool)
func (_Certification *CertificationCallerSession) IsCertificationCIDExists(certificationCID string) (bool, error) {
	return _Certification.Contract.IsCertificationCIDExists(&_Certification.CallOpts, certificationCID)
}

// DeactivateCertificationEvent is a paid mutator transaction binding the contract method 0xbac9ba98.
//
// Solidity: function deactivateCertificationEvent(string eventID) returns()
func (_Certification *CertificationTransactor) DeactivateCertificationEvent(opts *bind.TransactOpts, eventID string) (*types.Transaction, error) {
	return _Certification.contract.Transact(opts, "deactivateCertificationEvent", eventID)
}

// DeactivateCertificationEvent is a paid mutator transaction binding the contract method 0xbac9ba98.
//
// Solidity: function deactivateCertificationEvent(string eventID) returns()
func (_Certification *CertificationSession) DeactivateCertificationEvent(eventID string) (*types.Transaction, error) {
	return _Certification.Contract.DeactivateCertificationEvent(&_Certification.TransactOpts, eventID)
}

// DeactivateCertificationEvent is a paid mutator transaction binding the contract method 0xbac9ba98.
//
// Solidity: function deactivateCertificationEvent(string eventID) returns()
func (_Certification *CertificationTransactorSession) DeactivateCertificationEvent(eventID string) (*types.Transaction, error) {
	return _Certification.Contract.DeactivateCertificationEvent(&_Certification.TransactOpts, eventID)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_Certification *CertificationTransactor) StoreCertificationEvent(opts *bind.TransactOpts, eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _Certification.contract.Transact(opts, "storeCertificationEvent", eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_Certification *CertificationSession) StoreCertificationEvent(eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _Certification.Contract.StoreCertificationEvent(&_Certification.TransactOpts, eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_Certification *CertificationTransactorSession) StoreCertificationEvent(eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _Certification.Contract.StoreCertificationEvent(&_Certification.TransactOpts, eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// CertificationCertificationEventDeactivatedIterator is returned from FilterCertificationEventDeactivated and is used to iterate over the raw logs and unpacked data for CertificationEventDeactivated events raised by the Certification contract.
type CertificationCertificationEventDeactivatedIterator struct {
	Event *CertificationCertificationEventDeactivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CertificationCertificationEventDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CertificationCertificationEventDeactivated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CertificationCertificationEventDeactivated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CertificationCertificationEventDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CertificationCertificationEventDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CertificationCertificationEventDeactivated represents a CertificationEventDeactivated event raised by the Certification contract.
type CertificationCertificationEventDeactivated struct {
	EventID   string
	EntityID  string
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCertificationEventDeactivated is a free log retrieval operation binding the contract event 0x98b61e6ccb902eb502f90b318e95553cfec1b91ba40a570f7fe611016cdccca0.
//
// Solidity: event CertificationEventDeactivated(string eventID, string entityID, uint256 timestamp)
func (_Certification *CertificationFilterer) FilterCertificationEventDeactivated(opts *bind.FilterOpts) (*CertificationCertificationEventDeactivatedIterator, error) {

	logs, sub, err := _Certification.contract.FilterLogs(opts, "CertificationEventDeactivated")
	if err != nil {
		return nil, err
	}
	return &CertificationCertificationEventDeactivatedIterator{contract: _Certification.contract, event: "CertificationEventDeactivated", logs: logs, sub: sub}, nil
}

// WatchCertificationEventDeactivated is a free log subscription operation binding the contract event 0x98b61e6ccb902eb502f90b318e95553cfec1b91ba40a570f7fe611016cdccca0.
//
// Solidity: event CertificationEventDeactivated(string eventID, string entityID, uint256 timestamp)
func (_Certification *CertificationFilterer) WatchCertificationEventDeactivated(opts *bind.WatchOpts, sink chan<- *CertificationCertificationEventDeactivated) (event.Subscription, error) {

	logs, sub, err := _Certification.contract.WatchLogs(opts, "CertificationEventDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CertificationCertificationEventDeactivated)
				if err := _Certification.contract.UnpackLog(event, "CertificationEventDeactivated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCertificationEventDeactivated is a log parse operation binding the contract event 0x98b61e6ccb902eb502f90b318e95553cfec1b91ba40a570f7fe611016cdccca0.
//
// Solidity: event CertificationEventDeactivated(string eventID, string entityID, uint256 timestamp)
func (_Certification *CertificationFilterer) ParseCertificationEventDeactivated(log types.Log) (*CertificationCertificationEventDeactivated, error) {
	event := new(CertificationCertificationEventDeactivated)
	if err := _Certification.contract.UnpackLog(event, "CertificationEventDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CertificationCertificationEventStoredIterator is returned from FilterCertificationEventStored and is used to iterate over the raw logs and unpacked data for CertificationEventStored events raised by the Certification contract.
type CertificationCertificationEventStoredIterator struct {
	Event *CertificationCertificationEventStored // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CertificationCertificationEventStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CertificationCertificationEventStored)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CertificationCertificationEventStored)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CertificationCertificationEventStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CertificationCertificationEventStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CertificationCertificationEventStored represents a CertificationEventStored event raised by the Certification contract.
type CertificationCertificationEventStored struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	IsActive         bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCertificationEventStored is a free log retrieval operation binding the contract event 0x3105924db4236c20ebd33b32e084c8041fac64467f17b95da82af0358bc58d5a.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationFilterer) FilterCertificationEventStored(opts *bind.FilterOpts) (*CertificationCertificationEventStoredIterator, error) {

	logs, sub, err := _Certification.contract.FilterLogs(opts, "CertificationEventStored")
	if err != nil {
		return nil, err
	}
	return &CertificationCertificationEventStoredIterator{contract: _Certification.contract, event: "CertificationEventStored", logs: logs, sub: sub}, nil
}

// WatchCertificationEventStored is a free log subscription operation binding the contract event 0x3105924db4236c20ebd33b32e084c8041fac64467f17b95da82af0358bc58d5a.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationFilterer) WatchCertificationEventStored(opts *bind.WatchOpts, sink chan<- *CertificationCertificationEventStored) (event.Subscription, error) {

	logs, sub, err := _Certification.contract.WatchLogs(opts, "CertificationEventStored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CertificationCertificationEventStored)
				if err := _Certification.contract.UnpackLog(event, "CertificationEventStored", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCertificationEventStored is a log parse operation binding the contract event 0x3105924db4236c20ebd33b32e084c8041fac64467f17b95da82af0358bc58d5a.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn, bool isActive)
func (_Certification *CertificationFilterer) ParseCertificationEventStored(log types.Log) (*CertificationCertificationEventStored, error) {
	event := new(CertificationCertificationEventStored)
	if err := _Certification.contract.UnpackLog(event, "CertificationEventStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
