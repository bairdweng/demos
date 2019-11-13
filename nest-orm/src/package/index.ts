import { readdir, createReadStream, writeFile, pathExistsSync, existsSync, readdirSync, statSync } from 'fs-extra'
let process = require('child_process');
export class Package {
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