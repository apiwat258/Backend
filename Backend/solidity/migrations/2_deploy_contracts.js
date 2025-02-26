const UserRegistry = artifacts.require("UserRegistry");
const CertificationEvent = artifacts.require("CertificationEvent");

module.exports = async function (deployer) {
  await deployer.deploy(UserRegistry);
  const userRegistryInstance = await UserRegistry.deployed();

  await deployer.deploy(CertificationEvent, userRegistryInstance.address);
};