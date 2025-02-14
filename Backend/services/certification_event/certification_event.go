// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package certification_event

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

// CertificationEventMetaData contains all meta data concerning the CertificationEvent contract.
var CertificationEventMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"}],\"name\":\"CertificationEventStored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"certificationEvents\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdOn\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"entityID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"certificationCID\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"issuedDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"}],\"name\":\"storeCertificationEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"eventID\",\"type\":\"string\"}],\"name\":\"getCertificationEvent\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
}

// CertificationEventABI is the input ABI used to generate the binding from.
// Deprecated: Use CertificationEventMetaData.ABI instead.
var CertificationEventABI = CertificationEventMetaData.ABI

// CertificationEvent is an auto generated Go binding around an Ethereum contract.
type CertificationEvent struct {
	CertificationEventCaller     // Read-only binding to the contract
	CertificationEventTransactor // Write-only binding to the contract
	CertificationEventFilterer   // Log filterer for contract events
}

// CertificationEventCaller is an auto generated read-only Go binding around an Ethereum contract.
type CertificationEventCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationEventTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CertificationEventTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationEventFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CertificationEventFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificationEventSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CertificationEventSession struct {
	Contract     *CertificationEvent // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CertificationEventCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CertificationEventCallerSession struct {
	Contract *CertificationEventCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// CertificationEventTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CertificationEventTransactorSession struct {
	Contract     *CertificationEventTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// CertificationEventRaw is an auto generated low-level Go binding around an Ethereum contract.
type CertificationEventRaw struct {
	Contract *CertificationEvent // Generic contract binding to access the raw methods on
}

// CertificationEventCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CertificationEventCallerRaw struct {
	Contract *CertificationEventCaller // Generic read-only contract binding to access the raw methods on
}

// CertificationEventTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CertificationEventTransactorRaw struct {
	Contract *CertificationEventTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCertificationEvent creates a new instance of CertificationEvent, bound to a specific deployed contract.
func NewCertificationEvent(address common.Address, backend bind.ContractBackend) (*CertificationEvent, error) {
	contract, err := bindCertificationEvent(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CertificationEvent{CertificationEventCaller: CertificationEventCaller{contract: contract}, CertificationEventTransactor: CertificationEventTransactor{contract: contract}, CertificationEventFilterer: CertificationEventFilterer{contract: contract}}, nil
}

// NewCertificationEventCaller creates a new read-only instance of CertificationEvent, bound to a specific deployed contract.
func NewCertificationEventCaller(address common.Address, caller bind.ContractCaller) (*CertificationEventCaller, error) {
	contract, err := bindCertificationEvent(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CertificationEventCaller{contract: contract}, nil
}

// NewCertificationEventTransactor creates a new write-only instance of CertificationEvent, bound to a specific deployed contract.
func NewCertificationEventTransactor(address common.Address, transactor bind.ContractTransactor) (*CertificationEventTransactor, error) {
	contract, err := bindCertificationEvent(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CertificationEventTransactor{contract: contract}, nil
}

// NewCertificationEventFilterer creates a new log filterer instance of CertificationEvent, bound to a specific deployed contract.
func NewCertificationEventFilterer(address common.Address, filterer bind.ContractFilterer) (*CertificationEventFilterer, error) {
	contract, err := bindCertificationEvent(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CertificationEventFilterer{contract: contract}, nil
}

// bindCertificationEvent binds a generic wrapper to an already deployed contract.
func bindCertificationEvent(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CertificationEventMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CertificationEvent *CertificationEventRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CertificationEvent.Contract.CertificationEventCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CertificationEvent *CertificationEventRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CertificationEvent.Contract.CertificationEventTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CertificationEvent *CertificationEventRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CertificationEvent.Contract.CertificationEventTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CertificationEvent *CertificationEventCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CertificationEvent.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CertificationEvent *CertificationEventTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CertificationEvent.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CertificationEvent *CertificationEventTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CertificationEvent.Contract.contract.Transact(opts, method, params...)
}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventCaller) CertificationEvents(opts *bind.CallOpts, arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
}, error) {
	var out []interface{}
	err := _CertificationEvent.contract.Call(opts, &out, "certificationEvents", arg0)

	outstruct := new(struct {
		EventID          string
		EntityType       string
		EntityID         string
		CertificationCID string
		IssuedDate       *big.Int
		ExpiryDate       *big.Int
		CreatedOn        *big.Int
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

	return *outstruct, err

}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventSession) CertificationEvents(arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
}, error) {
	return _CertificationEvent.Contract.CertificationEvents(&_CertificationEvent.CallOpts, arg0)
}

// CertificationEvents is a free data retrieval call binding the contract method 0x379ca547.
//
// Solidity: function certificationEvents(string ) view returns(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventCallerSession) CertificationEvents(arg0 string) (struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
}, error) {
	return _CertificationEvent.Contract.CertificationEvents(&_CertificationEvent.CallOpts, arg0)
}

// GetCertificationEvent is a free data retrieval call binding the contract method 0x3e970613.
//
// Solidity: function getCertificationEvent(string eventID) view returns(string, string, string, string, uint256, uint256, uint256)
func (_CertificationEvent *CertificationEventCaller) GetCertificationEvent(opts *bind.CallOpts, eventID string) (string, string, string, string, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _CertificationEvent.contract.Call(opts, &out, "getCertificationEvent", eventID)

	if err != nil {
		return *new(string), *new(string), *new(string), *new(string), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)
	out4 := *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	out5 := *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	out6 := *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, out4, out5, out6, err

}

// GetCertificationEvent is a free data retrieval call binding the contract method 0x3e970613.
//
// Solidity: function getCertificationEvent(string eventID) view returns(string, string, string, string, uint256, uint256, uint256)
func (_CertificationEvent *CertificationEventSession) GetCertificationEvent(eventID string) (string, string, string, string, *big.Int, *big.Int, *big.Int, error) {
	return _CertificationEvent.Contract.GetCertificationEvent(&_CertificationEvent.CallOpts, eventID)
}

// GetCertificationEvent is a free data retrieval call binding the contract method 0x3e970613.
//
// Solidity: function getCertificationEvent(string eventID) view returns(string, string, string, string, uint256, uint256, uint256)
func (_CertificationEvent *CertificationEventCallerSession) GetCertificationEvent(eventID string) (string, string, string, string, *big.Int, *big.Int, *big.Int, error) {
	return _CertificationEvent.Contract.GetCertificationEvent(&_CertificationEvent.CallOpts, eventID)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_CertificationEvent *CertificationEventTransactor) StoreCertificationEvent(opts *bind.TransactOpts, eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _CertificationEvent.contract.Transact(opts, "storeCertificationEvent", eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_CertificationEvent *CertificationEventSession) StoreCertificationEvent(eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _CertificationEvent.Contract.StoreCertificationEvent(&_CertificationEvent.TransactOpts, eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// StoreCertificationEvent is a paid mutator transaction binding the contract method 0x046c0d9f.
//
// Solidity: function storeCertificationEvent(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate) returns()
func (_CertificationEvent *CertificationEventTransactorSession) StoreCertificationEvent(eventID string, entityType string, entityID string, certificationCID string, issuedDate *big.Int, expiryDate *big.Int) (*types.Transaction, error) {
	return _CertificationEvent.Contract.StoreCertificationEvent(&_CertificationEvent.TransactOpts, eventID, entityType, entityID, certificationCID, issuedDate, expiryDate)
}

// CertificationEventCertificationEventStoredIterator is returned from FilterCertificationEventStored and is used to iterate over the raw logs and unpacked data for CertificationEventStored events raised by the CertificationEvent contract.
type CertificationEventCertificationEventStoredIterator struct {
	Event *CertificationEventCertificationEventStored // Event containing the contract specifics and raw log

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
func (it *CertificationEventCertificationEventStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CertificationEventCertificationEventStored)
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
		it.Event = new(CertificationEventCertificationEventStored)
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
func (it *CertificationEventCertificationEventStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CertificationEventCertificationEventStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CertificationEventCertificationEventStored represents a CertificationEventStored event raised by the CertificationEvent contract.
type CertificationEventCertificationEventStored struct {
	EventID          string
	EntityType       string
	EntityID         string
	CertificationCID string
	IssuedDate       *big.Int
	ExpiryDate       *big.Int
	CreatedOn        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCertificationEventStored is a free log retrieval operation binding the contract event 0x101ae11cd6dd5de9b8f881e67ee4261a7ab8cec71c2147ad4e94dd21571fef7e.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventFilterer) FilterCertificationEventStored(opts *bind.FilterOpts) (*CertificationEventCertificationEventStoredIterator, error) {

	logs, sub, err := _CertificationEvent.contract.FilterLogs(opts, "CertificationEventStored")
	if err != nil {
		return nil, err
	}
	return &CertificationEventCertificationEventStoredIterator{contract: _CertificationEvent.contract, event: "CertificationEventStored", logs: logs, sub: sub}, nil
}

// WatchCertificationEventStored is a free log subscription operation binding the contract event 0x101ae11cd6dd5de9b8f881e67ee4261a7ab8cec71c2147ad4e94dd21571fef7e.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventFilterer) WatchCertificationEventStored(opts *bind.WatchOpts, sink chan<- *CertificationEventCertificationEventStored) (event.Subscription, error) {

	logs, sub, err := _CertificationEvent.contract.WatchLogs(opts, "CertificationEventStored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CertificationEventCertificationEventStored)
				if err := _CertificationEvent.contract.UnpackLog(event, "CertificationEventStored", log); err != nil {
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

// ParseCertificationEventStored is a log parse operation binding the contract event 0x101ae11cd6dd5de9b8f881e67ee4261a7ab8cec71c2147ad4e94dd21571fef7e.
//
// Solidity: event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn)
func (_CertificationEvent *CertificationEventFilterer) ParseCertificationEventStored(log types.Log) (*CertificationEventCertificationEventStored, error) {
	event := new(CertificationEventCertificationEventStored)
	if err := _CertificationEvent.contract.UnpackLog(event, "CertificationEventStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
