module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 7545,
      network_id: "*",
      gas: 8000000, // ✅ เพิ่มจาก 6M → 8M
      gasPrice: 20000000000 // ✅ คงค่าเดิมที่ 20 Gwei
    }
  },
  compilers: {
    solc: {
      version: "0.8.19"
    }
  }
};
