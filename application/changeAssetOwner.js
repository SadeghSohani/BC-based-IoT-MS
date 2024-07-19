const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./utils/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('./utils/AppUtil.js');

const channelName = 'mychannel';
const chaincodeName = 'chaincode';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';

var timers = {};
function timer(name) {
    timers[name + '_start'] = performance.now();
}

function timerEnd(name) {
    if (!timers[name + '_start']) return undefined;
    var time = performance.now() - timers[name + '_start'];
    var amount = timers[name + '_amount'] = timers[name + '_amount'] ? timers[name + '_amount'] + 1 : 1;
    var sum = timers[name + '_sum'] = timers[name + '_sum'] ? timers[name + '_sum'] + time : time;
    timers[name + '_avg'] = sum / amount;
    delete timers[name + '_start'];
    return time;
}

async function main(concurrentTxs) {
    try {
        const ccp = buildCCPOrg1();
        const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
        const wallet = await buildWallet(Wallets, walletPath);
        await enrollAdmin(caClient, wallet, mspOrg1);
        await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');
        const gateway = new Gateway();

        try {

            await gateway.connect(ccp, {
                wallet,
                identity: org1UserId,
                discovery: { enabled: true, asLocalhost: true }
            });

            const network = await gateway.getNetwork(channelName);
            const contract = network.getContract(chaincodeName);

            let result = await contract.evaluateTransaction('QueryAllAssets');
            let txs = JSON.parse(result)

            let latencies = [];
            timer("Submit");
            for (let i = 0; i < concurrentTxs; i++) {
                timer(`L#${i}`);
                contract.submitTransaction(
                    'ChangeAssetOwner', `${txs[i].Id}`, 
                    `${txs[i].owner}`, `${txs[i].owner}`
                ).then(result => {
                    console.log(result.toString());
                    latencies.push(timerEnd(`L#${i}`));
                    if (latencies.length === concurrentTxs) {
                        let responseTime = timerEnd("Submit");
                        let sum = 0;
                        let min = Infinity;
                        let max = -Infinity;
                        latencies.forEach(x => {
                            sum += x;
                            if (x < min)
                                min = x;
                            if (x > max)
                                max = x;
                        });
                        console.log(`Latency Min: ${min}`);
                        console.log(`Latency Max: ${max}`);
                        console.log(`Latency AVG: ${sum/latencies.length}`);
                        console.log(`TPS: ${latencies.length/(responseTime/1000)}`);
                    }
                    return;
                });
            }

        } finally {
            gateway.disconnect();
        }
    } catch (error) {
        console.error(`******** FAILED to run the application: ${error}`);
    }
}

main(1000);
