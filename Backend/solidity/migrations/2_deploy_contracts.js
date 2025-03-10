const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");
const Product = artifacts.require("Product");

module.exports = async function (deployer, network, accounts) {
  console.log("🚀 Starting contract deployment...");

  await deployer.deploy(UserRegistry);
  const userRegistryInstance = await UserRegistry.deployed();
  console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);

  await deployer.deploy(Product, userRegistryInstance.address);
  const productInstance = await Product.deployed();
  console.log("✅ Product Contract Deployed at:", productInstance.address);

  await deployer.deploy(RawMilk, userRegistryInstance.address);
  const rawMilkInstance = await RawMilk.deployed();
  console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);

  // ✅ เพิ่ม parameter userRegistryInstance.address
  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
  const certificationEventInstance = await CertificationEvent.deployed();
  console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);

  console.log("🎉 Deployment completed successfully!");
};
