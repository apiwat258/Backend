// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package productlot

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

// ProductLotProductLotInfo is an auto generated low-level Go binding around an user-defined struct.
type ProductLotProductLotInfo struct {
	LotId                  [32]byte
	ProductId              [32]byte
	Factory                common.Address
	Inspector              string
	InspectionDate         *big.Int
	Grade                  bool
	QualityAndNutritionCID string
	MilkTankIds            [][32]byte
	Status                 uint8
}

// ProductlotMetaData contains all meta data concerning the Productlot contract.
var ProductlotMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rawMilkContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_productContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"lotId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"inspector\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"inspectionDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"grade\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"qualityAndNutritionCID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"milkTankIds\",\"type\":\"bytes32[]\"}],\"name\":\"ProductLotCreated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"productContract\",\"outputs\":[{\"internalType\":\"contractProduct\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"productLotIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"productLots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"lotId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"inspector\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"inspectionDate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"grade\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"qualityAndNutritionCID\",\"type\":\"string\"},{\"internalType\":\"enumProductLot.ProductLotStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"rawMilkContract\",\"outputs\":[{\"internalType\":\"contractRawMilk\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"trackingContract\",\"outputs\":[{\"internalType\":\"contractTracking\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"userRegistry\",\"outputs\":[{\"internalType\":\"contractUserRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_lotId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_productId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_inspector\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_grade\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"_qualityAndNutritionCID\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"_milkTankIds\",\"type\":\"bytes32[]\"}],\"name\":\"createProductLot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_lotId\",\"type\":\"bytes32\"}],\"name\":\"getProductLot\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"lotId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"inspector\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"inspectionDate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"grade\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"qualityAndNutritionCID\",\"type\":\"string\"},{\"internalType\":\"bytes32[]\",\"name\":\"milkTankIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"enumProductLot.ProductLotStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structProductLot.ProductLotInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_lotId\",\"type\":\"bytes32\"}],\"name\":\"isProductLotExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"getProductLotsByFactory\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_lotId\",\"type\":\"bytes32\"}],\"name\":\"getMilkTanksByProductLot\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_trackingContract\",\"type\":\"address\"}],\"name\":\"setTrackingContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_lotId\",\"type\":\"bytes32\"}],\"name\":\"updateProductLotStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ProductlotABI is the input ABI used to generate the binding from.
// Deprecated: Use ProductlotMetaData.ABI instead.
var ProductlotABI = ProductlotMetaData.ABI

// Productlot is an auto generated Go binding around an Ethereum contract.
type Productlot struct {
	ProductlotCaller     // Read-only binding to the contract
	ProductlotTransactor // Write-only binding to the contract
	ProductlotFilterer   // Log filterer for contract events
}

// ProductlotCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProductlotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductlotTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProductlotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductlotFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProductlotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductlotSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProductlotSession struct {
	Contract     *Productlot       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProductlotCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProductlotCallerSession struct {
	Contract *ProductlotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ProductlotTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProductlotTransactorSession struct {
	Contract     *ProductlotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ProductlotRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProductlotRaw struct {
	Contract *Productlot // Generic contract binding to access the raw methods on
}

// ProductlotCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProductlotCallerRaw struct {
	Contract *ProductlotCaller // Generic read-only contract binding to access the raw methods on
}

// ProductlotTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProductlotTransactorRaw struct {
	Contract *ProductlotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProductlot creates a new instance of Productlot, bound to a specific deployed contract.
func NewProductlot(address common.Address, backend bind.ContractBackend) (*Productlot, error) {
	contract, err := bindProductlot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Productlot{ProductlotCaller: ProductlotCaller{contract: contract}, ProductlotTransactor: ProductlotTransactor{contract: contract}, ProductlotFilterer: ProductlotFilterer{contract: contract}}, nil
}

// NewProductlotCaller creates a new read-only instance of Productlot, bound to a specific deployed contract.
func NewProductlotCaller(address common.Address, caller bind.ContractCaller) (*ProductlotCaller, error) {
	contract, err := bindProductlot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProductlotCaller{contract: contract}, nil
}

// NewProductlotTransactor creates a new write-only instance of Productlot, bound to a specific deployed contract.
func NewProductlotTransactor(address common.Address, transactor bind.ContractTransactor) (*ProductlotTransactor, error) {
	contract, err := bindProductlot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProductlotTransactor{contract: contract}, nil
}

// NewProductlotFilterer creates a new log filterer instance of Productlot, bound to a specific deployed contract.
func NewProductlotFilterer(address common.Address, filterer bind.ContractFilterer) (*ProductlotFilterer, error) {
	contract, err := bindProductlot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProductlotFilterer{contract: contract}, nil
}

// bindProductlot binds a generic wrapper to an already deployed contract.
func bindProductlot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProductlotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Productlot *ProductlotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Productlot.Contract.ProductlotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Productlot *ProductlotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Productlot.Contract.ProductlotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Productlot *ProductlotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Productlot.Contract.ProductlotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Productlot *ProductlotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Productlot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Productlot *ProductlotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Productlot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Productlot *ProductlotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Productlot.Contract.contract.Transact(opts, method, params...)
}

// GetMilkTanksByProductLot is a free data retrieval call binding the contract method 0x718a947c.
//
// Solidity: function getMilkTanksByProductLot(bytes32 _lotId) view returns(bytes32[])
func (_Productlot *ProductlotCaller) GetMilkTanksByProductLot(opts *bind.CallOpts, _lotId [32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "getMilkTanksByProductLot", _lotId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetMilkTanksByProductLot is a free data retrieval call binding the contract method 0x718a947c.
//
// Solidity: function getMilkTanksByProductLot(bytes32 _lotId) view returns(bytes32[])
func (_Productlot *ProductlotSession) GetMilkTanksByProductLot(_lotId [32]byte) ([][32]byte, error) {
	return _Productlot.Contract.GetMilkTanksByProductLot(&_Productlot.CallOpts, _lotId)
}

// GetMilkTanksByProductLot is a free data retrieval call binding the contract method 0x718a947c.
//
// Solidity: function getMilkTanksByProductLot(bytes32 _lotId) view returns(bytes32[])
func (_Productlot *ProductlotCallerSession) GetMilkTanksByProductLot(_lotId [32]byte) ([][32]byte, error) {
	return _Productlot.Contract.GetMilkTanksByProductLot(&_Productlot.CallOpts, _lotId)
}

// GetProductLot is a free data retrieval call binding the contract method 0xd794e362.
//
// Solidity: function getProductLot(bytes32 _lotId) view returns((bytes32,bytes32,address,string,uint256,bool,string,bytes32[],uint8))
func (_Productlot *ProductlotCaller) GetProductLot(opts *bind.CallOpts, _lotId [32]byte) (ProductLotProductLotInfo, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "getProductLot", _lotId)

	if err != nil {
		return *new(ProductLotProductLotInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ProductLotProductLotInfo)).(*ProductLotProductLotInfo)

	return out0, err

}

// GetProductLot is a free data retrieval call binding the contract method 0xd794e362.
//
// Solidity: function getProductLot(bytes32 _lotId) view returns((bytes32,bytes32,address,string,uint256,bool,string,bytes32[],uint8))
func (_Productlot *ProductlotSession) GetProductLot(_lotId [32]byte) (ProductLotProductLotInfo, error) {
	return _Productlot.Contract.GetProductLot(&_Productlot.CallOpts, _lotId)
}

// GetProductLot is a free data retrieval call binding the contract method 0xd794e362.
//
// Solidity: function getProductLot(bytes32 _lotId) view returns((bytes32,bytes32,address,string,uint256,bool,string,bytes32[],uint8))
func (_Productlot *ProductlotCallerSession) GetProductLot(_lotId [32]byte) (ProductLotProductLotInfo, error) {
	return _Productlot.Contract.GetProductLot(&_Productlot.CallOpts, _lotId)
}

// GetProductLotsByFactory is a free data retrieval call binding the contract method 0x01b2a86f.
//
// Solidity: function getProductLotsByFactory(address _factory) view returns(bytes32[])
func (_Productlot *ProductlotCaller) GetProductLotsByFactory(opts *bind.CallOpts, _factory common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "getProductLotsByFactory", _factory)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetProductLotsByFactory is a free data retrieval call binding the contract method 0x01b2a86f.
//
// Solidity: function getProductLotsByFactory(address _factory) view returns(bytes32[])
func (_Productlot *ProductlotSession) GetProductLotsByFactory(_factory common.Address) ([][32]byte, error) {
	return _Productlot.Contract.GetProductLotsByFactory(&_Productlot.CallOpts, _factory)
}

// GetProductLotsByFactory is a free data retrieval call binding the contract method 0x01b2a86f.
//
// Solidity: function getProductLotsByFactory(address _factory) view returns(bytes32[])
func (_Productlot *ProductlotCallerSession) GetProductLotsByFactory(_factory common.Address) ([][32]byte, error) {
	return _Productlot.Contract.GetProductLotsByFactory(&_Productlot.CallOpts, _factory)
}

// IsProductLotExists is a free data retrieval call binding the contract method 0xfa0bff30.
//
// Solidity: function isProductLotExists(bytes32 _lotId) view returns(bool)
func (_Productlot *ProductlotCaller) IsProductLotExists(opts *bind.CallOpts, _lotId [32]byte) (bool, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "isProductLotExists", _lotId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsProductLotExists is a free data retrieval call binding the contract method 0xfa0bff30.
//
// Solidity: function isProductLotExists(bytes32 _lotId) view returns(bool)
func (_Productlot *ProductlotSession) IsProductLotExists(_lotId [32]byte) (bool, error) {
	return _Productlot.Contract.IsProductLotExists(&_Productlot.CallOpts, _lotId)
}

// IsProductLotExists is a free data retrieval call binding the contract method 0xfa0bff30.
//
// Solidity: function isProductLotExists(bytes32 _lotId) view returns(bool)
func (_Productlot *ProductlotCallerSession) IsProductLotExists(_lotId [32]byte) (bool, error) {
	return _Productlot.Contract.IsProductLotExists(&_Productlot.CallOpts, _lotId)
}

// ProductContract is a free data retrieval call binding the contract method 0x2a28d4df.
//
// Solidity: function productContract() view returns(address)
func (_Productlot *ProductlotCaller) ProductContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "productContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProductContract is a free data retrieval call binding the contract method 0x2a28d4df.
//
// Solidity: function productContract() view returns(address)
func (_Productlot *ProductlotSession) ProductContract() (common.Address, error) {
	return _Productlot.Contract.ProductContract(&_Productlot.CallOpts)
}

// ProductContract is a free data retrieval call binding the contract method 0x2a28d4df.
//
// Solidity: function productContract() view returns(address)
func (_Productlot *ProductlotCallerSession) ProductContract() (common.Address, error) {
	return _Productlot.Contract.ProductContract(&_Productlot.CallOpts)
}

// ProductLotIds is a free data retrieval call binding the contract method 0xa7a8e53e.
//
// Solidity: function productLotIds(uint256 ) view returns(bytes32)
func (_Productlot *ProductlotCaller) ProductLotIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "productLotIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProductLotIds is a free data retrieval call binding the contract method 0xa7a8e53e.
//
// Solidity: function productLotIds(uint256 ) view returns(bytes32)
func (_Productlot *ProductlotSession) ProductLotIds(arg0 *big.Int) ([32]byte, error) {
	return _Productlot.Contract.ProductLotIds(&_Productlot.CallOpts, arg0)
}

// ProductLotIds is a free data retrieval call binding the contract method 0xa7a8e53e.
//
// Solidity: function productLotIds(uint256 ) view returns(bytes32)
func (_Productlot *ProductlotCallerSession) ProductLotIds(arg0 *big.Int) ([32]byte, error) {
	return _Productlot.Contract.ProductLotIds(&_Productlot.CallOpts, arg0)
}

// ProductLots is a free data retrieval call binding the contract method 0xd5729549.
//
// Solidity: function productLots(bytes32 ) view returns(bytes32 lotId, bytes32 productId, address factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, uint8 status)
func (_Productlot *ProductlotCaller) ProductLots(opts *bind.CallOpts, arg0 [32]byte) (struct {
	LotId                  [32]byte
	ProductId              [32]byte
	Factory                common.Address
	Inspector              string
	InspectionDate         *big.Int
	Grade                  bool
	QualityAndNutritionCID string
	Status                 uint8
}, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "productLots", arg0)

	outstruct := new(struct {
		LotId                  [32]byte
		ProductId              [32]byte
		Factory                common.Address
		Inspector              string
		InspectionDate         *big.Int
		Grade                  bool
		QualityAndNutritionCID string
		Status                 uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LotId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ProductId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Factory = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Inspector = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.InspectionDate = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Grade = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.QualityAndNutritionCID = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

// ProductLots is a free data retrieval call binding the contract method 0xd5729549.
//
// Solidity: function productLots(bytes32 ) view returns(bytes32 lotId, bytes32 productId, address factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, uint8 status)
func (_Productlot *ProductlotSession) ProductLots(arg0 [32]byte) (struct {
	LotId                  [32]byte
	ProductId              [32]byte
	Factory                common.Address
	Inspector              string
	InspectionDate         *big.Int
	Grade                  bool
	QualityAndNutritionCID string
	Status                 uint8
}, error) {
	return _Productlot.Contract.ProductLots(&_Productlot.CallOpts, arg0)
}

// ProductLots is a free data retrieval call binding the contract method 0xd5729549.
//
// Solidity: function productLots(bytes32 ) view returns(bytes32 lotId, bytes32 productId, address factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, uint8 status)
func (_Productlot *ProductlotCallerSession) ProductLots(arg0 [32]byte) (struct {
	LotId                  [32]byte
	ProductId              [32]byte
	Factory                common.Address
	Inspector              string
	InspectionDate         *big.Int
	Grade                  bool
	QualityAndNutritionCID string
	Status                 uint8
}, error) {
	return _Productlot.Contract.ProductLots(&_Productlot.CallOpts, arg0)
}

// RawMilkContract is a free data retrieval call binding the contract method 0x1ce95d30.
//
// Solidity: function rawMilkContract() view returns(address)
func (_Productlot *ProductlotCaller) RawMilkContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "rawMilkContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RawMilkContract is a free data retrieval call binding the contract method 0x1ce95d30.
//
// Solidity: function rawMilkContract() view returns(address)
func (_Productlot *ProductlotSession) RawMilkContract() (common.Address, error) {
	return _Productlot.Contract.RawMilkContract(&_Productlot.CallOpts)
}

// RawMilkContract is a free data retrieval call binding the contract method 0x1ce95d30.
//
// Solidity: function rawMilkContract() view returns(address)
func (_Productlot *ProductlotCallerSession) RawMilkContract() (common.Address, error) {
	return _Productlot.Contract.RawMilkContract(&_Productlot.CallOpts)
}

// TrackingContract is a free data retrieval call binding the contract method 0x8c91be92.
//
// Solidity: function trackingContract() view returns(address)
func (_Productlot *ProductlotCaller) TrackingContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "trackingContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrackingContract is a free data retrieval call binding the contract method 0x8c91be92.
//
// Solidity: function trackingContract() view returns(address)
func (_Productlot *ProductlotSession) TrackingContract() (common.Address, error) {
	return _Productlot.Contract.TrackingContract(&_Productlot.CallOpts)
}

// TrackingContract is a free data retrieval call binding the contract method 0x8c91be92.
//
// Solidity: function trackingContract() view returns(address)
func (_Productlot *ProductlotCallerSession) TrackingContract() (common.Address, error) {
	return _Productlot.Contract.TrackingContract(&_Productlot.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Productlot *ProductlotCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Productlot.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Productlot *ProductlotSession) UserRegistry() (common.Address, error) {
	return _Productlot.Contract.UserRegistry(&_Productlot.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Productlot *ProductlotCallerSession) UserRegistry() (common.Address, error) {
	return _Productlot.Contract.UserRegistry(&_Productlot.CallOpts)
}

// CreateProductLot is a paid mutator transaction binding the contract method 0x43f9ad45.
//
// Solidity: function createProductLot(bytes32 _lotId, bytes32 _productId, string _inspector, bool _grade, string _qualityAndNutritionCID, bytes32[] _milkTankIds) returns()
func (_Productlot *ProductlotTransactor) CreateProductLot(opts *bind.TransactOpts, _lotId [32]byte, _productId [32]byte, _inspector string, _grade bool, _qualityAndNutritionCID string, _milkTankIds [][32]byte) (*types.Transaction, error) {
	return _Productlot.contract.Transact(opts, "createProductLot", _lotId, _productId, _inspector, _grade, _qualityAndNutritionCID, _milkTankIds)
}

// CreateProductLot is a paid mutator transaction binding the contract method 0x43f9ad45.
//
// Solidity: function createProductLot(bytes32 _lotId, bytes32 _productId, string _inspector, bool _grade, string _qualityAndNutritionCID, bytes32[] _milkTankIds) returns()
func (_Productlot *ProductlotSession) CreateProductLot(_lotId [32]byte, _productId [32]byte, _inspector string, _grade bool, _qualityAndNutritionCID string, _milkTankIds [][32]byte) (*types.Transaction, error) {
	return _Productlot.Contract.CreateProductLot(&_Productlot.TransactOpts, _lotId, _productId, _inspector, _grade, _qualityAndNutritionCID, _milkTankIds)
}

// CreateProductLot is a paid mutator transaction binding the contract method 0x43f9ad45.
//
// Solidity: function createProductLot(bytes32 _lotId, bytes32 _productId, string _inspector, bool _grade, string _qualityAndNutritionCID, bytes32[] _milkTankIds) returns()
func (_Productlot *ProductlotTransactorSession) CreateProductLot(_lotId [32]byte, _productId [32]byte, _inspector string, _grade bool, _qualityAndNutritionCID string, _milkTankIds [][32]byte) (*types.Transaction, error) {
	return _Productlot.Contract.CreateProductLot(&_Productlot.TransactOpts, _lotId, _productId, _inspector, _grade, _qualityAndNutritionCID, _milkTankIds)
}

// SetTrackingContract is a paid mutator transaction binding the contract method 0x42dc43b0.
//
// Solidity: function setTrackingContract(address _trackingContract) returns()
func (_Productlot *ProductlotTransactor) SetTrackingContract(opts *bind.TransactOpts, _trackingContract common.Address) (*types.Transaction, error) {
	return _Productlot.contract.Transact(opts, "setTrackingContract", _trackingContract)
}

// SetTrackingContract is a paid mutator transaction binding the contract method 0x42dc43b0.
//
// Solidity: function setTrackingContract(address _trackingContract) returns()
func (_Productlot *ProductlotSession) SetTrackingContract(_trackingContract common.Address) (*types.Transaction, error) {
	return _Productlot.Contract.SetTrackingContract(&_Productlot.TransactOpts, _trackingContract)
}

// SetTrackingContract is a paid mutator transaction binding the contract method 0x42dc43b0.
//
// Solidity: function setTrackingContract(address _trackingContract) returns()
func (_Productlot *ProductlotTransactorSession) SetTrackingContract(_trackingContract common.Address) (*types.Transaction, error) {
	return _Productlot.Contract.SetTrackingContract(&_Productlot.TransactOpts, _trackingContract)
}

// UpdateProductLotStatus is a paid mutator transaction binding the contract method 0x52df0ee6.
//
// Solidity: function updateProductLotStatus(bytes32 _lotId) returns()
func (_Productlot *ProductlotTransactor) UpdateProductLotStatus(opts *bind.TransactOpts, _lotId [32]byte) (*types.Transaction, error) {
	return _Productlot.contract.Transact(opts, "updateProductLotStatus", _lotId)
}

// UpdateProductLotStatus is a paid mutator transaction binding the contract method 0x52df0ee6.
//
// Solidity: function updateProductLotStatus(bytes32 _lotId) returns()
func (_Productlot *ProductlotSession) UpdateProductLotStatus(_lotId [32]byte) (*types.Transaction, error) {
	return _Productlot.Contract.UpdateProductLotStatus(&_Productlot.TransactOpts, _lotId)
}

// UpdateProductLotStatus is a paid mutator transaction binding the contract method 0x52df0ee6.
//
// Solidity: function updateProductLotStatus(bytes32 _lotId) returns()
func (_Productlot *ProductlotTransactorSession) UpdateProductLotStatus(_lotId [32]byte) (*types.Transaction, error) {
	return _Productlot.Contract.UpdateProductLotStatus(&_Productlot.TransactOpts, _lotId)
}

// ProductlotProductLotCreatedIterator is returned from FilterProductLotCreated and is used to iterate over the raw logs and unpacked data for ProductLotCreated events raised by the Productlot contract.
type ProductlotProductLotCreatedIterator struct {
	Event *ProductlotProductLotCreated // Event containing the contract specifics and raw log

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
func (it *ProductlotProductLotCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProductlotProductLotCreated)
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
		it.Event = new(ProductlotProductLotCreated)
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
func (it *ProductlotProductLotCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProductlotProductLotCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProductlotProductLotCreated represents a ProductLotCreated event raised by the Productlot contract.
type ProductlotProductLotCreated struct {
	LotId                  [32]byte
	ProductId              [32]byte
	Factory                common.Address
	Inspector              string
	InspectionDate         *big.Int
	Grade                  bool
	QualityAndNutritionCID string
	MilkTankIds            [][32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterProductLotCreated is a free log retrieval operation binding the contract event 0x23013b10378a942712c02f38e17299de4d74256e5007665263df05dd00827c44.
//
// Solidity: event ProductLotCreated(bytes32 indexed lotId, bytes32 indexed productId, address indexed factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, bytes32[] milkTankIds)
func (_Productlot *ProductlotFilterer) FilterProductLotCreated(opts *bind.FilterOpts, lotId [][32]byte, productId [][32]byte, factory []common.Address) (*ProductlotProductLotCreatedIterator, error) {

	var lotIdRule []interface{}
	for _, lotIdItem := range lotId {
		lotIdRule = append(lotIdRule, lotIdItem)
	}
	var productIdRule []interface{}
	for _, productIdItem := range productId {
		productIdRule = append(productIdRule, productIdItem)
	}
	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _Productlot.contract.FilterLogs(opts, "ProductLotCreated", lotIdRule, productIdRule, factoryRule)
	if err != nil {
		return nil, err
	}
	return &ProductlotProductLotCreatedIterator{contract: _Productlot.contract, event: "ProductLotCreated", logs: logs, sub: sub}, nil
}

// WatchProductLotCreated is a free log subscription operation binding the contract event 0x23013b10378a942712c02f38e17299de4d74256e5007665263df05dd00827c44.
//
// Solidity: event ProductLotCreated(bytes32 indexed lotId, bytes32 indexed productId, address indexed factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, bytes32[] milkTankIds)
func (_Productlot *ProductlotFilterer) WatchProductLotCreated(opts *bind.WatchOpts, sink chan<- *ProductlotProductLotCreated, lotId [][32]byte, productId [][32]byte, factory []common.Address) (event.Subscription, error) {

	var lotIdRule []interface{}
	for _, lotIdItem := range lotId {
		lotIdRule = append(lotIdRule, lotIdItem)
	}
	var productIdRule []interface{}
	for _, productIdItem := range productId {
		productIdRule = append(productIdRule, productIdItem)
	}
	var factoryRule []interface{}
	for _, factoryItem := range factory {
		factoryRule = append(factoryRule, factoryItem)
	}

	logs, sub, err := _Productlot.contract.WatchLogs(opts, "ProductLotCreated", lotIdRule, productIdRule, factoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProductlotProductLotCreated)
				if err := _Productlot.contract.UnpackLog(event, "ProductLotCreated", log); err != nil {
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

// ParseProductLotCreated is a log parse operation binding the contract event 0x23013b10378a942712c02f38e17299de4d74256e5007665263df05dd00827c44.
//
// Solidity: event ProductLotCreated(bytes32 indexed lotId, bytes32 indexed productId, address indexed factory, string inspector, uint256 inspectionDate, bool grade, string qualityAndNutritionCID, bytes32[] milkTankIds)
func (_Productlot *ProductlotFilterer) ParseProductLotCreated(log types.Log) (*ProductlotProductLotCreated, error) {
	event := new(ProductlotProductLotCreated)
	if err := _Productlot.contract.UnpackLog(event, "ProductLotCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
