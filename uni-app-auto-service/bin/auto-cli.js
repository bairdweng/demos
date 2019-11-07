const program = require('commander')
const fs = require('fs-extra')
const path = require('path')

const { Package } = require('../lib/package/index')

program
  .version(fs.readJsonSync(path.join(__dirname, '../package.json')).version, '-v, --version')

program
  .command('run')
  .description('uni自动打包')
  .action(() => {
    Package.start()
  })
program.parse(process.argv)