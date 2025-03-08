const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");

module.exports = async function (deployer, network, accounts) {
  console.log("🚀 Starting contract deployment...");

  // ✅ ดีพลอย UserRegistry ใหม่
  await deployer.deploy(UserRegistry);
  const userRegistryInstance = await UserRegistry.deployed();
  console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);

  // ✅ ดีพลอย CertificationEvent ใหม่ โดยใช้ address ของ UserRegistry
  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
  const certificationEventInstance = await CertificationEvent.deployed();
  console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);

  // ✅ ดีพลอย RawMilk ใหม่ โดยใช้ address ของ UserRegistry
  await deployer.deploy(RawMilk, userRegistryInstance.address);
  const rawMilkInstance = await RawMilk.deployed();
  console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);

  console.log("🎉 All contracts deployed successfully!");
};
