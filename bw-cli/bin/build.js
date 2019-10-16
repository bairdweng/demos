#! /usr/bin/env node
const fs = require('fs-extra')
const program = require('commander')
const path = require('path')
program
  .version(fs.readJsonSync(path.join(__dirname, '../package.json')).version, '-v, --version')
  .command('create <projectName>')
  .description('使用 BW 创建工程')
  .action((projectName) => {
    require('../lib/creator')(projectName)
  })

program.parse(process.argv)