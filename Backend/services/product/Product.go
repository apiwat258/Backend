// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package product

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

// ProductProductInfo is an auto generated low-level Go binding around an user-defined struct.
type ProductProductInfo struct {
	ProductId     [32]byte
	FactoryWallet common.Address
	ProductName   string
	ProductCID    string
	Category      string
}

// ProductMetaData contains all meta data concerning the Product contract.
var ProductMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"factoryWallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"productName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"productCID\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"category\",\"type\":\"string\"}],\"name\":\"ProductCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"productIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"products\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"factoryWallet\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"productName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"productCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"category\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[],\"name\":\"userRegistry\",\"outputs\":[{\"internalType\":\"contractUserRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factoryWallet\",\"type\":\"address\"}],\"name\":\"getProductsByFactory\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"},{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_productId\",\"type\":\"bytes32\"}],\"name\":\"getProductById\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"productId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"factoryWallet\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"productName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"productCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"category\",\"type\":\"string\"}],\"internalType\":\"structProduct.ProductInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_productId\",\"type\":\"bytes32\"}],\"name\":\"isProductExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"constant\":true},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_productId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_productName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_productCID\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_category\",\"type\":\"string\"}],\"name\":\"createProduct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ProductABI is the input ABI used to generate the binding from.
// Deprecated: Use ProductMetaData.ABI instead.
var ProductABI = ProductMetaData.ABI

// Product is an auto generated Go binding around an Ethereum contract.
type Product struct {
	ProductCaller     // Read-only binding to the contract
	ProductTransactor // Write-only binding to the contract
	ProductFilterer   // Log filterer for contract events
}

// ProductCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProductCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProductTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProductFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProductSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProductSession struct {
	Contract     *Product          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProductCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProductCallerSession struct {
	Contract *ProductCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ProductTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProductTransactorSession struct {
	Contract     *ProductTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ProductRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProductRaw struct {
	Contract *Product // Generic contract binding to access the raw methods on
}

// ProductCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProductCallerRaw struct {
	Contract *ProductCaller // Generic read-only contract binding to access the raw methods on
}

// ProductTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProductTransactorRaw struct {
	Contract *ProductTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProduct creates a new instance of Product, bound to a specific deployed contract.
func NewProduct(address common.Address, backend bind.ContractBackend) (*Product, error) {
	contract, err := bindProduct(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Product{ProductCaller: ProductCaller{contract: contract}, ProductTransactor: ProductTransactor{contract: contract}, ProductFilterer: ProductFilterer{contract: contract}}, nil
}

// NewProductCaller creates a new read-only instance of Product, bound to a specific deployed contract.
func NewProductCaller(address common.Address, caller bind.ContractCaller) (*ProductCaller, error) {
	contract, err := bindProduct(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProductCaller{contract: contract}, nil
}

// NewProductTransactor creates a new write-only instance of Product, bound to a specific deployed contract.
func NewProductTransactor(address common.Address, transactor bind.ContractTransactor) (*ProductTransactor, error) {
	contract, err := bindProduct(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProductTransactor{contract: contract}, nil
}

// NewProductFilterer creates a new log filterer instance of Product, bound to a specific deployed contract.
func NewProductFilterer(address common.Address, filterer bind.ContractFilterer) (*ProductFilterer, error) {
	contract, err := bindProduct(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProductFilterer{contract: contract}, nil
}

// bindProduct binds a generic wrapper to an already deployed contract.
func bindProduct(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProductMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Product *ProductRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Product.Contract.ProductCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Product *ProductRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Product.Contract.ProductTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Product *ProductRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Product.Contract.ProductTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Product *ProductCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Product.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Product *ProductTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Product.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Product *ProductTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Product.Contract.contract.Transact(opts, method, params...)
}

// GetProductById is a free data retrieval call binding the contract method 0x5ea57703.
//
// Solidity: function getProductById(bytes32 _productId) view returns((bytes32,address,string,string,string))
func (_Product *ProductCaller) GetProductById(opts *bind.CallOpts, _productId [32]byte) (ProductProductInfo, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "getProductById", _productId)

	if err != nil {
		return *new(ProductProductInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ProductProductInfo)).(*ProductProductInfo)

	return out0, err

}

// GetProductById is a free data retrieval call binding the contract method 0x5ea57703.
//
// Solidity: function getProductById(bytes32 _productId) view returns((bytes32,address,string,string,string))
func (_Product *ProductSession) GetProductById(_productId [32]byte) (ProductProductInfo, error) {
	return _Product.Contract.GetProductById(&_Product.CallOpts, _productId)
}

// GetProductById is a free data retrieval call binding the contract method 0x5ea57703.
//
// Solidity: function getProductById(bytes32 _productId) view returns((bytes32,address,string,string,string))
func (_Product *ProductCallerSession) GetProductById(_productId [32]byte) (ProductProductInfo, error) {
	return _Product.Contract.GetProductById(&_Product.CallOpts, _productId)
}

// GetProductsByFactory is a free data retrieval call binding the contract method 0x99513ad9.
//
// Solidity: function getProductsByFactory(address _factoryWallet) view returns(bytes32[], string[], string[])
func (_Product *ProductCaller) GetProductsByFactory(opts *bind.CallOpts, _factoryWallet common.Address) ([][32]byte, []string, []string, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "getProductsByFactory", _factoryWallet)

	if err != nil {
		return *new([][32]byte), *new([]string), *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	out1 := *abi.ConvertType(out[1], new([]string)).(*[]string)
	out2 := *abi.ConvertType(out[2], new([]string)).(*[]string)

	return out0, out1, out2, err

}

// GetProductsByFactory is a free data retrieval call binding the contract method 0x99513ad9.
//
// Solidity: function getProductsByFactory(address _factoryWallet) view returns(bytes32[], string[], string[])
func (_Product *ProductSession) GetProductsByFactory(_factoryWallet common.Address) ([][32]byte, []string, []string, error) {
	return _Product.Contract.GetProductsByFactory(&_Product.CallOpts, _factoryWallet)
}

// GetProductsByFactory is a free data retrieval call binding the contract method 0x99513ad9.
//
// Solidity: function getProductsByFactory(address _factoryWallet) view returns(bytes32[], string[], string[])
func (_Product *ProductCallerSession) GetProductsByFactory(_factoryWallet common.Address) ([][32]byte, []string, []string, error) {
	return _Product.Contract.GetProductsByFactory(&_Product.CallOpts, _factoryWallet)
}

// IsProductExist is a free data retrieval call binding the contract method 0xdd91b929.
//
// Solidity: function isProductExist(bytes32 _productId) view returns(bool)
func (_Product *ProductCaller) IsProductExist(opts *bind.CallOpts, _productId [32]byte) (bool, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "isProductExist", _productId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsProductExist is a free data retrieval call binding the contract method 0xdd91b929.
//
// Solidity: function isProductExist(bytes32 _productId) view returns(bool)
func (_Product *ProductSession) IsProductExist(_productId [32]byte) (bool, error) {
	return _Product.Contract.IsProductExist(&_Product.CallOpts, _productId)
}

// IsProductExist is a free data retrieval call binding the contract method 0xdd91b929.
//
// Solidity: function isProductExist(bytes32 _productId) view returns(bool)
func (_Product *ProductCallerSession) IsProductExist(_productId [32]byte) (bool, error) {
	return _Product.Contract.IsProductExist(&_Product.CallOpts, _productId)
}

// ProductIds is a free data retrieval call binding the contract method 0xcb51bba9.
//
// Solidity: function productIds(uint256 ) view returns(bytes32)
func (_Product *ProductCaller) ProductIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "productIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProductIds is a free data retrieval call binding the contract method 0xcb51bba9.
//
// Solidity: function productIds(uint256 ) view returns(bytes32)
func (_Product *ProductSession) ProductIds(arg0 *big.Int) ([32]byte, error) {
	return _Product.Contract.ProductIds(&_Product.CallOpts, arg0)
}

// ProductIds is a free data retrieval call binding the contract method 0xcb51bba9.
//
// Solidity: function productIds(uint256 ) view returns(bytes32)
func (_Product *ProductCallerSession) ProductIds(arg0 *big.Int) ([32]byte, error) {
	return _Product.Contract.ProductIds(&_Product.CallOpts, arg0)
}

// Products is a free data retrieval call binding the contract method 0x79054391.
//
// Solidity: function products(bytes32 ) view returns(bytes32 productId, address factoryWallet, string productName, string productCID, string category)
func (_Product *ProductCaller) Products(opts *bind.CallOpts, arg0 [32]byte) (struct {
	ProductId     [32]byte
	FactoryWallet common.Address
	ProductName   string
	ProductCID    string
	Category      string
}, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "products", arg0)

	outstruct := new(struct {
		ProductId     [32]byte
		FactoryWallet common.Address
		ProductName   string
		ProductCID    string
		Category      string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ProductId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.FactoryWallet = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ProductName = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ProductCID = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.Category = *abi.ConvertType(out[4], new(string)).(*string)

	return *outstruct, err

}

// Products is a free data retrieval call binding the contract method 0x79054391.
//
// Solidity: function products(bytes32 ) view returns(bytes32 productId, address factoryWallet, string productName, string productCID, string category)
func (_Product *ProductSession) Products(arg0 [32]byte) (struct {
	ProductId     [32]byte
	FactoryWallet common.Address
	ProductName   string
	ProductCID    string
	Category      string
}, error) {
	return _Product.Contract.Products(&_Product.CallOpts, arg0)
}

// Products is a free data retrieval call binding the contract method 0x79054391.
//
// Solidity: function products(bytes32 ) view returns(bytes32 productId, address factoryWallet, string productName, string productCID, string category)
func (_Product *ProductCallerSession) Products(arg0 [32]byte) (struct {
	ProductId     [32]byte
	FactoryWallet common.Address
	ProductName   string
	ProductCID    string
	Category      string
}, error) {
	return _Product.Contract.Products(&_Product.CallOpts, arg0)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Product *ProductCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Product.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Product *ProductSession) UserRegistry() (common.Address, error) {
	return _Product.Contract.UserRegistry(&_Product.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Product *ProductCallerSession) UserRegistry() (common.Address, error) {
	return _Product.Contract.UserRegistry(&_Product.CallOpts)
}

// CreateProduct is a paid mutator transaction binding the contract method 0x290344c7.
//
// Solidity: function createProduct(bytes32 _productId, string _productName, string _productCID, string _category) returns()
func (_Product *ProductTransactor) CreateProduct(opts *bind.TransactOpts, _productId [32]byte, _productName string, _productCID string, _category string) (*types.Transaction, error) {
	return _Product.contract.Transact(opts, "createProduct", _productId, _productName, _productCID, _category)
}

// CreateProduct is a paid mutator transaction binding the contract method 0x290344c7.
//
// Solidity: function createProduct(bytes32 _productId, string _productName, string _productCID, string _category) returns()
func (_Product *ProductSession) CreateProduct(_productId [32]byte, _productName string, _productCID string, _category string) (*types.Transaction, error) {
	return _Product.Contract.CreateProduct(&_Product.TransactOpts, _productId, _productName, _productCID, _category)
}

// CreateProduct is a paid mutator transaction binding the contract method 0x290344c7.
//
// Solidity: function createProduct(bytes32 _productId, string _productName, string _productCID, string _category) returns()
func (_Product *ProductTransactorSession) CreateProduct(_productId [32]byte, _productName string, _productCID string, _category string) (*types.Transaction, error) {
	return _Product.Contract.CreateProduct(&_Product.TransactOpts, _productId, _productName, _productCID, _category)
}

// ProductProductCreatedIterator is returned from FilterProductCreated and is used to iterate over the raw logs and unpacked data for ProductCreated events raised by the Product contract.
type ProductProductCreatedIterator struct {
	Event *ProductProductCreated // Event containing the contract specifics and raw log

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
func (it *ProductProductCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProductProductCreated)
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
		it.Event = new(ProductProductCreated)
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
func (it *ProductProductCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProductProductCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProductProductCreated represents a ProductCreated event raised by the Product contract.
type ProductProductCreated struct {
	ProductId     [32]byte
	FactoryWallet common.Address
	ProductName   string
	ProductCID    string
	Category      string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterProductCreated is a free log retrieval operation binding the contract event 0xce36d41273b724ef277b938c439afb2342a6d7a476a905409f0b300fbbadcc72.
//
// Solidity: event ProductCreated(bytes32 indexed productId, address indexed factoryWallet, string productName, string productCID, string category)
func (_Product *ProductFilterer) FilterProductCreated(opts *bind.FilterOpts, productId [][32]byte, factoryWallet []common.Address) (*ProductProductCreatedIterator, error) {

	var productIdRule []interface{}
	for _, productIdItem := range productId {
		productIdRule = append(productIdRule, productIdItem)
	}
	var factoryWalletRule []interface{}
	for _, factoryWalletItem := range factoryWallet {
		factoryWalletRule = append(factoryWalletRule, factoryWalletItem)
	}

	logs, sub, err := _Product.contract.FilterLogs(opts, "ProductCreated", productIdRule, factoryWalletRule)
	if err != nil {
		return nil, err
	}
	return &ProductProductCreatedIterator{contract: _Product.contract, event: "ProductCreated", logs: logs, sub: sub}, nil
}

// WatchProductCreated is a free log subscription operation binding the contract event 0xce36d41273b724ef277b938c439afb2342a6d7a476a905409f0b300fbbadcc72.
//
// Solidity: event ProductCreated(bytes32 indexed productId, address indexed factoryWallet, string productName, string productCID, string category)
func (_Product *ProductFilterer) WatchProductCreated(opts *bind.WatchOpts, sink chan<- *ProductProductCreated, productId [][32]byte, factoryWallet []common.Address) (event.Subscription, error) {

	var productIdRule []interface{}
	for _, productIdItem := range productId {
		productIdRule = append(productIdRule, productIdItem)
	}
	var factoryWalletRule []interface{}
	for _, factoryWalletItem := range factoryWallet {
		factoryWalletRule = append(factoryWalletRule, factoryWalletItem)
	}

	logs, sub, err := _Product.contract.WatchLogs(opts, "ProductCreated", productIdRule, factoryWalletRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProductProductCreated)
				if err := _Product.contract.UnpackLog(event, "ProductCreated", log); err != nil {
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

// ParseProductCreated is a log parse operation binding the contract event 0xce36d41273b724ef277b938c439afb2342a6d7a476a905409f0b300fbbadcc72.
//
// Solidity: event ProductCreated(bytes32 indexed productId, address indexed factoryWallet, string productName, string productCID, string category)
func (_Product *ProductFilterer) ParseProductCreated(log types.Log) (*ProductProductCreated, error) {
	event := new(ProductProductCreated)
	if err := _Product.contract.UnpackLog(event, "ProductCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
