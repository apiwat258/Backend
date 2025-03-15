// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tracking

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

// TrackingLogisticsCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type TrackingLogisticsCheckpoint struct {
	TrackingId        [32]byte
	LogisticsProvider common.Address
	PickupTime        *big.Int
	DeliveryTime      *big.Int
	Quantity          *big.Int
	Temperature       *big.Int
	PersonInCharge    string
	CheckType         uint8
	ReceiverCID       string
}

// TrackingRetailerConfirmation is an auto generated low-level Go binding around an user-defined struct.
type TrackingRetailerConfirmation struct {
	TrackingId     [32]byte
	RetailerId     string
	ReceivedTime   *big.Int
	QualityCID     string
	PersonInCharge string
}

// TrackingTrackingEvent is an auto generated low-level Go binding around an user-defined struct.
type TrackingTrackingEvent struct {
	TrackingId   [32]byte
	ProductLotId [32]byte
	RetailerId   string
	QrCodeCID    string
	Status       uint8
}

// TrackingMetaData contains all meta data concerning the Tracking contract.
var TrackingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_productLotContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"name\":\"LogisticsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"qualityCID\",\"type\":\"string\"}],\"name\":\"RetailerReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"productLotId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"qrCodeCID\",\"type\":\"string\"}],\"name\":\"TrackingCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allTrackingIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"logisticsCheckpoints\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"productLotContract\",\"outputs\":[{\"internalType\":\"contractProductLot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"retailerConfirmations\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"receivedTime\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"qualityCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"retailerTracking\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"trackingEvents\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"productLotId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"qrCodeCID\",\"type\":\"string\"},{\"internalType\":\"enumTracking.TrackingStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"userRegistry\",\"outputs\":[{\"internalType\":\"contractUserRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_productLotId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_retailerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_qrCodeCID\",\"type\":\"string\"}],\"name\":\"createTrackingEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"_temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"_personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"_checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_receiverCID\",\"type\":\"string\"}],\"name\":\"updateLogisticsCheckpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_retailerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_qualityCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_personInCharge\",\"type\":\"string\"}],\"name\":\"retailerReceiveProduct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_trackingId\",\"type\":\"bytes32\"}],\"name\":\"getTrackingById\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"productLotId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"qrCodeCID\",\"type\":\"string\"},{\"internalType\":\"enumTracking.TrackingStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structTracking.TrackingEvent\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"internalType\":\"structTracking.LogisticsCheckpoint[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"retailerId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"receivedTime\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"qualityCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"}],\"internalType\":\"structTracking.RetailerConfirmation\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_retailerId\",\"type\":\"string\"}],\"name\":\"getTrackingByRetailer\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_productLotId\",\"type\":\"bytes32\"}],\"name\":\"getTrackingByLotId\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"resultTrackingIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"string[]\",\"name\":\"retailerIds\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"qrCodeCIDs\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"getAllTrackingIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_trackingId\",\"type\":\"bytes32\"}],\"name\":\"getLogisticsCheckpointsByTrackingId\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"internalType\":\"structTracking.LogisticsCheckpoint[]\",\"name\":\"beforeCheckpoints\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"internalType\":\"structTracking.LogisticsCheckpoint[]\",\"name\":\"duringCheckpoints\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"trackingId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"logisticsProvider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"pickupTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deliveryTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"temperature\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"personInCharge\",\"type\":\"string\"},{\"internalType\":\"enumTracking.LogisticsCheckType\",\"name\":\"checkType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"receiverCID\",\"type\":\"string\"}],\"internalType\":\"structTracking.LogisticsCheckpoint[]\",\"name\":\"afterCheckpoints\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true}]",
}

// TrackingABI is the input ABI used to generate the binding from.
// Deprecated: Use TrackingMetaData.ABI instead.
var TrackingABI = TrackingMetaData.ABI

// Tracking is an auto generated Go binding around an Ethereum contract.
type Tracking struct {
	TrackingCaller     // Read-only binding to the contract
	TrackingTransactor // Write-only binding to the contract
	TrackingFilterer   // Log filterer for contract events
}

// TrackingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TrackingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrackingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TrackingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrackingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TrackingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrackingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TrackingSession struct {
	Contract     *Tracking         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TrackingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TrackingCallerSession struct {
	Contract *TrackingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TrackingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TrackingTransactorSession struct {
	Contract     *TrackingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TrackingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TrackingRaw struct {
	Contract *Tracking // Generic contract binding to access the raw methods on
}

// TrackingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TrackingCallerRaw struct {
	Contract *TrackingCaller // Generic read-only contract binding to access the raw methods on
}

// TrackingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TrackingTransactorRaw struct {
	Contract *TrackingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTracking creates a new instance of Tracking, bound to a specific deployed contract.
func NewTracking(address common.Address, backend bind.ContractBackend) (*Tracking, error) {
	contract, err := bindTracking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Tracking{TrackingCaller: TrackingCaller{contract: contract}, TrackingTransactor: TrackingTransactor{contract: contract}, TrackingFilterer: TrackingFilterer{contract: contract}}, nil
}

// NewTrackingCaller creates a new read-only instance of Tracking, bound to a specific deployed contract.
func NewTrackingCaller(address common.Address, caller bind.ContractCaller) (*TrackingCaller, error) {
	contract, err := bindTracking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrackingCaller{contract: contract}, nil
}

// NewTrackingTransactor creates a new write-only instance of Tracking, bound to a specific deployed contract.
func NewTrackingTransactor(address common.Address, transactor bind.ContractTransactor) (*TrackingTransactor, error) {
	contract, err := bindTracking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrackingTransactor{contract: contract}, nil
}

// NewTrackingFilterer creates a new log filterer instance of Tracking, bound to a specific deployed contract.
func NewTrackingFilterer(address common.Address, filterer bind.ContractFilterer) (*TrackingFilterer, error) {
	contract, err := bindTracking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrackingFilterer{contract: contract}, nil
}

// bindTracking binds a generic wrapper to an already deployed contract.
func bindTracking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TrackingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Tracking *TrackingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Tracking.Contract.TrackingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Tracking *TrackingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tracking.Contract.TrackingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Tracking *TrackingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Tracking.Contract.TrackingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Tracking *TrackingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Tracking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Tracking *TrackingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Tracking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Tracking *TrackingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Tracking.Contract.contract.Transact(opts, method, params...)
}

// AllTrackingIds is a free data retrieval call binding the contract method 0x1a5c3d0e.
//
// Solidity: function allTrackingIds(uint256 ) view returns(bytes32)
func (_Tracking *TrackingCaller) AllTrackingIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "allTrackingIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AllTrackingIds is a free data retrieval call binding the contract method 0x1a5c3d0e.
//
// Solidity: function allTrackingIds(uint256 ) view returns(bytes32)
func (_Tracking *TrackingSession) AllTrackingIds(arg0 *big.Int) ([32]byte, error) {
	return _Tracking.Contract.AllTrackingIds(&_Tracking.CallOpts, arg0)
}

// AllTrackingIds is a free data retrieval call binding the contract method 0x1a5c3d0e.
//
// Solidity: function allTrackingIds(uint256 ) view returns(bytes32)
func (_Tracking *TrackingCallerSession) AllTrackingIds(arg0 *big.Int) ([32]byte, error) {
	return _Tracking.Contract.AllTrackingIds(&_Tracking.CallOpts, arg0)
}

// GetAllTrackingIds is a free data retrieval call binding the contract method 0x54cb24a4.
//
// Solidity: function getAllTrackingIds() view returns(bytes32[])
func (_Tracking *TrackingCaller) GetAllTrackingIds(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "getAllTrackingIds")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetAllTrackingIds is a free data retrieval call binding the contract method 0x54cb24a4.
//
// Solidity: function getAllTrackingIds() view returns(bytes32[])
func (_Tracking *TrackingSession) GetAllTrackingIds() ([][32]byte, error) {
	return _Tracking.Contract.GetAllTrackingIds(&_Tracking.CallOpts)
}

// GetAllTrackingIds is a free data retrieval call binding the contract method 0x54cb24a4.
//
// Solidity: function getAllTrackingIds() view returns(bytes32[])
func (_Tracking *TrackingCallerSession) GetAllTrackingIds() ([][32]byte, error) {
	return _Tracking.Contract.GetAllTrackingIds(&_Tracking.CallOpts)
}

// GetLogisticsCheckpointsByTrackingId is a free data retrieval call binding the contract method 0x56ef21b4.
//
// Solidity: function getLogisticsCheckpointsByTrackingId(bytes32 _trackingId) view returns((bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] beforeCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] duringCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] afterCheckpoints)
func (_Tracking *TrackingCaller) GetLogisticsCheckpointsByTrackingId(opts *bind.CallOpts, _trackingId [32]byte) (struct {
	BeforeCheckpoints []TrackingLogisticsCheckpoint
	DuringCheckpoints []TrackingLogisticsCheckpoint
	AfterCheckpoints  []TrackingLogisticsCheckpoint
}, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "getLogisticsCheckpointsByTrackingId", _trackingId)

	outstruct := new(struct {
		BeforeCheckpoints []TrackingLogisticsCheckpoint
		DuringCheckpoints []TrackingLogisticsCheckpoint
		AfterCheckpoints  []TrackingLogisticsCheckpoint
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BeforeCheckpoints = *abi.ConvertType(out[0], new([]TrackingLogisticsCheckpoint)).(*[]TrackingLogisticsCheckpoint)
	outstruct.DuringCheckpoints = *abi.ConvertType(out[1], new([]TrackingLogisticsCheckpoint)).(*[]TrackingLogisticsCheckpoint)
	outstruct.AfterCheckpoints = *abi.ConvertType(out[2], new([]TrackingLogisticsCheckpoint)).(*[]TrackingLogisticsCheckpoint)

	return *outstruct, err

}

// GetLogisticsCheckpointsByTrackingId is a free data retrieval call binding the contract method 0x56ef21b4.
//
// Solidity: function getLogisticsCheckpointsByTrackingId(bytes32 _trackingId) view returns((bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] beforeCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] duringCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] afterCheckpoints)
func (_Tracking *TrackingSession) GetLogisticsCheckpointsByTrackingId(_trackingId [32]byte) (struct {
	BeforeCheckpoints []TrackingLogisticsCheckpoint
	DuringCheckpoints []TrackingLogisticsCheckpoint
	AfterCheckpoints  []TrackingLogisticsCheckpoint
}, error) {
	return _Tracking.Contract.GetLogisticsCheckpointsByTrackingId(&_Tracking.CallOpts, _trackingId)
}

// GetLogisticsCheckpointsByTrackingId is a free data retrieval call binding the contract method 0x56ef21b4.
//
// Solidity: function getLogisticsCheckpointsByTrackingId(bytes32 _trackingId) view returns((bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] beforeCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] duringCheckpoints, (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[] afterCheckpoints)
func (_Tracking *TrackingCallerSession) GetLogisticsCheckpointsByTrackingId(_trackingId [32]byte) (struct {
	BeforeCheckpoints []TrackingLogisticsCheckpoint
	DuringCheckpoints []TrackingLogisticsCheckpoint
	AfterCheckpoints  []TrackingLogisticsCheckpoint
}, error) {
	return _Tracking.Contract.GetLogisticsCheckpointsByTrackingId(&_Tracking.CallOpts, _trackingId)
}

// GetTrackingById is a free data retrieval call binding the contract method 0xc93c35eb.
//
// Solidity: function getTrackingById(bytes32 _trackingId) view returns((bytes32,bytes32,string,string,uint8), (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[], (bytes32,string,uint256,string,string))
func (_Tracking *TrackingCaller) GetTrackingById(opts *bind.CallOpts, _trackingId [32]byte) (TrackingTrackingEvent, []TrackingLogisticsCheckpoint, TrackingRetailerConfirmation, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "getTrackingById", _trackingId)

	if err != nil {
		return *new(TrackingTrackingEvent), *new([]TrackingLogisticsCheckpoint), *new(TrackingRetailerConfirmation), err
	}

	out0 := *abi.ConvertType(out[0], new(TrackingTrackingEvent)).(*TrackingTrackingEvent)
	out1 := *abi.ConvertType(out[1], new([]TrackingLogisticsCheckpoint)).(*[]TrackingLogisticsCheckpoint)
	out2 := *abi.ConvertType(out[2], new(TrackingRetailerConfirmation)).(*TrackingRetailerConfirmation)

	return out0, out1, out2, err

}

// GetTrackingById is a free data retrieval call binding the contract method 0xc93c35eb.
//
// Solidity: function getTrackingById(bytes32 _trackingId) view returns((bytes32,bytes32,string,string,uint8), (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[], (bytes32,string,uint256,string,string))
func (_Tracking *TrackingSession) GetTrackingById(_trackingId [32]byte) (TrackingTrackingEvent, []TrackingLogisticsCheckpoint, TrackingRetailerConfirmation, error) {
	return _Tracking.Contract.GetTrackingById(&_Tracking.CallOpts, _trackingId)
}

// GetTrackingById is a free data retrieval call binding the contract method 0xc93c35eb.
//
// Solidity: function getTrackingById(bytes32 _trackingId) view returns((bytes32,bytes32,string,string,uint8), (bytes32,address,uint256,uint256,uint256,int256,string,uint8,string)[], (bytes32,string,uint256,string,string))
func (_Tracking *TrackingCallerSession) GetTrackingById(_trackingId [32]byte) (TrackingTrackingEvent, []TrackingLogisticsCheckpoint, TrackingRetailerConfirmation, error) {
	return _Tracking.Contract.GetTrackingById(&_Tracking.CallOpts, _trackingId)
}

// GetTrackingByLotId is a free data retrieval call binding the contract method 0x65f0f7d5.
//
// Solidity: function getTrackingByLotId(bytes32 _productLotId) view returns(bytes32[] resultTrackingIds, string[] retailerIds, string[] qrCodeCIDs)
func (_Tracking *TrackingCaller) GetTrackingByLotId(opts *bind.CallOpts, _productLotId [32]byte) (struct {
	ResultTrackingIds [][32]byte
	RetailerIds       []string
	QrCodeCIDs        []string
}, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "getTrackingByLotId", _productLotId)

	outstruct := new(struct {
		ResultTrackingIds [][32]byte
		RetailerIds       []string
		QrCodeCIDs        []string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ResultTrackingIds = *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	outstruct.RetailerIds = *abi.ConvertType(out[1], new([]string)).(*[]string)
	outstruct.QrCodeCIDs = *abi.ConvertType(out[2], new([]string)).(*[]string)

	return *outstruct, err

}

// GetTrackingByLotId is a free data retrieval call binding the contract method 0x65f0f7d5.
//
// Solidity: function getTrackingByLotId(bytes32 _productLotId) view returns(bytes32[] resultTrackingIds, string[] retailerIds, string[] qrCodeCIDs)
func (_Tracking *TrackingSession) GetTrackingByLotId(_productLotId [32]byte) (struct {
	ResultTrackingIds [][32]byte
	RetailerIds       []string
	QrCodeCIDs        []string
}, error) {
	return _Tracking.Contract.GetTrackingByLotId(&_Tracking.CallOpts, _productLotId)
}

// GetTrackingByLotId is a free data retrieval call binding the contract method 0x65f0f7d5.
//
// Solidity: function getTrackingByLotId(bytes32 _productLotId) view returns(bytes32[] resultTrackingIds, string[] retailerIds, string[] qrCodeCIDs)
func (_Tracking *TrackingCallerSession) GetTrackingByLotId(_productLotId [32]byte) (struct {
	ResultTrackingIds [][32]byte
	RetailerIds       []string
	QrCodeCIDs        []string
}, error) {
	return _Tracking.Contract.GetTrackingByLotId(&_Tracking.CallOpts, _productLotId)
}

// GetTrackingByRetailer is a free data retrieval call binding the contract method 0x83fb9889.
//
// Solidity: function getTrackingByRetailer(string _retailerId) view returns(bytes32[])
func (_Tracking *TrackingCaller) GetTrackingByRetailer(opts *bind.CallOpts, _retailerId string) ([][32]byte, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "getTrackingByRetailer", _retailerId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetTrackingByRetailer is a free data retrieval call binding the contract method 0x83fb9889.
//
// Solidity: function getTrackingByRetailer(string _retailerId) view returns(bytes32[])
func (_Tracking *TrackingSession) GetTrackingByRetailer(_retailerId string) ([][32]byte, error) {
	return _Tracking.Contract.GetTrackingByRetailer(&_Tracking.CallOpts, _retailerId)
}

// GetTrackingByRetailer is a free data retrieval call binding the contract method 0x83fb9889.
//
// Solidity: function getTrackingByRetailer(string _retailerId) view returns(bytes32[])
func (_Tracking *TrackingCallerSession) GetTrackingByRetailer(_retailerId string) ([][32]byte, error) {
	return _Tracking.Contract.GetTrackingByRetailer(&_Tracking.CallOpts, _retailerId)
}

// LogisticsCheckpoints is a free data retrieval call binding the contract method 0xf30459e3.
//
// Solidity: function logisticsCheckpoints(bytes32 , uint256 ) view returns(bytes32 trackingId, address logisticsProvider, uint256 pickupTime, uint256 deliveryTime, uint256 quantity, int256 temperature, string personInCharge, uint8 checkType, string receiverCID)
func (_Tracking *TrackingCaller) LogisticsCheckpoints(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (struct {
	TrackingId        [32]byte
	LogisticsProvider common.Address
	PickupTime        *big.Int
	DeliveryTime      *big.Int
	Quantity          *big.Int
	Temperature       *big.Int
	PersonInCharge    string
	CheckType         uint8
	ReceiverCID       string
}, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "logisticsCheckpoints", arg0, arg1)

	outstruct := new(struct {
		TrackingId        [32]byte
		LogisticsProvider common.Address
		PickupTime        *big.Int
		DeliveryTime      *big.Int
		Quantity          *big.Int
		Temperature       *big.Int
		PersonInCharge    string
		CheckType         uint8
		ReceiverCID       string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TrackingId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.LogisticsProvider = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.PickupTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.DeliveryTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Quantity = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Temperature = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.PersonInCharge = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.CheckType = *abi.ConvertType(out[7], new(uint8)).(*uint8)
	outstruct.ReceiverCID = *abi.ConvertType(out[8], new(string)).(*string)

	return *outstruct, err

}

// LogisticsCheckpoints is a free data retrieval call binding the contract method 0xf30459e3.
//
// Solidity: function logisticsCheckpoints(bytes32 , uint256 ) view returns(bytes32 trackingId, address logisticsProvider, uint256 pickupTime, uint256 deliveryTime, uint256 quantity, int256 temperature, string personInCharge, uint8 checkType, string receiverCID)
func (_Tracking *TrackingSession) LogisticsCheckpoints(arg0 [32]byte, arg1 *big.Int) (struct {
	TrackingId        [32]byte
	LogisticsProvider common.Address
	PickupTime        *big.Int
	DeliveryTime      *big.Int
	Quantity          *big.Int
	Temperature       *big.Int
	PersonInCharge    string
	CheckType         uint8
	ReceiverCID       string
}, error) {
	return _Tracking.Contract.LogisticsCheckpoints(&_Tracking.CallOpts, arg0, arg1)
}

// LogisticsCheckpoints is a free data retrieval call binding the contract method 0xf30459e3.
//
// Solidity: function logisticsCheckpoints(bytes32 , uint256 ) view returns(bytes32 trackingId, address logisticsProvider, uint256 pickupTime, uint256 deliveryTime, uint256 quantity, int256 temperature, string personInCharge, uint8 checkType, string receiverCID)
func (_Tracking *TrackingCallerSession) LogisticsCheckpoints(arg0 [32]byte, arg1 *big.Int) (struct {
	TrackingId        [32]byte
	LogisticsProvider common.Address
	PickupTime        *big.Int
	DeliveryTime      *big.Int
	Quantity          *big.Int
	Temperature       *big.Int
	PersonInCharge    string
	CheckType         uint8
	ReceiverCID       string
}, error) {
	return _Tracking.Contract.LogisticsCheckpoints(&_Tracking.CallOpts, arg0, arg1)
}

// ProductLotContract is a free data retrieval call binding the contract method 0xcc1faf8e.
//
// Solidity: function productLotContract() view returns(address)
func (_Tracking *TrackingCaller) ProductLotContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "productLotContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProductLotContract is a free data retrieval call binding the contract method 0xcc1faf8e.
//
// Solidity: function productLotContract() view returns(address)
func (_Tracking *TrackingSession) ProductLotContract() (common.Address, error) {
	return _Tracking.Contract.ProductLotContract(&_Tracking.CallOpts)
}

// ProductLotContract is a free data retrieval call binding the contract method 0xcc1faf8e.
//
// Solidity: function productLotContract() view returns(address)
func (_Tracking *TrackingCallerSession) ProductLotContract() (common.Address, error) {
	return _Tracking.Contract.ProductLotContract(&_Tracking.CallOpts)
}

// RetailerConfirmations is a free data retrieval call binding the contract method 0xf1671737.
//
// Solidity: function retailerConfirmations(bytes32 ) view returns(bytes32 trackingId, string retailerId, uint256 receivedTime, string qualityCID, string personInCharge)
func (_Tracking *TrackingCaller) RetailerConfirmations(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TrackingId     [32]byte
	RetailerId     string
	ReceivedTime   *big.Int
	QualityCID     string
	PersonInCharge string
}, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "retailerConfirmations", arg0)

	outstruct := new(struct {
		TrackingId     [32]byte
		RetailerId     string
		ReceivedTime   *big.Int
		QualityCID     string
		PersonInCharge string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TrackingId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RetailerId = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.ReceivedTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.QualityCID = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.PersonInCharge = *abi.ConvertType(out[4], new(string)).(*string)

	return *outstruct, err

}

// RetailerConfirmations is a free data retrieval call binding the contract method 0xf1671737.
//
// Solidity: function retailerConfirmations(bytes32 ) view returns(bytes32 trackingId, string retailerId, uint256 receivedTime, string qualityCID, string personInCharge)
func (_Tracking *TrackingSession) RetailerConfirmations(arg0 [32]byte) (struct {
	TrackingId     [32]byte
	RetailerId     string
	ReceivedTime   *big.Int
	QualityCID     string
	PersonInCharge string
}, error) {
	return _Tracking.Contract.RetailerConfirmations(&_Tracking.CallOpts, arg0)
}

// RetailerConfirmations is a free data retrieval call binding the contract method 0xf1671737.
//
// Solidity: function retailerConfirmations(bytes32 ) view returns(bytes32 trackingId, string retailerId, uint256 receivedTime, string qualityCID, string personInCharge)
func (_Tracking *TrackingCallerSession) RetailerConfirmations(arg0 [32]byte) (struct {
	TrackingId     [32]byte
	RetailerId     string
	ReceivedTime   *big.Int
	QualityCID     string
	PersonInCharge string
}, error) {
	return _Tracking.Contract.RetailerConfirmations(&_Tracking.CallOpts, arg0)
}

// RetailerTracking is a free data retrieval call binding the contract method 0xf43dc2bc.
//
// Solidity: function retailerTracking(string , uint256 ) view returns(bytes32)
func (_Tracking *TrackingCaller) RetailerTracking(opts *bind.CallOpts, arg0 string, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "retailerTracking", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RetailerTracking is a free data retrieval call binding the contract method 0xf43dc2bc.
//
// Solidity: function retailerTracking(string , uint256 ) view returns(bytes32)
func (_Tracking *TrackingSession) RetailerTracking(arg0 string, arg1 *big.Int) ([32]byte, error) {
	return _Tracking.Contract.RetailerTracking(&_Tracking.CallOpts, arg0, arg1)
}

// RetailerTracking is a free data retrieval call binding the contract method 0xf43dc2bc.
//
// Solidity: function retailerTracking(string , uint256 ) view returns(bytes32)
func (_Tracking *TrackingCallerSession) RetailerTracking(arg0 string, arg1 *big.Int) ([32]byte, error) {
	return _Tracking.Contract.RetailerTracking(&_Tracking.CallOpts, arg0, arg1)
}

// TrackingEvents is a free data retrieval call binding the contract method 0x9ed6f8ab.
//
// Solidity: function trackingEvents(bytes32 ) view returns(bytes32 trackingId, bytes32 productLotId, string retailerId, string qrCodeCID, uint8 status)
func (_Tracking *TrackingCaller) TrackingEvents(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TrackingId   [32]byte
	ProductLotId [32]byte
	RetailerId   string
	QrCodeCID    string
	Status       uint8
}, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "trackingEvents", arg0)

	outstruct := new(struct {
		TrackingId   [32]byte
		ProductLotId [32]byte
		RetailerId   string
		QrCodeCID    string
		Status       uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TrackingId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ProductLotId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.RetailerId = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.QrCodeCID = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// TrackingEvents is a free data retrieval call binding the contract method 0x9ed6f8ab.
//
// Solidity: function trackingEvents(bytes32 ) view returns(bytes32 trackingId, bytes32 productLotId, string retailerId, string qrCodeCID, uint8 status)
func (_Tracking *TrackingSession) TrackingEvents(arg0 [32]byte) (struct {
	TrackingId   [32]byte
	ProductLotId [32]byte
	RetailerId   string
	QrCodeCID    string
	Status       uint8
}, error) {
	return _Tracking.Contract.TrackingEvents(&_Tracking.CallOpts, arg0)
}

// TrackingEvents is a free data retrieval call binding the contract method 0x9ed6f8ab.
//
// Solidity: function trackingEvents(bytes32 ) view returns(bytes32 trackingId, bytes32 productLotId, string retailerId, string qrCodeCID, uint8 status)
func (_Tracking *TrackingCallerSession) TrackingEvents(arg0 [32]byte) (struct {
	TrackingId   [32]byte
	ProductLotId [32]byte
	RetailerId   string
	QrCodeCID    string
	Status       uint8
}, error) {
	return _Tracking.Contract.TrackingEvents(&_Tracking.CallOpts, arg0)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Tracking *TrackingCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Tracking.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Tracking *TrackingSession) UserRegistry() (common.Address, error) {
	return _Tracking.Contract.UserRegistry(&_Tracking.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Tracking *TrackingCallerSession) UserRegistry() (common.Address, error) {
	return _Tracking.Contract.UserRegistry(&_Tracking.CallOpts)
}

// CreateTrackingEvent is a paid mutator transaction binding the contract method 0x465b6e73.
//
// Solidity: function createTrackingEvent(bytes32 _trackingId, bytes32 _productLotId, string _retailerId, string _qrCodeCID) returns()
func (_Tracking *TrackingTransactor) CreateTrackingEvent(opts *bind.TransactOpts, _trackingId [32]byte, _productLotId [32]byte, _retailerId string, _qrCodeCID string) (*types.Transaction, error) {
	return _Tracking.contract.Transact(opts, "createTrackingEvent", _trackingId, _productLotId, _retailerId, _qrCodeCID)
}

// CreateTrackingEvent is a paid mutator transaction binding the contract method 0x465b6e73.
//
// Solidity: function createTrackingEvent(bytes32 _trackingId, bytes32 _productLotId, string _retailerId, string _qrCodeCID) returns()
func (_Tracking *TrackingSession) CreateTrackingEvent(_trackingId [32]byte, _productLotId [32]byte, _retailerId string, _qrCodeCID string) (*types.Transaction, error) {
	return _Tracking.Contract.CreateTrackingEvent(&_Tracking.TransactOpts, _trackingId, _productLotId, _retailerId, _qrCodeCID)
}

// CreateTrackingEvent is a paid mutator transaction binding the contract method 0x465b6e73.
//
// Solidity: function createTrackingEvent(bytes32 _trackingId, bytes32 _productLotId, string _retailerId, string _qrCodeCID) returns()
func (_Tracking *TrackingTransactorSession) CreateTrackingEvent(_trackingId [32]byte, _productLotId [32]byte, _retailerId string, _qrCodeCID string) (*types.Transaction, error) {
	return _Tracking.Contract.CreateTrackingEvent(&_Tracking.TransactOpts, _trackingId, _productLotId, _retailerId, _qrCodeCID)
}

// RetailerReceiveProduct is a paid mutator transaction binding the contract method 0x7ea96547.
//
// Solidity: function retailerReceiveProduct(bytes32 _trackingId, string _retailerId, string _qualityCID, string _personInCharge) returns()
func (_Tracking *TrackingTransactor) RetailerReceiveProduct(opts *bind.TransactOpts, _trackingId [32]byte, _retailerId string, _qualityCID string, _personInCharge string) (*types.Transaction, error) {
	return _Tracking.contract.Transact(opts, "retailerReceiveProduct", _trackingId, _retailerId, _qualityCID, _personInCharge)
}

// RetailerReceiveProduct is a paid mutator transaction binding the contract method 0x7ea96547.
//
// Solidity: function retailerReceiveProduct(bytes32 _trackingId, string _retailerId, string _qualityCID, string _personInCharge) returns()
func (_Tracking *TrackingSession) RetailerReceiveProduct(_trackingId [32]byte, _retailerId string, _qualityCID string, _personInCharge string) (*types.Transaction, error) {
	return _Tracking.Contract.RetailerReceiveProduct(&_Tracking.TransactOpts, _trackingId, _retailerId, _qualityCID, _personInCharge)
}

// RetailerReceiveProduct is a paid mutator transaction binding the contract method 0x7ea96547.
//
// Solidity: function retailerReceiveProduct(bytes32 _trackingId, string _retailerId, string _qualityCID, string _personInCharge) returns()
func (_Tracking *TrackingTransactorSession) RetailerReceiveProduct(_trackingId [32]byte, _retailerId string, _qualityCID string, _personInCharge string) (*types.Transaction, error) {
	return _Tracking.Contract.RetailerReceiveProduct(&_Tracking.TransactOpts, _trackingId, _retailerId, _qualityCID, _personInCharge)
}

// UpdateLogisticsCheckpoint is a paid mutator transaction binding the contract method 0xa7e604b2.
//
// Solidity: function updateLogisticsCheckpoint(bytes32 _trackingId, uint256 _pickupTime, uint256 _deliveryTime, uint256 _quantity, int256 _temperature, string _personInCharge, uint8 _checkType, string _receiverCID) returns()
func (_Tracking *TrackingTransactor) UpdateLogisticsCheckpoint(opts *bind.TransactOpts, _trackingId [32]byte, _pickupTime *big.Int, _deliveryTime *big.Int, _quantity *big.Int, _temperature *big.Int, _personInCharge string, _checkType uint8, _receiverCID string) (*types.Transaction, error) {
	return _Tracking.contract.Transact(opts, "updateLogisticsCheckpoint", _trackingId, _pickupTime, _deliveryTime, _quantity, _temperature, _personInCharge, _checkType, _receiverCID)
}

// UpdateLogisticsCheckpoint is a paid mutator transaction binding the contract method 0xa7e604b2.
//
// Solidity: function updateLogisticsCheckpoint(bytes32 _trackingId, uint256 _pickupTime, uint256 _deliveryTime, uint256 _quantity, int256 _temperature, string _personInCharge, uint8 _checkType, string _receiverCID) returns()
func (_Tracking *TrackingSession) UpdateLogisticsCheckpoint(_trackingId [32]byte, _pickupTime *big.Int, _deliveryTime *big.Int, _quantity *big.Int, _temperature *big.Int, _personInCharge string, _checkType uint8, _receiverCID string) (*types.Transaction, error) {
	return _Tracking.Contract.UpdateLogisticsCheckpoint(&_Tracking.TransactOpts, _trackingId, _pickupTime, _deliveryTime, _quantity, _temperature, _personInCharge, _checkType, _receiverCID)
}

// UpdateLogisticsCheckpoint is a paid mutator transaction binding the contract method 0xa7e604b2.
//
// Solidity: function updateLogisticsCheckpoint(bytes32 _trackingId, uint256 _pickupTime, uint256 _deliveryTime, uint256 _quantity, int256 _temperature, string _personInCharge, uint8 _checkType, string _receiverCID) returns()
func (_Tracking *TrackingTransactorSession) UpdateLogisticsCheckpoint(_trackingId [32]byte, _pickupTime *big.Int, _deliveryTime *big.Int, _quantity *big.Int, _temperature *big.Int, _personInCharge string, _checkType uint8, _receiverCID string) (*types.Transaction, error) {
	return _Tracking.Contract.UpdateLogisticsCheckpoint(&_Tracking.TransactOpts, _trackingId, _pickupTime, _deliveryTime, _quantity, _temperature, _personInCharge, _checkType, _receiverCID)
}

// TrackingLogisticsUpdatedIterator is returned from FilterLogisticsUpdated and is used to iterate over the raw logs and unpacked data for LogisticsUpdated events raised by the Tracking contract.
type TrackingLogisticsUpdatedIterator struct {
	Event *TrackingLogisticsUpdated // Event containing the contract specifics and raw log

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
func (it *TrackingLogisticsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrackingLogisticsUpdated)
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
		it.Event = new(TrackingLogisticsUpdated)
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
func (it *TrackingLogisticsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrackingLogisticsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrackingLogisticsUpdated represents a LogisticsUpdated event raised by the Tracking contract.
type TrackingLogisticsUpdated struct {
	TrackingId        [32]byte
	LogisticsProvider common.Address
	CheckType         uint8
	ReceiverCID       string
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogisticsUpdated is a free log retrieval operation binding the contract event 0x710a10364a49834e8c6e0fafff5354ca1e160f26a1ba1d7b33c7100c80d7ea6d.
//
// Solidity: event LogisticsUpdated(bytes32 indexed trackingId, address indexed logisticsProvider, uint8 checkType, string receiverCID)
func (_Tracking *TrackingFilterer) FilterLogisticsUpdated(opts *bind.FilterOpts, trackingId [][32]byte, logisticsProvider []common.Address) (*TrackingLogisticsUpdatedIterator, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}
	var logisticsProviderRule []interface{}
	for _, logisticsProviderItem := range logisticsProvider {
		logisticsProviderRule = append(logisticsProviderRule, logisticsProviderItem)
	}

	logs, sub, err := _Tracking.contract.FilterLogs(opts, "LogisticsUpdated", trackingIdRule, logisticsProviderRule)
	if err != nil {
		return nil, err
	}
	return &TrackingLogisticsUpdatedIterator{contract: _Tracking.contract, event: "LogisticsUpdated", logs: logs, sub: sub}, nil
}

// WatchLogisticsUpdated is a free log subscription operation binding the contract event 0x710a10364a49834e8c6e0fafff5354ca1e160f26a1ba1d7b33c7100c80d7ea6d.
//
// Solidity: event LogisticsUpdated(bytes32 indexed trackingId, address indexed logisticsProvider, uint8 checkType, string receiverCID)
func (_Tracking *TrackingFilterer) WatchLogisticsUpdated(opts *bind.WatchOpts, sink chan<- *TrackingLogisticsUpdated, trackingId [][32]byte, logisticsProvider []common.Address) (event.Subscription, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}
	var logisticsProviderRule []interface{}
	for _, logisticsProviderItem := range logisticsProvider {
		logisticsProviderRule = append(logisticsProviderRule, logisticsProviderItem)
	}

	logs, sub, err := _Tracking.contract.WatchLogs(opts, "LogisticsUpdated", trackingIdRule, logisticsProviderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrackingLogisticsUpdated)
				if err := _Tracking.contract.UnpackLog(event, "LogisticsUpdated", log); err != nil {
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

// ParseLogisticsUpdated is a log parse operation binding the contract event 0x710a10364a49834e8c6e0fafff5354ca1e160f26a1ba1d7b33c7100c80d7ea6d.
//
// Solidity: event LogisticsUpdated(bytes32 indexed trackingId, address indexed logisticsProvider, uint8 checkType, string receiverCID)
func (_Tracking *TrackingFilterer) ParseLogisticsUpdated(log types.Log) (*TrackingLogisticsUpdated, error) {
	event := new(TrackingLogisticsUpdated)
	if err := _Tracking.contract.UnpackLog(event, "LogisticsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrackingRetailerReceivedIterator is returned from FilterRetailerReceived and is used to iterate over the raw logs and unpacked data for RetailerReceived events raised by the Tracking contract.
type TrackingRetailerReceivedIterator struct {
	Event *TrackingRetailerReceived // Event containing the contract specifics and raw log

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
func (it *TrackingRetailerReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrackingRetailerReceived)
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
		it.Event = new(TrackingRetailerReceived)
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
func (it *TrackingRetailerReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrackingRetailerReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrackingRetailerReceived represents a RetailerReceived event raised by the Tracking contract.
type TrackingRetailerReceived struct {
	TrackingId [32]byte
	RetailerId string
	QualityCID string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRetailerReceived is a free log retrieval operation binding the contract event 0x6931b58dbc3c7b4d9caec3ed0638b6990770e147b532536809aa388e4fa97b95.
//
// Solidity: event RetailerReceived(bytes32 indexed trackingId, string retailerId, string qualityCID)
func (_Tracking *TrackingFilterer) FilterRetailerReceived(opts *bind.FilterOpts, trackingId [][32]byte) (*TrackingRetailerReceivedIterator, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}

	logs, sub, err := _Tracking.contract.FilterLogs(opts, "RetailerReceived", trackingIdRule)
	if err != nil {
		return nil, err
	}
	return &TrackingRetailerReceivedIterator{contract: _Tracking.contract, event: "RetailerReceived", logs: logs, sub: sub}, nil
}

// WatchRetailerReceived is a free log subscription operation binding the contract event 0x6931b58dbc3c7b4d9caec3ed0638b6990770e147b532536809aa388e4fa97b95.
//
// Solidity: event RetailerReceived(bytes32 indexed trackingId, string retailerId, string qualityCID)
func (_Tracking *TrackingFilterer) WatchRetailerReceived(opts *bind.WatchOpts, sink chan<- *TrackingRetailerReceived, trackingId [][32]byte) (event.Subscription, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}

	logs, sub, err := _Tracking.contract.WatchLogs(opts, "RetailerReceived", trackingIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrackingRetailerReceived)
				if err := _Tracking.contract.UnpackLog(event, "RetailerReceived", log); err != nil {
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

// ParseRetailerReceived is a log parse operation binding the contract event 0x6931b58dbc3c7b4d9caec3ed0638b6990770e147b532536809aa388e4fa97b95.
//
// Solidity: event RetailerReceived(bytes32 indexed trackingId, string retailerId, string qualityCID)
func (_Tracking *TrackingFilterer) ParseRetailerReceived(log types.Log) (*TrackingRetailerReceived, error) {
	event := new(TrackingRetailerReceived)
	if err := _Tracking.contract.UnpackLog(event, "RetailerReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrackingTrackingCreatedIterator is returned from FilterTrackingCreated and is used to iterate over the raw logs and unpacked data for TrackingCreated events raised by the Tracking contract.
type TrackingTrackingCreatedIterator struct {
	Event *TrackingTrackingCreated // Event containing the contract specifics and raw log

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
func (it *TrackingTrackingCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrackingTrackingCreated)
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
		it.Event = new(TrackingTrackingCreated)
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
func (it *TrackingTrackingCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrackingTrackingCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrackingTrackingCreated represents a TrackingCreated event raised by the Tracking contract.
type TrackingTrackingCreated struct {
	TrackingId   [32]byte
	ProductLotId [32]byte
	RetailerId   string
	QrCodeCID    string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTrackingCreated is a free log retrieval operation binding the contract event 0xa2fc3ade0f3e8f9c31cda358ccb72883469278eb67b094b4d6d90a5cab7d5a53.
//
// Solidity: event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, string retailerId, string qrCodeCID)
func (_Tracking *TrackingFilterer) FilterTrackingCreated(opts *bind.FilterOpts, trackingId [][32]byte, productLotId [][32]byte) (*TrackingTrackingCreatedIterator, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}
	var productLotIdRule []interface{}
	for _, productLotIdItem := range productLotId {
		productLotIdRule = append(productLotIdRule, productLotIdItem)
	}

	logs, sub, err := _Tracking.contract.FilterLogs(opts, "TrackingCreated", trackingIdRule, productLotIdRule)
	if err != nil {
		return nil, err
	}
	return &TrackingTrackingCreatedIterator{contract: _Tracking.contract, event: "TrackingCreated", logs: logs, sub: sub}, nil
}

// WatchTrackingCreated is a free log subscription operation binding the contract event 0xa2fc3ade0f3e8f9c31cda358ccb72883469278eb67b094b4d6d90a5cab7d5a53.
//
// Solidity: event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, string retailerId, string qrCodeCID)
func (_Tracking *TrackingFilterer) WatchTrackingCreated(opts *bind.WatchOpts, sink chan<- *TrackingTrackingCreated, trackingId [][32]byte, productLotId [][32]byte) (event.Subscription, error) {

	var trackingIdRule []interface{}
	for _, trackingIdItem := range trackingId {
		trackingIdRule = append(trackingIdRule, trackingIdItem)
	}
	var productLotIdRule []interface{}
	for _, productLotIdItem := range productLotId {
		productLotIdRule = append(productLotIdRule, productLotIdItem)
	}

	logs, sub, err := _Tracking.contract.WatchLogs(opts, "TrackingCreated", trackingIdRule, productLotIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrackingTrackingCreated)
				if err := _Tracking.contract.UnpackLog(event, "TrackingCreated", log); err != nil {
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

// ParseTrackingCreated is a log parse operation binding the contract event 0xa2fc3ade0f3e8f9c31cda358ccb72883469278eb67b094b4d6d90a5cab7d5a53.
//
// Solidity: event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, string retailerId, string qrCodeCID)
func (_Tracking *TrackingFilterer) ParseTrackingCreated(log types.Log) (*TrackingTrackingCreated, error) {
	event := new(TrackingTrackingCreated)
	if err := _Tracking.contract.UnpackLog(event, "TrackingCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
