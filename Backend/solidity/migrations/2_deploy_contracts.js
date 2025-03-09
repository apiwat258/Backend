const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");
const Product = artifacts.require("Product"); // ✅ เพิ่ม Product

module.exports = async function (deployer, network, accounts) {
  console.log("🚀 Starting contract deployment...");

  // ✅ ตรวจสอบว่ามี UserRegistry แล้ว
  let userRegistryInstance;
  try {
    userRegistryInstance = await UserRegistry.deployed();
    console.log("✅ Using existing UserRegistry at:", userRegistryInstance.address);
  } catch (error) {
    console.log("🚨 UserRegistry not found, deploying a new one...");
    await deployer.deploy(UserRegistry);
    userRegistryInstance = await UserRegistry.deployed();
    console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);
  }

  // ✅ ตรวจสอบว่า Product ถูกดีพลอยหรือยัง
  let productInstance;
  try {
    productInstance = await Product.deployed();
    console.log("✅ Using existing Product contract at:", productInstance.address);
  } catch (error) {
    console.log("🚀 Deploying Product Contract...");
    await deployer.deploy(Product, userRegistryInstance.address);
    productInstance = await Product.deployed();
    console.log("✅ Product Contract Deployed at:", productInstance.address);
  }

  console.log("🎉 Deployment completed!");
};
