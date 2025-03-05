const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");

module.exports = async function (deployer, network, accounts) {
  let userRegistryInstance;
  let certificationEventInstance;
  let rawMilkInstance;

  // ✅ ดีพลอย UserRegistry ก่อน
  await deployer.deploy(UserRegistry);
  userRegistryInstance = await UserRegistry.deployed();
  console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);

  // ✅ ดีพลอย CertificationEvent โดยใช้ address ของ UserRegistry
  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
  certificationEventInstance = await CertificationEvent.deployed();
  console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);

  // ✅ ดีพลอย RawMilk โดยใช้ address ของ UserRegistry
  await deployer.deploy(RawMilk, userRegistryInstance.address);
  rawMilkInstance = await RawMilk.deployed();
  console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);
};
