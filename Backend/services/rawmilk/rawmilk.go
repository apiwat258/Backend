// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rawmilk

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

// RawmilkMetaData contains all meta data concerning the Rawmilk contract.
var RawmilkMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"farmWallet\",\"type\":\"address\"}],\"name\":\"FarmRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rawMilkID\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"farmWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumRawMilkSupplyChain.MilkStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"}],\"name\":\"RawMilkAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rawMilkID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumRawMilkSupplyChain.MilkStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"RawMilkStatusUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"farms\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"rawMilkRecords\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"farmWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"temperature\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protein\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"},{\"internalType\":\"enumRawMilkSupplyChain.MilkStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_farmWallet\",\"type\":\"address\"}],\"name\":\"registerFarm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_rawMilkID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_temperature\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_pH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_protein\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_ipfsCid\",\"type\":\"string\"}],\"name\":\"addRawMilk\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_rawMilkID\",\"type\":\"bytes32\"}],\"name\":\"getRawMilk\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"farmWallet\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"temperature\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protein\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"},{\"internalType\":\"enumRawMilkSupplyChain.MilkStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_rawMilkID\",\"type\":\"bytes32\"},{\"internalType\":\"enumRawMilkSupplyChain.MilkStatus\",\"name\":\"_newStatus\",\"type\":\"uint8\"}],\"name\":\"updateRawMilkStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RawmilkABI is the input ABI used to generate the binding from.
// Deprecated: Use RawmilkMetaData.ABI instead.
var RawmilkABI = RawmilkMetaData.ABI

// Rawmilk is an auto generated Go binding around an Ethereum contract.
type Rawmilk struct {
	RawmilkCaller     // Read-only binding to the contract
	RawmilkTransactor // Write-only binding to the contract
	RawmilkFilterer   // Log filterer for contract events
}

// RawmilkCaller is an auto generated read-only Go binding around an Ethereum contract.
type RawmilkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RawmilkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RawmilkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RawmilkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RawmilkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RawmilkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RawmilkSession struct {
	Contract     *Rawmilk          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RawmilkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RawmilkCallerSession struct {
	Contract *RawmilkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// RawmilkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RawmilkTransactorSession struct {
	Contract     *RawmilkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RawmilkRaw is an auto generated low-level Go binding around an Ethereum contract.
type RawmilkRaw struct {
	Contract *Rawmilk // Generic contract binding to access the raw methods on
}

// RawmilkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RawmilkCallerRaw struct {
	Contract *RawmilkCaller // Generic read-only contract binding to access the raw methods on
}

// RawmilkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RawmilkTransactorRaw struct {
	Contract *RawmilkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRawmilk creates a new instance of Rawmilk, bound to a specific deployed contract.
func NewRawmilk(address common.Address, backend bind.ContractBackend) (*Rawmilk, error) {
	contract, err := bindRawmilk(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rawmilk{RawmilkCaller: RawmilkCaller{contract: contract}, RawmilkTransactor: RawmilkTransactor{contract: contract}, RawmilkFilterer: RawmilkFilterer{contract: contract}}, nil
}

// NewRawmilkCaller creates a new read-only instance of Rawmilk, bound to a specific deployed contract.
func NewRawmilkCaller(address common.Address, caller bind.ContractCaller) (*RawmilkCaller, error) {
	contract, err := bindRawmilk(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RawmilkCaller{contract: contract}, nil
}

// NewRawmilkTransactor creates a new write-only instance of Rawmilk, bound to a specific deployed contract.
func NewRawmilkTransactor(address common.Address, transactor bind.ContractTransactor) (*RawmilkTransactor, error) {
	contract, err := bindRawmilk(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RawmilkTransactor{contract: contract}, nil
}

// NewRawmilkFilterer creates a new log filterer instance of Rawmilk, bound to a specific deployed contract.
func NewRawmilkFilterer(address common.Address, filterer bind.ContractFilterer) (*RawmilkFilterer, error) {
	contract, err := bindRawmilk(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RawmilkFilterer{contract: contract}, nil
}

// bindRawmilk binds a generic wrapper to an already deployed contract.
func bindRawmilk(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RawmilkMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rawmilk *RawmilkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rawmilk.Contract.RawmilkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rawmilk *RawmilkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rawmilk.Contract.RawmilkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rawmilk *RawmilkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rawmilk.Contract.RawmilkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rawmilk *RawmilkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rawmilk.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rawmilk *RawmilkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rawmilk.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rawmilk *RawmilkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rawmilk.Contract.contract.Transact(opts, method, params...)
}

// Farms is a free data retrieval call binding the contract method 0x421adfa0.
//
// Solidity: function farms(address ) view returns(bool)
func (_Rawmilk *RawmilkCaller) Farms(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "farms", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Farms is a free data retrieval call binding the contract method 0x421adfa0.
//
// Solidity: function farms(address ) view returns(bool)
func (_Rawmilk *RawmilkSession) Farms(arg0 common.Address) (bool, error) {
	return _Rawmilk.Contract.Farms(&_Rawmilk.CallOpts, arg0)
}

// Farms is a free data retrieval call binding the contract method 0x421adfa0.
//
// Solidity: function farms(address ) view returns(bool)
func (_Rawmilk *RawmilkCallerSession) Farms(arg0 common.Address) (bool, error) {
	return _Rawmilk.Contract.Farms(&_Rawmilk.CallOpts, arg0)
}

// GetRawMilk is a free data retrieval call binding the contract method 0xb0ae6828.
//
// Solidity: function getRawMilk(bytes32 _rawMilkID) view returns(address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkCaller) GetRawMilk(opts *bind.CallOpts, _rawMilkID [32]byte) (struct {
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "getRawMilk", _rawMilkID)

	outstruct := new(struct {
		FarmWallet  common.Address
		Temperature *big.Int
		PH          *big.Int
		Fat         *big.Int
		Protein     *big.Int
		IpfsCid     string
		Status      uint8
		Timestamp   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FarmWallet = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Temperature = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PH = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Fat = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Protein = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IpfsCid = *abi.ConvertType(out[5], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.Timestamp = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRawMilk is a free data retrieval call binding the contract method 0xb0ae6828.
//
// Solidity: function getRawMilk(bytes32 _rawMilkID) view returns(address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkSession) GetRawMilk(_rawMilkID [32]byte) (struct {
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	return _Rawmilk.Contract.GetRawMilk(&_Rawmilk.CallOpts, _rawMilkID)
}

// GetRawMilk is a free data retrieval call binding the contract method 0xb0ae6828.
//
// Solidity: function getRawMilk(bytes32 _rawMilkID) view returns(address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkCallerSession) GetRawMilk(_rawMilkID [32]byte) (struct {
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	return _Rawmilk.Contract.GetRawMilk(&_Rawmilk.CallOpts, _rawMilkID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rawmilk *RawmilkCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rawmilk *RawmilkSession) Owner() (common.Address, error) {
	return _Rawmilk.Contract.Owner(&_Rawmilk.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rawmilk *RawmilkCallerSession) Owner() (common.Address, error) {
	return _Rawmilk.Contract.Owner(&_Rawmilk.CallOpts)
}

// RawMilkRecords is a free data retrieval call binding the contract method 0xdbaee037.
//
// Solidity: function rawMilkRecords(bytes32 ) view returns(bytes32 id, address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkCaller) RawMilkRecords(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Id          [32]byte
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "rawMilkRecords", arg0)

	outstruct := new(struct {
		Id          [32]byte
		FarmWallet  common.Address
		Temperature *big.Int
		PH          *big.Int
		Fat         *big.Int
		Protein     *big.Int
		IpfsCid     string
		Status      uint8
		Timestamp   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FarmWallet = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Temperature = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PH = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Fat = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Protein = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.IpfsCid = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.Timestamp = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RawMilkRecords is a free data retrieval call binding the contract method 0xdbaee037.
//
// Solidity: function rawMilkRecords(bytes32 ) view returns(bytes32 id, address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkSession) RawMilkRecords(arg0 [32]byte) (struct {
	Id          [32]byte
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	return _Rawmilk.Contract.RawMilkRecords(&_Rawmilk.CallOpts, arg0)
}

// RawMilkRecords is a free data retrieval call binding the contract method 0xdbaee037.
//
// Solidity: function rawMilkRecords(bytes32 ) view returns(bytes32 id, address farmWallet, uint256 temperature, uint256 pH, uint256 fat, uint256 protein, string ipfsCid, uint8 status, uint256 timestamp)
func (_Rawmilk *RawmilkCallerSession) RawMilkRecords(arg0 [32]byte) (struct {
	Id          [32]byte
	FarmWallet  common.Address
	Temperature *big.Int
	PH          *big.Int
	Fat         *big.Int
	Protein     *big.Int
	IpfsCid     string
	Status      uint8
	Timestamp   *big.Int
}, error) {
	return _Rawmilk.Contract.RawMilkRecords(&_Rawmilk.CallOpts, arg0)
}

// AddRawMilk is a paid mutator transaction binding the contract method 0x0fd802c8.
//
// Solidity: function addRawMilk(bytes32 _rawMilkID, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, string _ipfsCid) returns()
func (_Rawmilk *RawmilkTransactor) AddRawMilk(opts *bind.TransactOpts, _rawMilkID [32]byte, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _ipfsCid string) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "addRawMilk", _rawMilkID, _temperature, _pH, _fat, _protein, _ipfsCid)
}

// AddRawMilk is a paid mutator transaction binding the contract method 0x0fd802c8.
//
// Solidity: function addRawMilk(bytes32 _rawMilkID, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, string _ipfsCid) returns()
func (_Rawmilk *RawmilkSession) AddRawMilk(_rawMilkID [32]byte, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _ipfsCid string) (*types.Transaction, error) {
	return _Rawmilk.Contract.AddRawMilk(&_Rawmilk.TransactOpts, _rawMilkID, _temperature, _pH, _fat, _protein, _ipfsCid)
}

// AddRawMilk is a paid mutator transaction binding the contract method 0x0fd802c8.
//
// Solidity: function addRawMilk(bytes32 _rawMilkID, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, string _ipfsCid) returns()
func (_Rawmilk *RawmilkTransactorSession) AddRawMilk(_rawMilkID [32]byte, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _ipfsCid string) (*types.Transaction, error) {
	return _Rawmilk.Contract.AddRawMilk(&_Rawmilk.TransactOpts, _rawMilkID, _temperature, _pH, _fat, _protein, _ipfsCid)
}

// RegisterFarm is a paid mutator transaction binding the contract method 0xedcd5d17.
//
// Solidity: function registerFarm(address _farmWallet) returns()
func (_Rawmilk *RawmilkTransactor) RegisterFarm(opts *bind.TransactOpts, _farmWallet common.Address) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "registerFarm", _farmWallet)
}

// RegisterFarm is a paid mutator transaction binding the contract method 0xedcd5d17.
//
// Solidity: function registerFarm(address _farmWallet) returns()
func (_Rawmilk *RawmilkSession) RegisterFarm(_farmWallet common.Address) (*types.Transaction, error) {
	return _Rawmilk.Contract.RegisterFarm(&_Rawmilk.TransactOpts, _farmWallet)
}

// RegisterFarm is a paid mutator transaction binding the contract method 0xedcd5d17.
//
// Solidity: function registerFarm(address _farmWallet) returns()
func (_Rawmilk *RawmilkTransactorSession) RegisterFarm(_farmWallet common.Address) (*types.Transaction, error) {
	return _Rawmilk.Contract.RegisterFarm(&_Rawmilk.TransactOpts, _farmWallet)
}

// UpdateRawMilkStatus is a paid mutator transaction binding the contract method 0xc90b88a7.
//
// Solidity: function updateRawMilkStatus(bytes32 _rawMilkID, uint8 _newStatus) returns()
func (_Rawmilk *RawmilkTransactor) UpdateRawMilkStatus(opts *bind.TransactOpts, _rawMilkID [32]byte, _newStatus uint8) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "updateRawMilkStatus", _rawMilkID, _newStatus)
}

// UpdateRawMilkStatus is a paid mutator transaction binding the contract method 0xc90b88a7.
//
// Solidity: function updateRawMilkStatus(bytes32 _rawMilkID, uint8 _newStatus) returns()
func (_Rawmilk *RawmilkSession) UpdateRawMilkStatus(_rawMilkID [32]byte, _newStatus uint8) (*types.Transaction, error) {
	return _Rawmilk.Contract.UpdateRawMilkStatus(&_Rawmilk.TransactOpts, _rawMilkID, _newStatus)
}

// UpdateRawMilkStatus is a paid mutator transaction binding the contract method 0xc90b88a7.
//
// Solidity: function updateRawMilkStatus(bytes32 _rawMilkID, uint8 _newStatus) returns()
func (_Rawmilk *RawmilkTransactorSession) UpdateRawMilkStatus(_rawMilkID [32]byte, _newStatus uint8) (*types.Transaction, error) {
	return _Rawmilk.Contract.UpdateRawMilkStatus(&_Rawmilk.TransactOpts, _rawMilkID, _newStatus)
}

// RawmilkFarmRegisteredIterator is returned from FilterFarmRegistered and is used to iterate over the raw logs and unpacked data for FarmRegistered events raised by the Rawmilk contract.
type RawmilkFarmRegisteredIterator struct {
	Event *RawmilkFarmRegistered // Event containing the contract specifics and raw log

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
func (it *RawmilkFarmRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkFarmRegistered)
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
		it.Event = new(RawmilkFarmRegistered)
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
func (it *RawmilkFarmRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkFarmRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkFarmRegistered represents a FarmRegistered event raised by the Rawmilk contract.
type RawmilkFarmRegistered struct {
	FarmWallet common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFarmRegistered is a free log retrieval operation binding the contract event 0x34deede75423e085ceb61f47c1ddb4466e42dcf614d55b7dbb1df16e64e9d5d4.
//
// Solidity: event FarmRegistered(address indexed farmWallet)
func (_Rawmilk *RawmilkFilterer) FilterFarmRegistered(opts *bind.FilterOpts, farmWallet []common.Address) (*RawmilkFarmRegisteredIterator, error) {

	var farmWalletRule []interface{}
	for _, farmWalletItem := range farmWallet {
		farmWalletRule = append(farmWalletRule, farmWalletItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "FarmRegistered", farmWalletRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkFarmRegisteredIterator{contract: _Rawmilk.contract, event: "FarmRegistered", logs: logs, sub: sub}, nil
}

// WatchFarmRegistered is a free log subscription operation binding the contract event 0x34deede75423e085ceb61f47c1ddb4466e42dcf614d55b7dbb1df16e64e9d5d4.
//
// Solidity: event FarmRegistered(address indexed farmWallet)
func (_Rawmilk *RawmilkFilterer) WatchFarmRegistered(opts *bind.WatchOpts, sink chan<- *RawmilkFarmRegistered, farmWallet []common.Address) (event.Subscription, error) {

	var farmWalletRule []interface{}
	for _, farmWalletItem := range farmWallet {
		farmWalletRule = append(farmWalletRule, farmWalletItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "FarmRegistered", farmWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkFarmRegistered)
				if err := _Rawmilk.contract.UnpackLog(event, "FarmRegistered", log); err != nil {
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

// ParseFarmRegistered is a log parse operation binding the contract event 0x34deede75423e085ceb61f47c1ddb4466e42dcf614d55b7dbb1df16e64e9d5d4.
//
// Solidity: event FarmRegistered(address indexed farmWallet)
func (_Rawmilk *RawmilkFilterer) ParseFarmRegistered(log types.Log) (*RawmilkFarmRegistered, error) {
	event := new(RawmilkFarmRegistered)
	if err := _Rawmilk.contract.UnpackLog(event, "FarmRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RawmilkRawMilkAddedIterator is returned from FilterRawMilkAdded and is used to iterate over the raw logs and unpacked data for RawMilkAdded events raised by the Rawmilk contract.
type RawmilkRawMilkAddedIterator struct {
	Event *RawmilkRawMilkAdded // Event containing the contract specifics and raw log

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
func (it *RawmilkRawMilkAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkRawMilkAdded)
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
		it.Event = new(RawmilkRawMilkAdded)
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
func (it *RawmilkRawMilkAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkRawMilkAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkRawMilkAdded represents a RawMilkAdded event raised by the Rawmilk contract.
type RawmilkRawMilkAdded struct {
	RawMilkID  [32]byte
	FarmWallet common.Address
	Status     uint8
	IpfsCid    string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRawMilkAdded is a free log retrieval operation binding the contract event 0x44bf3f1c16b69967a6b0ba8d8030943e8d61dbf76ac4e0d5970c0172e2c99af2.
//
// Solidity: event RawMilkAdded(bytes32 indexed rawMilkID, address indexed farmWallet, uint8 status, string ipfsCid)
func (_Rawmilk *RawmilkFilterer) FilterRawMilkAdded(opts *bind.FilterOpts, rawMilkID [][32]byte, farmWallet []common.Address) (*RawmilkRawMilkAddedIterator, error) {

	var rawMilkIDRule []interface{}
	for _, rawMilkIDItem := range rawMilkID {
		rawMilkIDRule = append(rawMilkIDRule, rawMilkIDItem)
	}
	var farmWalletRule []interface{}
	for _, farmWalletItem := range farmWallet {
		farmWalletRule = append(farmWalletRule, farmWalletItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "RawMilkAdded", rawMilkIDRule, farmWalletRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkRawMilkAddedIterator{contract: _Rawmilk.contract, event: "RawMilkAdded", logs: logs, sub: sub}, nil
}

// WatchRawMilkAdded is a free log subscription operation binding the contract event 0x44bf3f1c16b69967a6b0ba8d8030943e8d61dbf76ac4e0d5970c0172e2c99af2.
//
// Solidity: event RawMilkAdded(bytes32 indexed rawMilkID, address indexed farmWallet, uint8 status, string ipfsCid)
func (_Rawmilk *RawmilkFilterer) WatchRawMilkAdded(opts *bind.WatchOpts, sink chan<- *RawmilkRawMilkAdded, rawMilkID [][32]byte, farmWallet []common.Address) (event.Subscription, error) {

	var rawMilkIDRule []interface{}
	for _, rawMilkIDItem := range rawMilkID {
		rawMilkIDRule = append(rawMilkIDRule, rawMilkIDItem)
	}
	var farmWalletRule []interface{}
	for _, farmWalletItem := range farmWallet {
		farmWalletRule = append(farmWalletRule, farmWalletItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "RawMilkAdded", rawMilkIDRule, farmWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkRawMilkAdded)
				if err := _Rawmilk.contract.UnpackLog(event, "RawMilkAdded", log); err != nil {
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

// ParseRawMilkAdded is a log parse operation binding the contract event 0x44bf3f1c16b69967a6b0ba8d8030943e8d61dbf76ac4e0d5970c0172e2c99af2.
//
// Solidity: event RawMilkAdded(bytes32 indexed rawMilkID, address indexed farmWallet, uint8 status, string ipfsCid)
func (_Rawmilk *RawmilkFilterer) ParseRawMilkAdded(log types.Log) (*RawmilkRawMilkAdded, error) {
	event := new(RawmilkRawMilkAdded)
	if err := _Rawmilk.contract.UnpackLog(event, "RawMilkAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RawmilkRawMilkStatusUpdatedIterator is returned from FilterRawMilkStatusUpdated and is used to iterate over the raw logs and unpacked data for RawMilkStatusUpdated events raised by the Rawmilk contract.
type RawmilkRawMilkStatusUpdatedIterator struct {
	Event *RawmilkRawMilkStatusUpdated // Event containing the contract specifics and raw log

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
func (it *RawmilkRawMilkStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkRawMilkStatusUpdated)
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
		it.Event = new(RawmilkRawMilkStatusUpdated)
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
func (it *RawmilkRawMilkStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkRawMilkStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkRawMilkStatusUpdated represents a RawMilkStatusUpdated event raised by the Rawmilk contract.
type RawmilkRawMilkStatusUpdated struct {
	RawMilkID [32]byte
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRawMilkStatusUpdated is a free log retrieval operation binding the contract event 0xf51f46360a6d2d804a8686b98ecb0ca33aeb53a6f3ff77cb0b25091f96f804fd.
//
// Solidity: event RawMilkStatusUpdated(bytes32 indexed rawMilkID, uint8 newStatus)
func (_Rawmilk *RawmilkFilterer) FilterRawMilkStatusUpdated(opts *bind.FilterOpts, rawMilkID [][32]byte) (*RawmilkRawMilkStatusUpdatedIterator, error) {

	var rawMilkIDRule []interface{}
	for _, rawMilkIDItem := range rawMilkID {
		rawMilkIDRule = append(rawMilkIDRule, rawMilkIDItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "RawMilkStatusUpdated", rawMilkIDRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkRawMilkStatusUpdatedIterator{contract: _Rawmilk.contract, event: "RawMilkStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchRawMilkStatusUpdated is a free log subscription operation binding the contract event 0xf51f46360a6d2d804a8686b98ecb0ca33aeb53a6f3ff77cb0b25091f96f804fd.
//
// Solidity: event RawMilkStatusUpdated(bytes32 indexed rawMilkID, uint8 newStatus)
func (_Rawmilk *RawmilkFilterer) WatchRawMilkStatusUpdated(opts *bind.WatchOpts, sink chan<- *RawmilkRawMilkStatusUpdated, rawMilkID [][32]byte) (event.Subscription, error) {

	var rawMilkIDRule []interface{}
	for _, rawMilkIDItem := range rawMilkID {
		rawMilkIDRule = append(rawMilkIDRule, rawMilkIDItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "RawMilkStatusUpdated", rawMilkIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkRawMilkStatusUpdated)
				if err := _Rawmilk.contract.UnpackLog(event, "RawMilkStatusUpdated", log); err != nil {
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

// ParseRawMilkStatusUpdated is a log parse operation binding the contract event 0xf51f46360a6d2d804a8686b98ecb0ca33aeb53a6f3ff77cb0b25091f96f804fd.
//
// Solidity: event RawMilkStatusUpdated(bytes32 indexed rawMilkID, uint8 newStatus)
func (_Rawmilk *RawmilkFilterer) ParseRawMilkStatusUpdated(log types.Log) (*RawmilkRawMilkStatusUpdated, error) {
	event := new(RawmilkRawMilkStatusUpdated)
	if err := _Rawmilk.contract.UnpackLog(event, "RawMilkStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
