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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"tankId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumRawMilk.MilkStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"qualityReportCID\",\"type\":\"string\"}],\"name\":\"MilkQualityVerified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"tankId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"farmer\",\"type\":\"address\"}],\"name\":\"MilkTankCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"tankId\",\"type\":\"bytes32\"}],\"name\":\"MilkTankUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"milkTanks\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"tankId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"farmer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumRawMilk.MilkStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"qualityReportCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"qrCodeCID\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tankIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"userRegistry\",\"outputs\":[{\"internalType\":\"contractUserRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_quantity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_temperature\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_pH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_protein\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_bacteria\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_contaminants\",\"type\":\"bool\"}],\"name\":\"validateMilkData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_tankId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_personInCharge\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_qrCodeCID\",\"type\":\"string\"}],\"name\":\"createMilkTank\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_tankId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_qrCodeCID\",\"type\":\"string\"}],\"name\":\"updateMilkTank\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_tankId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"_approved\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"_qualityReportCID\",\"type\":\"string\"}],\"name\":\"verifyMilkQuality\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_tankId\",\"type\":\"bytes32\"}],\"name\":\"getMilkTank\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"enumRawMilk.MilkStatus\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"getAllMilkTanks\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_farmer\",\"type\":\"address\"}],\"name\":\"getMilkTanksByFarmer\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
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

// GetAllMilkTanks is a free data retrieval call binding the contract method 0xae2cc61e.
//
// Solidity: function getAllMilkTanks() view returns(bytes32[])
func (_Rawmilk *RawmilkCaller) GetAllMilkTanks(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "getAllMilkTanks")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetAllMilkTanks is a free data retrieval call binding the contract method 0xae2cc61e.
//
// Solidity: function getAllMilkTanks() view returns(bytes32[])
func (_Rawmilk *RawmilkSession) GetAllMilkTanks() ([][32]byte, error) {
	return _Rawmilk.Contract.GetAllMilkTanks(&_Rawmilk.CallOpts)
}

// GetAllMilkTanks is a free data retrieval call binding the contract method 0xae2cc61e.
//
// Solidity: function getAllMilkTanks() view returns(bytes32[])
func (_Rawmilk *RawmilkCallerSession) GetAllMilkTanks() ([][32]byte, error) {
	return _Rawmilk.Contract.GetAllMilkTanks(&_Rawmilk.CallOpts)
}

// GetMilkTank is a free data retrieval call binding the contract method 0x5a03541f.
//
// Solidity: function getMilkTank(bytes32 _tankId) view returns(bytes32, address, string, uint8, string, string)
func (_Rawmilk *RawmilkCaller) GetMilkTank(opts *bind.CallOpts, _tankId [32]byte) ([32]byte, common.Address, string, uint8, string, string, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "getMilkTank", _tankId)

	if err != nil {
		return *new([32]byte), *new(common.Address), *new(string), *new(uint8), *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(uint8)).(*uint8)
	out4 := *abi.ConvertType(out[4], new(string)).(*string)
	out5 := *abi.ConvertType(out[5], new(string)).(*string)

	return out0, out1, out2, out3, out4, out5, err

}

// GetMilkTank is a free data retrieval call binding the contract method 0x5a03541f.
//
// Solidity: function getMilkTank(bytes32 _tankId) view returns(bytes32, address, string, uint8, string, string)
func (_Rawmilk *RawmilkSession) GetMilkTank(_tankId [32]byte) ([32]byte, common.Address, string, uint8, string, string, error) {
	return _Rawmilk.Contract.GetMilkTank(&_Rawmilk.CallOpts, _tankId)
}

// GetMilkTank is a free data retrieval call binding the contract method 0x5a03541f.
//
// Solidity: function getMilkTank(bytes32 _tankId) view returns(bytes32, address, string, uint8, string, string)
func (_Rawmilk *RawmilkCallerSession) GetMilkTank(_tankId [32]byte) ([32]byte, common.Address, string, uint8, string, string, error) {
	return _Rawmilk.Contract.GetMilkTank(&_Rawmilk.CallOpts, _tankId)
}

// GetMilkTanksByFarmer is a free data retrieval call binding the contract method 0x70d61bba.
//
// Solidity: function getMilkTanksByFarmer(address _farmer) view returns(bytes32[])
func (_Rawmilk *RawmilkCaller) GetMilkTanksByFarmer(opts *bind.CallOpts, _farmer common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "getMilkTanksByFarmer", _farmer)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetMilkTanksByFarmer is a free data retrieval call binding the contract method 0x70d61bba.
//
// Solidity: function getMilkTanksByFarmer(address _farmer) view returns(bytes32[])
func (_Rawmilk *RawmilkSession) GetMilkTanksByFarmer(_farmer common.Address) ([][32]byte, error) {
	return _Rawmilk.Contract.GetMilkTanksByFarmer(&_Rawmilk.CallOpts, _farmer)
}

// GetMilkTanksByFarmer is a free data retrieval call binding the contract method 0x70d61bba.
//
// Solidity: function getMilkTanksByFarmer(address _farmer) view returns(bytes32[])
func (_Rawmilk *RawmilkCallerSession) GetMilkTanksByFarmer(_farmer common.Address) ([][32]byte, error) {
	return _Rawmilk.Contract.GetMilkTanksByFarmer(&_Rawmilk.CallOpts, _farmer)
}

// MilkTanks is a free data retrieval call binding the contract method 0xa34a7fb9.
//
// Solidity: function milkTanks(bytes32 ) view returns(bytes32 tankId, address farmer, string personInCharge, uint8 status, string qualityReportCID, string qrCodeCID)
func (_Rawmilk *RawmilkCaller) MilkTanks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TankId           [32]byte
	Farmer           common.Address
	PersonInCharge   string
	Status           uint8
	QualityReportCID string
	QrCodeCID        string
}, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "milkTanks", arg0)

	outstruct := new(struct {
		TankId           [32]byte
		Farmer           common.Address
		PersonInCharge   string
		Status           uint8
		QualityReportCID string
		QrCodeCID        string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TankId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Farmer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.PersonInCharge = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.QualityReportCID = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.QrCodeCID = *abi.ConvertType(out[5], new(string)).(*string)

	return *outstruct, err

}

// MilkTanks is a free data retrieval call binding the contract method 0xa34a7fb9.
//
// Solidity: function milkTanks(bytes32 ) view returns(bytes32 tankId, address farmer, string personInCharge, uint8 status, string qualityReportCID, string qrCodeCID)
func (_Rawmilk *RawmilkSession) MilkTanks(arg0 [32]byte) (struct {
	TankId           [32]byte
	Farmer           common.Address
	PersonInCharge   string
	Status           uint8
	QualityReportCID string
	QrCodeCID        string
}, error) {
	return _Rawmilk.Contract.MilkTanks(&_Rawmilk.CallOpts, arg0)
}

// MilkTanks is a free data retrieval call binding the contract method 0xa34a7fb9.
//
// Solidity: function milkTanks(bytes32 ) view returns(bytes32 tankId, address farmer, string personInCharge, uint8 status, string qualityReportCID, string qrCodeCID)
func (_Rawmilk *RawmilkCallerSession) MilkTanks(arg0 [32]byte) (struct {
	TankId           [32]byte
	Farmer           common.Address
	PersonInCharge   string
	Status           uint8
	QualityReportCID string
	QrCodeCID        string
}, error) {
	return _Rawmilk.Contract.MilkTanks(&_Rawmilk.CallOpts, arg0)
}

// TankIds is a free data retrieval call binding the contract method 0x4b89e5f6.
//
// Solidity: function tankIds(uint256 ) view returns(bytes32)
func (_Rawmilk *RawmilkCaller) TankIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "tankIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TankIds is a free data retrieval call binding the contract method 0x4b89e5f6.
//
// Solidity: function tankIds(uint256 ) view returns(bytes32)
func (_Rawmilk *RawmilkSession) TankIds(arg0 *big.Int) ([32]byte, error) {
	return _Rawmilk.Contract.TankIds(&_Rawmilk.CallOpts, arg0)
}

// TankIds is a free data retrieval call binding the contract method 0x4b89e5f6.
//
// Solidity: function tankIds(uint256 ) view returns(bytes32)
func (_Rawmilk *RawmilkCallerSession) TankIds(arg0 *big.Int) ([32]byte, error) {
	return _Rawmilk.Contract.TankIds(&_Rawmilk.CallOpts, arg0)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Rawmilk *RawmilkCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Rawmilk *RawmilkSession) UserRegistry() (common.Address, error) {
	return _Rawmilk.Contract.UserRegistry(&_Rawmilk.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Rawmilk *RawmilkCallerSession) UserRegistry() (common.Address, error) {
	return _Rawmilk.Contract.UserRegistry(&_Rawmilk.CallOpts)
}

// ValidateMilkData is a free data retrieval call binding the contract method 0x356f30c6.
//
// Solidity: function validateMilkData(uint256 _quantity, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, bool _bacteria, bool _contaminants) pure returns(bool valid, string message)
func (_Rawmilk *RawmilkCaller) ValidateMilkData(opts *bind.CallOpts, _quantity *big.Int, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _bacteria bool, _contaminants bool) (struct {
	Valid   bool
	Message string
}, error) {
	var out []interface{}
	err := _Rawmilk.contract.Call(opts, &out, "validateMilkData", _quantity, _temperature, _pH, _fat, _protein, _bacteria, _contaminants)

	outstruct := new(struct {
		Valid   bool
		Message string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Valid = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Message = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// ValidateMilkData is a free data retrieval call binding the contract method 0x356f30c6.
//
// Solidity: function validateMilkData(uint256 _quantity, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, bool _bacteria, bool _contaminants) pure returns(bool valid, string message)
func (_Rawmilk *RawmilkSession) ValidateMilkData(_quantity *big.Int, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _bacteria bool, _contaminants bool) (struct {
	Valid   bool
	Message string
}, error) {
	return _Rawmilk.Contract.ValidateMilkData(&_Rawmilk.CallOpts, _quantity, _temperature, _pH, _fat, _protein, _bacteria, _contaminants)
}

// ValidateMilkData is a free data retrieval call binding the contract method 0x356f30c6.
//
// Solidity: function validateMilkData(uint256 _quantity, uint256 _temperature, uint256 _pH, uint256 _fat, uint256 _protein, bool _bacteria, bool _contaminants) pure returns(bool valid, string message)
func (_Rawmilk *RawmilkCallerSession) ValidateMilkData(_quantity *big.Int, _temperature *big.Int, _pH *big.Int, _fat *big.Int, _protein *big.Int, _bacteria bool, _contaminants bool) (struct {
	Valid   bool
	Message string
}, error) {
	return _Rawmilk.Contract.ValidateMilkData(&_Rawmilk.CallOpts, _quantity, _temperature, _pH, _fat, _protein, _bacteria, _contaminants)
}

// CreateMilkTank is a paid mutator transaction binding the contract method 0x916b4775.
//
// Solidity: function createMilkTank(bytes32 _tankId, string _personInCharge, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkTransactor) CreateMilkTank(opts *bind.TransactOpts, _tankId [32]byte, _personInCharge string, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "createMilkTank", _tankId, _personInCharge, _qrCodeCID)
}

// CreateMilkTank is a paid mutator transaction binding the contract method 0x916b4775.
//
// Solidity: function createMilkTank(bytes32 _tankId, string _personInCharge, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkSession) CreateMilkTank(_tankId [32]byte, _personInCharge string, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.CreateMilkTank(&_Rawmilk.TransactOpts, _tankId, _personInCharge, _qrCodeCID)
}

// CreateMilkTank is a paid mutator transaction binding the contract method 0x916b4775.
//
// Solidity: function createMilkTank(bytes32 _tankId, string _personInCharge, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkTransactorSession) CreateMilkTank(_tankId [32]byte, _personInCharge string, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.CreateMilkTank(&_Rawmilk.TransactOpts, _tankId, _personInCharge, _qrCodeCID)
}

// UpdateMilkTank is a paid mutator transaction binding the contract method 0xef5d4f37.
//
// Solidity: function updateMilkTank(bytes32 _tankId, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkTransactor) UpdateMilkTank(opts *bind.TransactOpts, _tankId [32]byte, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "updateMilkTank", _tankId, _qrCodeCID)
}

// UpdateMilkTank is a paid mutator transaction binding the contract method 0xef5d4f37.
//
// Solidity: function updateMilkTank(bytes32 _tankId, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkSession) UpdateMilkTank(_tankId [32]byte, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.UpdateMilkTank(&_Rawmilk.TransactOpts, _tankId, _qrCodeCID)
}

// UpdateMilkTank is a paid mutator transaction binding the contract method 0xef5d4f37.
//
// Solidity: function updateMilkTank(bytes32 _tankId, string _qrCodeCID) returns()
func (_Rawmilk *RawmilkTransactorSession) UpdateMilkTank(_tankId [32]byte, _qrCodeCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.UpdateMilkTank(&_Rawmilk.TransactOpts, _tankId, _qrCodeCID)
}

// VerifyMilkQuality is a paid mutator transaction binding the contract method 0xeec34719.
//
// Solidity: function verifyMilkQuality(bytes32 _tankId, bool _approved, string _qualityReportCID) returns()
func (_Rawmilk *RawmilkTransactor) VerifyMilkQuality(opts *bind.TransactOpts, _tankId [32]byte, _approved bool, _qualityReportCID string) (*types.Transaction, error) {
	return _Rawmilk.contract.Transact(opts, "verifyMilkQuality", _tankId, _approved, _qualityReportCID)
}

// VerifyMilkQuality is a paid mutator transaction binding the contract method 0xeec34719.
//
// Solidity: function verifyMilkQuality(bytes32 _tankId, bool _approved, string _qualityReportCID) returns()
func (_Rawmilk *RawmilkSession) VerifyMilkQuality(_tankId [32]byte, _approved bool, _qualityReportCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.VerifyMilkQuality(&_Rawmilk.TransactOpts, _tankId, _approved, _qualityReportCID)
}

// VerifyMilkQuality is a paid mutator transaction binding the contract method 0xeec34719.
//
// Solidity: function verifyMilkQuality(bytes32 _tankId, bool _approved, string _qualityReportCID) returns()
func (_Rawmilk *RawmilkTransactorSession) VerifyMilkQuality(_tankId [32]byte, _approved bool, _qualityReportCID string) (*types.Transaction, error) {
	return _Rawmilk.Contract.VerifyMilkQuality(&_Rawmilk.TransactOpts, _tankId, _approved, _qualityReportCID)
}

// RawmilkMilkQualityVerifiedIterator is returned from FilterMilkQualityVerified and is used to iterate over the raw logs and unpacked data for MilkQualityVerified events raised by the Rawmilk contract.
type RawmilkMilkQualityVerifiedIterator struct {
	Event *RawmilkMilkQualityVerified // Event containing the contract specifics and raw log

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
func (it *RawmilkMilkQualityVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkMilkQualityVerified)
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
		it.Event = new(RawmilkMilkQualityVerified)
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
func (it *RawmilkMilkQualityVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkMilkQualityVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkMilkQualityVerified represents a MilkQualityVerified event raised by the Rawmilk contract.
type RawmilkMilkQualityVerified struct {
	TankId           [32]byte
	Status           uint8
	QualityReportCID string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterMilkQualityVerified is a free log retrieval operation binding the contract event 0xf8e56d2671d841b96def7d6bcc2753dd3030266f946afb9fbe2849d106e3ade7.
//
// Solidity: event MilkQualityVerified(bytes32 indexed tankId, uint8 status, string qualityReportCID)
func (_Rawmilk *RawmilkFilterer) FilterMilkQualityVerified(opts *bind.FilterOpts, tankId [][32]byte) (*RawmilkMilkQualityVerifiedIterator, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "MilkQualityVerified", tankIdRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkMilkQualityVerifiedIterator{contract: _Rawmilk.contract, event: "MilkQualityVerified", logs: logs, sub: sub}, nil
}

// WatchMilkQualityVerified is a free log subscription operation binding the contract event 0xf8e56d2671d841b96def7d6bcc2753dd3030266f946afb9fbe2849d106e3ade7.
//
// Solidity: event MilkQualityVerified(bytes32 indexed tankId, uint8 status, string qualityReportCID)
func (_Rawmilk *RawmilkFilterer) WatchMilkQualityVerified(opts *bind.WatchOpts, sink chan<- *RawmilkMilkQualityVerified, tankId [][32]byte) (event.Subscription, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "MilkQualityVerified", tankIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkMilkQualityVerified)
				if err := _Rawmilk.contract.UnpackLog(event, "MilkQualityVerified", log); err != nil {
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

// ParseMilkQualityVerified is a log parse operation binding the contract event 0xf8e56d2671d841b96def7d6bcc2753dd3030266f946afb9fbe2849d106e3ade7.
//
// Solidity: event MilkQualityVerified(bytes32 indexed tankId, uint8 status, string qualityReportCID)
func (_Rawmilk *RawmilkFilterer) ParseMilkQualityVerified(log types.Log) (*RawmilkMilkQualityVerified, error) {
	event := new(RawmilkMilkQualityVerified)
	if err := _Rawmilk.contract.UnpackLog(event, "MilkQualityVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RawmilkMilkTankCreatedIterator is returned from FilterMilkTankCreated and is used to iterate over the raw logs and unpacked data for MilkTankCreated events raised by the Rawmilk contract.
type RawmilkMilkTankCreatedIterator struct {
	Event *RawmilkMilkTankCreated // Event containing the contract specifics and raw log

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
func (it *RawmilkMilkTankCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkMilkTankCreated)
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
		it.Event = new(RawmilkMilkTankCreated)
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
func (it *RawmilkMilkTankCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkMilkTankCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkMilkTankCreated represents a MilkTankCreated event raised by the Rawmilk contract.
type RawmilkMilkTankCreated struct {
	TankId [32]byte
	Farmer common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMilkTankCreated is a free log retrieval operation binding the contract event 0x37b0128d4b624d18911b44dd0342b16405d37240c96673cf8d9b578307c45599.
//
// Solidity: event MilkTankCreated(bytes32 indexed tankId, address indexed farmer)
func (_Rawmilk *RawmilkFilterer) FilterMilkTankCreated(opts *bind.FilterOpts, tankId [][32]byte, farmer []common.Address) (*RawmilkMilkTankCreatedIterator, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}
	var farmerRule []interface{}
	for _, farmerItem := range farmer {
		farmerRule = append(farmerRule, farmerItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "MilkTankCreated", tankIdRule, farmerRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkMilkTankCreatedIterator{contract: _Rawmilk.contract, event: "MilkTankCreated", logs: logs, sub: sub}, nil
}

// WatchMilkTankCreated is a free log subscription operation binding the contract event 0x37b0128d4b624d18911b44dd0342b16405d37240c96673cf8d9b578307c45599.
//
// Solidity: event MilkTankCreated(bytes32 indexed tankId, address indexed farmer)
func (_Rawmilk *RawmilkFilterer) WatchMilkTankCreated(opts *bind.WatchOpts, sink chan<- *RawmilkMilkTankCreated, tankId [][32]byte, farmer []common.Address) (event.Subscription, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}
	var farmerRule []interface{}
	for _, farmerItem := range farmer {
		farmerRule = append(farmerRule, farmerItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "MilkTankCreated", tankIdRule, farmerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkMilkTankCreated)
				if err := _Rawmilk.contract.UnpackLog(event, "MilkTankCreated", log); err != nil {
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

// ParseMilkTankCreated is a log parse operation binding the contract event 0x37b0128d4b624d18911b44dd0342b16405d37240c96673cf8d9b578307c45599.
//
// Solidity: event MilkTankCreated(bytes32 indexed tankId, address indexed farmer)
func (_Rawmilk *RawmilkFilterer) ParseMilkTankCreated(log types.Log) (*RawmilkMilkTankCreated, error) {
	event := new(RawmilkMilkTankCreated)
	if err := _Rawmilk.contract.UnpackLog(event, "MilkTankCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RawmilkMilkTankUpdatedIterator is returned from FilterMilkTankUpdated and is used to iterate over the raw logs and unpacked data for MilkTankUpdated events raised by the Rawmilk contract.
type RawmilkMilkTankUpdatedIterator struct {
	Event *RawmilkMilkTankUpdated // Event containing the contract specifics and raw log

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
func (it *RawmilkMilkTankUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RawmilkMilkTankUpdated)
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
		it.Event = new(RawmilkMilkTankUpdated)
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
func (it *RawmilkMilkTankUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RawmilkMilkTankUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RawmilkMilkTankUpdated represents a MilkTankUpdated event raised by the Rawmilk contract.
type RawmilkMilkTankUpdated struct {
	TankId [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMilkTankUpdated is a free log retrieval operation binding the contract event 0x8636dbfe7028ccd3c14e5c0a2144b02d540e202d580b12552583f35e25e94ce6.
//
// Solidity: event MilkTankUpdated(bytes32 indexed tankId)
func (_Rawmilk *RawmilkFilterer) FilterMilkTankUpdated(opts *bind.FilterOpts, tankId [][32]byte) (*RawmilkMilkTankUpdatedIterator, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}

	logs, sub, err := _Rawmilk.contract.FilterLogs(opts, "MilkTankUpdated", tankIdRule)
	if err != nil {
		return nil, err
	}
	return &RawmilkMilkTankUpdatedIterator{contract: _Rawmilk.contract, event: "MilkTankUpdated", logs: logs, sub: sub}, nil
}

// WatchMilkTankUpdated is a free log subscription operation binding the contract event 0x8636dbfe7028ccd3c14e5c0a2144b02d540e202d580b12552583f35e25e94ce6.
//
// Solidity: event MilkTankUpdated(bytes32 indexed tankId)
func (_Rawmilk *RawmilkFilterer) WatchMilkTankUpdated(opts *bind.WatchOpts, sink chan<- *RawmilkMilkTankUpdated, tankId [][32]byte) (event.Subscription, error) {

	var tankIdRule []interface{}
	for _, tankIdItem := range tankId {
		tankIdRule = append(tankIdRule, tankIdItem)
	}

	logs, sub, err := _Rawmilk.contract.WatchLogs(opts, "MilkTankUpdated", tankIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RawmilkMilkTankUpdated)
				if err := _Rawmilk.contract.UnpackLog(event, "MilkTankUpdated", log); err != nil {
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

// ParseMilkTankUpdated is a log parse operation binding the contract event 0x8636dbfe7028ccd3c14e5c0a2144b02d540e202d580b12552583f35e25e94ce6.
//
// Solidity: event MilkTankUpdated(bytes32 indexed tankId)
func (_Rawmilk *RawmilkFilterer) ParseMilkTankUpdated(log types.Log) (*RawmilkMilkTankUpdated, error) {
	event := new(RawmilkMilkTankUpdated)
	if err := _Rawmilk.contract.UnpackLog(event, "MilkTankUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
