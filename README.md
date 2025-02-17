# Blockchain-Based Agricultural Supply Chain Backend

## Project Overview

This project is a backend system for a **blockchain-based agricultural supply chain** that ensures transparency, security, and traceability of organic food products. The system integrates **PostgreSQL**, **IPFS**, and **Blockchain** (Gnosis Chain for testing) to store, manage, and verify supply chain data.

## Tech Stack

- **Backend:** Golang (Fiber Framework)
- **Database:** PostgreSQL
- **Blockchain:** Solidity (Smart Contracts deployed on Gnosis Chain)
- **File Storage:** IPFS (for storing certification & quality inspection documents)
- **Authentication:** JWT-based authentication
- **API Framework:** Fiber

## Project Structure

```
Backend/
│── api/
│   ├── controllers/       # API controllers
│   ├── routes/            # API route definitions
│── config/                # Configuration files
│── database/              # Database connection setup
│── middleware/            # Authentication & security
│── models/                # Database models
│── solidity/              # Smart Contracts
│── services/              # Blockchain & IPFS services
│── utils/                 # Utility functions
│── main.go                # Application entry point
│── .env                   # Environment variables
```

## Database Structure

### **PostgreSQL (For storing core entities)**

| Table Name             | Description                                     |
| ---------------------- | ----------------------------------------------- |
| `users`                | Stores user authentication details              |
| `farmer`               | Stores farmer profile information               |
| `dairyfactory`         | Stores dairy factory details                    |
| `logisticsprovider`    | Stores logistics company information            |
| `retailer`             | Stores retailer details                         |
| `organiccertification` | Stores organic certification details            |
| `externalid`           | Tracks external IDs used in supply chain events |

### **Blockchain (For immutable data storage)**

| Entity               | Description                                          |
| -------------------- | ---------------------------------------------------- |
| `RawMilk`            | Tracks raw milk production details                   |
| `Product`            | Stores processed product information                 |
| `ProductLot`         | Tracks product batches                               |
| `Shipping Events`    | Stores supply chain events like transport & delivery |
| `CertificationEvent` | Logs certification updates                           |

### **IPFS (For storing important documents)**

| File Type              | Storage |
| ---------------------- | ------- |
| Organic Certifications | IPFS    |
| Quality Inspections    | IPFS    |
| Nutritional Info       | IPFS    |

## Setup Instructions

### **1. Clone the Repository**

```sh
git clone https://github.com/your-repo.git
cd Backend
```

### **2. Set Up Environment Variables (.env file)**

Create a `.env` file and configure the necessary variables:

#### **For Ganache Testing (Local Blockchain Development)**

```sh
PORT=8080
DATABASE_URL="postgres://user:password@localhost:5432/dairy_supplychain"
BLOCKCHAIN_RPC_URL=http://127.0.0.1:7545
PRIVATE_KEY=b5a0908e76578dce71f1b8e2c81caf94e29b6769bf7e068f89e6e75076c19e4e
CERT_CONTRACT_ADDRESS=0x37a93e7FbCAA6BB4A7806C7Eabe58a067630A7Ff
JWT_SECRET=rI7Qaj1LvHxgNAxqiu3C3+T3oI3jWKzqes1La3jiIrk=
RAWMILK_CONTRACT_ADDRESS=0x61052244e424b8B3BC7457F9afB1177608395CA3
```

#### **For Gnosis Chain (Testnet Deployment)**

```sh
PORT=8080
DATABASE_URL="postgres://user:password@localhost:5432/dairy_supplychain"
BLOCKCHAIN_RPC_URL=https://rpc.gnosischain.com
PRIVATE_KEY=your_gnosis_chain_private_key
CERT_CONTRACT_ADDRESS=0xYourGnosisCertContractAddress
JWT_SECRET=your_secure_jwt_secret
RAWMILK_CONTRACT_ADDRESS=0xYourGnosisRawMilkContractAddress
```

### **3. Install Dependencies**

```sh
go mod tidy
```

### **4. Run the Server**

```sh
go run main.go
```

### **5. API Routes Overview**

| Endpoint                        | Method | Description                    |
| ------------------------------- | ------ | ------------------------------ |
| `/api/v1/auth/register`         | POST   | Register a new user            |
| `/api/v1/auth/login`            | POST   | Authenticate a user            |
| `/api/v1/auth/update-role`      | POST   | Assign user roles              |
| `/api/v1/farmers`               | POST   | Register a new farmer          |
| `/api/v1/rawmilk`               | POST   | Add raw milk entry             |
| `/api/v1/certifications/upload` | POST   | Upload a certification to IPFS |
| `/api/v1/certifications/create` | POST   | Create a certification record  |

## Smart Contract Integration

- Smart contracts are deployed on **Gnosis Chain** (for testing).
- Blockchain interactions include **tracking raw milk, products, and certification events**.

## Security Measures

- **JWT Authentication** for secure access
- **Role-Based Access Control (RBAC)**
- **CORS Protection**
- **Blockchain for Immutable Record Keeping**

## Future Improvements

- Integrate real-time tracking for logistics
- Implement **Zero-Knowledge Proofs (ZKP)** for privacy-preserving verification
- Enhance **UI for Farmers & Retailers**

## Contributors

- **Backend Developer:** [Your Name]
- **Smart Contract Developer:** [Contributor Name]
- **Blockchain & IPFS Specialist:** [Contributor Name]

## License

MIT License

