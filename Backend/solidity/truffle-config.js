module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 7545,
      network_id: "*",
      gas: 6000000, // ✅ เพิ่ม gas limit
      gasPrice: 20000000000 // 20 Gwei
    }
  },
  compilers: {
    solc: {
      version: "0.8.19"
    }
  }
};
