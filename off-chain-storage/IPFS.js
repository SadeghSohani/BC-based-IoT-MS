const axios = require("axios");
const FormData = require("form-data");
const fs = require("fs");
const JWT = "";

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

async function upload() {
    try {
        const formData = new FormData();

        const file = fs.createReadStream("./1.txt");
        formData.append("file", file);

        const pinataMetadata = JSON.stringify({
            name: "1.txt",
        });
        formData.append("pinataMetadata", pinataMetadata);

        const pinataOptions = JSON.stringify({
            cidVersion: 1,
        });
        formData.append("pinataOptions", pinataOptions);

        timer("Submit");
        const res = await axios.post(
            "https://api.pinata.cloud/pinning/pinFileToIPFS",
            formData,
            {
                headers: {
                    Authorization: `Bearer ${JWT}`,
                },
            }
        );
        console.log(res.data);
        console.log(timerEnd("Submit"));
    } catch (error) {
        console.log(error);
    }
}

async function download() {
    try {
        timer("Submit");
        const res = await fetch(
            "https://gateway.pinata.cloud/ipfs/bafybeiedu3evepxtpa52swq2uv6jq2rdwiagf6x77dh5n2vgn3hipoalxu"
        );
        const resData = await res.text();
        // console.log(resData);
        console.log(timerEnd("Submit"));
    } catch (error) {
        console.log(error);
    }
}

// upload();
download();