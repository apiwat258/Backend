const CertificationEvent = artifacts.require("CertificationEvent");
const RawMilkSupplyChain = artifacts.require("RawMilkSupplyChain");

module.exports = function (deployer) {
  deployer.deploy(CertificationEvent);
  deployer.deploy(RawMilkSupplyChain);
};
