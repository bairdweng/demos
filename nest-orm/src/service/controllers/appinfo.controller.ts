
import { Controller, Post, UseInterceptors, UploadedFile, HttpStatus, Response } from '@nestjs/common';

import { FileInterceptor } from '@nestjs/platform-express';

import { unit } from '../common/unit'
import { createWriteStream } from 'fs';
import { join } from 'path';
import { DirType } from '../common/enums'

import { getConnection } from "typeorm";
import { appInfoService } from "../service/appinfo.service"
import { from } from 'rxjs';
@Controller()
export class AppInfoController {
  // 上传iOS的p12文件
  @Post('uploadp12')
  @UseInterceptors(FileInterceptor('file'))
  uploadFile(@UploadedFile() file, @Response() res) {
    if (file.mimetype === 'application/x-pkcs12') {
      let suffix = unit.getFileSuffix(file.originalname)
      if (suffix === null) {
        res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "文件名不规范", null));
        return;
      }
      let fileName = unit.getFileName(suffix);
      const witeFile = createWriteStream(join(__dirname, '../../../', unit.getDir(DirType.p12), fileName))
      witeFile.write(file.buffer)
      res.status(HttpStatus.OK).json(unit.resData(true, "", { "fileName": fileName }));
    }
    else {
      res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "请上传.p12文件", null));
    }
  }
  // 上传iOS描述文件
  @Post('uploadProFile')
  @UseInterceptors(FileInterceptor('file'))
  uploadProFile(@UploadedFile() file, @Response() res) {
    if (file.mimetype === 'application/octet-stream') {
      let suffix = unit.getFileSuffix(file.originalname)
      if (suffix === null) {
        res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "文件名不规范", null));
        return;
      }
      let fileName = unit.getFileName(suffix);
      const witeFile = createWriteStream(join(__dirname, '../../../', unit.getDir(DirType.proFile), fileName))
      witeFile.write(file.buffer)
      res.status(HttpStatus.OK).json(unit.resData(true, "", { "fileName": fileName }));
    }
    else {
      res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "请上传mobileprovision文件", null));
    }
  }
  // 上传icon
  @Post('uploadImage')
  @UseInterceptors(FileInterceptor('file'))
  uploadImage(@UploadedFile() file, @Response() res) {
    if (file.mimetype === 'image/png' || file.mimetype === 'image/jpeg') {
      let suffix = unit.getFileSuffix(file.originalname)
      if (suffix === null) {
        res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "文件名不规范", null));
        return;
      }
      let fileName = unit.getFileName(suffix);
      const witeFile = createWriteStream(join(__dirname, '../../../', unit.getDir(DirType.image), fileName))
      witeFile.write(file.buffer)
      res.status(HttpStatus.OK).json(unit.resData(true, "", { "fileName": fileName }));
    }
    else {
      res.status(HttpStatus.ACCEPTED).json(unit.resData(false, "请上传图片", null));
    }
  }
  // 添加app信息
  @Post('addAppInfo')
  addAppInfo(@Response() res) {

    // console.log(dd);
    let service = new appInfoService();
    service.getInfo().then(ret => {
      res.status(HttpStatus.OK).json(unit.resData(true, "", ret));

    }).catch(e => {
      console.log(e);

      res.status(HttpStatus.ACCEPTED).json(unit.resData(false, String(e), {}));
    });
  }
}