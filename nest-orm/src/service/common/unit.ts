
import { DirType } from './enums'
export class unit {
  // 根据时间戳生成文件。
  static getFileName(typeName: string): string {
    let timestamp = new Date().getTime();
    return String(timestamp) + "." + typeName;
  }
  // 定义请求返回的model
  static resData(isOk: boolean, message: string, data: Object) {
    return isOk === true ? {
      code: 0,
      message: '',
      data
    } : {
        code: -1,
        message: message,
        data
      }
  }
  // 获取文件后缀
  static getFileSuffix(fileName: string) {
    let strs = fileName.split(".");
    if (strs.length === 2) {
      return strs[1];
    }
    else {
      return null;
    }
  }
  // 获取文件存放目录
  static getDir(type: DirType) {
    if (type === DirType.p12) {
      return "dir/p12"
    }
    else if (type === DirType.proFile) {
      return "dir/proFile"
    }
    else {
      return "dir/images"
    }
  }
}