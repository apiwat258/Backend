// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package userregistry

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

// UserregistryMetaData contains all meta data concerning the Userregistry contract.
var UserregistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumUserRegistry.UserRole\",\"name\":\"role\",\"type\":\"uint8\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"enumUserRegistry.UserRole\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"isRegistered\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"enumUserRegistry.UserRole\",\"name\":\"role\",\"type\":\"uint8\"}],\"name\":\"registerUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"}],\"name\":\"getUserRole\",\"outputs\":[{\"internalType\":\"enumUserRegistry.UserRole\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
}

// UserregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use UserregistryMetaData.ABI instead.
var UserregistryABI = UserregistryMetaData.ABI

// Userregistry is an auto generated Go binding around an Ethereum contract.
type Userregistry struct {
	UserregistryCaller     // Read-only binding to the contract
	UserregistryTransactor // Write-only binding to the contract
	UserregistryFilterer   // Log filterer for contract events
}

// UserregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type UserregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UserregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UserregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UserregistrySession struct {
	Contract     *Userregistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UserregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UserregistryCallerSession struct {
	Contract *UserregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// UserregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UserregistryTransactorSession struct {
	Contract     *UserregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// UserregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type UserregistryRaw struct {
	Contract *Userregistry // Generic contract binding to access the raw methods on
}

// UserregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UserregistryCallerRaw struct {
	Contract *UserregistryCaller // Generic read-only contract binding to access the raw methods on
}

// UserregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UserregistryTransactorRaw struct {
	Contract *UserregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUserregistry creates a new instance of Userregistry, bound to a specific deployed contract.
func NewUserregistry(address common.Address, backend bind.ContractBackend) (*Userregistry, error) {
	contract, err := bindUserregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Userregistry{UserregistryCaller: UserregistryCaller{contract: contract}, UserregistryTransactor: UserregistryTransactor{contract: contract}, UserregistryFilterer: UserregistryFilterer{contract: contract}}, nil
}

// NewUserregistryCaller creates a new read-only instance of Userregistry, bound to a specific deployed contract.
func NewUserregistryCaller(address common.Address, caller bind.ContractCaller) (*UserregistryCaller, error) {
	contract, err := bindUserregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UserregistryCaller{contract: contract}, nil
}

// NewUserregistryTransactor creates a new write-only instance of Userregistry, bound to a specific deployed contract.
func NewUserregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*UserregistryTransactor, error) {
	contract, err := bindUserregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UserregistryTransactor{contract: contract}, nil
}

// NewUserregistryFilterer creates a new log filterer instance of Userregistry, bound to a specific deployed contract.
func NewUserregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*UserregistryFilterer, error) {
	contract, err := bindUserregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UserregistryFilterer{contract: contract}, nil
}

// bindUserregistry binds a generic wrapper to an already deployed contract.
func bindUserregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UserregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Userregistry *UserregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Userregistry.Contract.UserregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Userregistry *UserregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Userregistry.Contract.UserregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Userregistry *UserregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Userregistry.Contract.UserregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Userregistry *UserregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Userregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Userregistry *UserregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Userregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Userregistry *UserregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Userregistry.Contract.contract.Transact(opts, method, params...)
}

// GetUserRole is a free data retrieval call binding the contract method 0x27820851.
//
// Solidity: function getUserRole(address wallet) view returns(uint8)
func (_Userregistry *UserregistryCaller) GetUserRole(opts *bind.CallOpts, wallet common.Address) (uint8, error) {
	var out []interface{}
	err := _Userregistry.contract.Call(opts, &out, "getUserRole", wallet)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetUserRole is a free data retrieval call binding the contract method 0x27820851.
//
// Solidity: function getUserRole(address wallet) view returns(uint8)
func (_Userregistry *UserregistrySession) GetUserRole(wallet common.Address) (uint8, error) {
	return _Userregistry.Contract.GetUserRole(&_Userregistry.CallOpts, wallet)
}

// GetUserRole is a free data retrieval call binding the contract method 0x27820851.
//
// Solidity: function getUserRole(address wallet) view returns(uint8)
func (_Userregistry *UserregistryCallerSession) GetUserRole(wallet common.Address) (uint8, error) {
	return _Userregistry.Contract.GetUserRole(&_Userregistry.CallOpts, wallet)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(address wallet) view returns(bool)
func (_Userregistry *UserregistryCaller) IsUserRegistered(opts *bind.CallOpts, wallet common.Address) (bool, error) {
	var out []interface{}
	err := _Userregistry.contract.Call(opts, &out, "isUserRegistered", wallet)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(address wallet) view returns(bool)
func (_Userregistry *UserregistrySession) IsUserRegistered(wallet common.Address) (bool, error) {
	return _Userregistry.Contract.IsUserRegistered(&_Userregistry.CallOpts, wallet)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0x163f7522.
//
// Solidity: function isUserRegistered(address wallet) view returns(bool)
func (_Userregistry *UserregistryCallerSession) IsUserRegistered(wallet common.Address) (bool, error) {
	return _Userregistry.Contract.IsUserRegistered(&_Userregistry.CallOpts, wallet)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(address wallet, uint8 role, bool isRegistered)
func (_Userregistry *UserregistryCaller) Users(opts *bind.CallOpts, arg0 common.Address) (struct {
	Wallet       common.Address
	Role         uint8
	IsRegistered bool
}, error) {
	var out []interface{}
	err := _Userregistry.contract.Call(opts, &out, "users", arg0)

	outstruct := new(struct {
		Wallet       common.Address
		Role         uint8
		IsRegistered bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Wallet = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Role = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.IsRegistered = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(address wallet, uint8 role, bool isRegistered)
func (_Userregistry *UserregistrySession) Users(arg0 common.Address) (struct {
	Wallet       common.Address
	Role         uint8
	IsRegistered bool
}, error) {
	return _Userregistry.Contract.Users(&_Userregistry.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(address wallet, uint8 role, bool isRegistered)
func (_Userregistry *UserregistryCallerSession) Users(arg0 common.Address) (struct {
	Wallet       common.Address
	Role         uint8
	IsRegistered bool
}, error) {
	return _Userregistry.Contract.Users(&_Userregistry.CallOpts, arg0)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x02b8c8cd.
//
// Solidity: function registerUser(uint8 role) returns()
func (_Userregistry *UserregistryTransactor) RegisterUser(opts *bind.TransactOpts, role uint8) (*types.Transaction, error) {
	return _Userregistry.contract.Transact(opts, "registerUser", role)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x02b8c8cd.
//
// Solidity: function registerUser(uint8 role) returns()
func (_Userregistry *UserregistrySession) RegisterUser(role uint8) (*types.Transaction, error) {
	return _Userregistry.Contract.RegisterUser(&_Userregistry.TransactOpts, role)
}

// RegisterUser is a paid mutator transaction binding the contract method 0x02b8c8cd.
//
// Solidity: function registerUser(uint8 role) returns()
func (_Userregistry *UserregistryTransactorSession) RegisterUser(role uint8) (*types.Transaction, error) {
	return _Userregistry.Contract.RegisterUser(&_Userregistry.TransactOpts, role)
}

// UserregistryUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the Userregistry contract.
type UserregistryUserRegisteredIterator struct {
	Event *UserregistryUserRegistered // Event containing the contract specifics and raw log

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
func (it *UserregistryUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UserregistryUserRegistered)
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
		it.Event = new(UserregistryUserRegistered)
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
func (it *UserregistryUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UserregistryUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UserregistryUserRegistered represents a UserRegistered event raised by the Userregistry contract.
type UserregistryUserRegistered struct {
	Wallet common.Address
	Role   uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x8e18574f8a7333e6a87112742a55853a643e2f0550f4e1a3a2c3f1bb2994441d.
//
// Solidity: event UserRegistered(address indexed wallet, uint8 role)
func (_Userregistry *UserregistryFilterer) FilterUserRegistered(opts *bind.FilterOpts, wallet []common.Address) (*UserregistryUserRegisteredIterator, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Userregistry.contract.FilterLogs(opts, "UserRegistered", walletRule)
	if err != nil {
		return nil, err
	}
	return &UserregistryUserRegisteredIterator{contract: _Userregistry.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x8e18574f8a7333e6a87112742a55853a643e2f0550f4e1a3a2c3f1bb2994441d.
//
// Solidity: event UserRegistered(address indexed wallet, uint8 role)
func (_Userregistry *UserregistryFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *UserregistryUserRegistered, wallet []common.Address) (event.Subscription, error) {

	var walletRule []interface{}
	for _, walletItem := range wallet {
		walletRule = append(walletRule, walletItem)
	}

	logs, sub, err := _Userregistry.contract.WatchLogs(opts, "UserRegistered", walletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UserregistryUserRegistered)
				if err := _Userregistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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

// ParseUserRegistered is a log parse operation binding the contract event 0x8e18574f8a7333e6a87112742a55853a643e2f0550f4e1a3a2c3f1bb2994441d.
//
// Solidity: event UserRegistered(address indexed wallet, uint8 role)
func (_Userregistry *UserregistryFilterer) ParseUserRegistered(log types.Log) (*UserregistryUserRegistered, error) {
	event := new(UserregistryUserRegistered)
	if err := _Userregistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
