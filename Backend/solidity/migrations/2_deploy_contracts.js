const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");
const Product = artifacts.require("Product");
const ProductLot = artifacts.require("ProductLot");

module.exports = async function (deployer, network, accounts) {
  console.log("🚀 Starting contract deployment...");

  // ✅ 1. Deploy UserRegistry
  await deployer.deploy(UserRegistry);
  const userRegistryInstance = await UserRegistry.deployed();
  console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);

  // ✅ 2. Deploy Product Contract (ต้องมี UserRegistry)
  await deployer.deploy(Product, userRegistryInstance.address);
  const productInstance = await Product.deployed();
  console.log("✅ Product Contract Deployed at:", productInstance.address);

  // ✅ 3. Deploy RawMilk Contract (ต้องมี UserRegistry)
  await deployer.deploy(RawMilk, userRegistryInstance.address);
  const rawMilkInstance = await RawMilk.deployed();
  console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);

  // ✅ 4. Deploy CertificationEvent Contract (ต้องมี UserRegistry)
  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
  const certificationEventInstance = await CertificationEvent.deployed();
  console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);

  // ✅ 5. Deploy ProductLot Contract (ต้องมี UserRegistry, RawMilk, และ Product)
  await deployer.deploy(ProductLot, userRegistryInstance.address, rawMilkInstance.address, productInstance.address);
  const productLotInstance = await ProductLot.deployed();
  console.log("✅ ProductLot Contract Deployed at:", productLotInstance.address);

  console.log("🎉 Deployment completed successfully!");
};
