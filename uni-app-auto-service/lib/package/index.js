"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let process = require('child_process');
class Package {
    static start() {
        let ls = process.exec('cd platforms/ios/fastlaneDemo && fastlane inHouse');
        ls.stdout.on('data', (data) => {
            console.log(`stdout: ${data}`);
        });
        ls.stderr.on('data', (data) => {
            console.error(`stderr: ${data}`);
        });
        ls.on('close', (code) => {
            console.log(`子进程退出，使用退出码 ${code}`);
        });
    }
}
exports.Package = Package;
