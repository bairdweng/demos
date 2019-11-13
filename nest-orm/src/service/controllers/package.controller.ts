
import { Controller, Get, Post, Body, HttpStatus, Response } from '@nestjs/common';


@Controller()

export class PackageController {
  @Get('package')
  public async package(@Response() res) {
    let process = require('child_process');
    let ls = process.exec('npm run start');
    ls.stdout.on('data', (data) => {
      console.log(`stdout: ${data}`);
    });
    ls.stderr.on('data', (data) => {
      console.error(`stderr: ${data}`);
    });
    ls.on('close', (code) => {
      res.status(HttpStatus.OK).json({ "code": code });
    })
  }
}