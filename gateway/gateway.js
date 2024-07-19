const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./utils/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('./utils/AppUtil.js');
const axios = require("axios");
const express = require('express');
const bodyParser = require('body-parser');
const app = express();
app.use(bodyParser.json());

const port = 3000;
const channelName = 'mychannel';
const chaincodeName = 'chaincode';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';
const key = "";
var outLinks = [];

async function main() {
    try {
        const ccp = buildCCPOrg1();

        const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

        const wallet = await buildWallet(Wallets, walletPath);

        await enrollAdmin(caClient, wallet, mspOrg1);

        await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

        const gateway = new Gateway();

        await gateway.connect(ccp, {
            wallet,
            identity: org1UserId,
            discovery: { enabled: true, asLocalhost: true }
        });

        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        await contract.addContractListener(async (event) => {
            if (event.eventName === key) {
                const message = event.payload.toString('utf8');
                const instruction = message.split(':');
                if (instruction[0] === "Send") {
                    outLinks.push(instruction[1])
                } else if (instruction[0] === "Stop") {
                    outLinks = outLinks.filter((c) => { return c !== instruction[1] })
                }
            }
        });

        console.log("Gateway is connected ===================================")
    } catch (error) {
        console.error(`******** FAILED to run the application: ${error}`);
    }
}

app.listen(port, () => {
    main();
    console.log(`Server is running on port ${port}`);
});

app.post('/sensors/data', async (req, res) => {
    for (let link in outLinks) {
        const res = await axios.post(link, req.body);
    }
});