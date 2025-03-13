const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");
const Product = artifacts.require("Product");
const ProductLot = artifacts.require("ProductLot");
const Tracking = artifacts.require("Tracking"); // ✅ เพิ่ม Tracking Contract

module.exports = async function (deployer, network, accounts) {
  console.log("🚀 Starting contract deployment...");

  // ✅ 1. Deploy UserRegistry
  await deployer.deploy(UserRegistry);
  const userRegistryInstance = await UserRegistry.deployed();
  console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);

  // ✅ 2. Deploy Product Contract
  await deployer.deploy(Product, userRegistryInstance.address);
  const productInstance = await Product.deployed();
  console.log("✅ Product Contract Deployed at:", productInstance.address);

  // ✅ 3. Deploy RawMilk Contract
  await deployer.deploy(RawMilk, userRegistryInstance.address);
  const rawMilkInstance = await RawMilk.deployed();
  console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);

  // ✅ 4. Deploy CertificationEvent Contract
  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
  const certificationEventInstance = await CertificationEvent.deployed();
  console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);

  // ✅ 5. Deploy ProductLot Contract
  await deployer.deploy(ProductLot, userRegistryInstance.address, rawMilkInstance.address, productInstance.address);
  const productLotInstance = await ProductLot.deployed();
  console.log("✅ ProductLot Contract Deployed at:", productLotInstance.address);

  // ✅ 6. Deploy Tracking Contract (ใหม่)
  await deployer.deploy(Tracking, userRegistryInstance.address, productLotInstance.address);
  const trackingInstance = await Tracking.deployed();
  console.log("✅ Tracking Contract Deployed at:", trackingInstance.address);

  console.log("🎉 Deployment completed successfully!");
};
