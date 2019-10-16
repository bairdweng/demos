#! /usr/bin/env node
const program = require('commander');
const {debug,run} = require('../unit/unit.js')
const pack = require('../lib/compile/webpack').Compile
program
  .command('debug')
  .description('调试')
  .action(() => {
     debug()
  })
program
  .command('watch <environment>')
  .description('监听')
  .action((environment) => {
    // 主要是为了监听。
    pack.watch(environment)
    run()
  }) 
program
  .command('build <environment>')
  .description('打包')
  .action((environment) => {
    pack.build(environment)
  })   
program.parse(process.argv)