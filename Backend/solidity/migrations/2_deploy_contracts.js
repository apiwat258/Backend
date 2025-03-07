const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilk = artifacts.require("RawMilk");

module.exports = async function (deployer, network, accounts) {
  let userRegistryInstance;
  let certificationEventInstance;
  let rawMilkInstance;

  // ✅ ตรวจสอบว่า UserRegistry ถูกดีพลอยแล้วหรือยัง
  try {
    userRegistryInstance = await UserRegistry.deployed();
    console.log("✅ UserRegistry Contract already deployed at:", userRegistryInstance.address);
  } catch (error) {
    // ถ้ายังไม่ถูกดีพลอย ก็ให้ดีพลอย UserRegistry ใหม่
    await deployer.deploy(UserRegistry);
    userRegistryInstance = await UserRegistry.deployed();
    console.log("✅ UserRegistry Contract Deployed at:", userRegistryInstance.address);
  }

  // ✅ ตรวจสอบว่า CertificationEvent ถูกดีพลอยแล้วหรือยัง
  try {
    certificationEventInstance = await CertificationEvent.deployed();
    console.log("✅ CertificationEvent Contract already deployed at:", certificationEventInstance.address);
  } catch (error) {
    // ถ้ายังไม่ถูกดีพลอย ก็ให้ดีพลอย CertificationEvent ใหม่
    await deployer.deploy(CertificationEvent, userRegistryInstance.address);
    certificationEventInstance = await CertificationEvent.deployed();
    console.log("✅ CertificationEvent Contract Deployed at:", certificationEventInstance.address);
  }

  // ✅ ตรวจสอบว่า RawMilk ถูกดีพลอยแล้วหรือยัง
  try {
    rawMilkInstance = await RawMilk.deployed();
    console.log("✅ RawMilk Contract already deployed at:", rawMilkInstance.address);
  } catch (error) {
    // ถ้ายังไม่ถูกดีพลอย ก็ให้ดีพลอย RawMilk ใหม่
    await deployer.deploy(RawMilk, userRegistryInstance.address);
    rawMilkInstance = await RawMilk.deployed();
    console.log("✅ RawMilk Contract Deployed at:", rawMilkInstance.address);
  }
};
