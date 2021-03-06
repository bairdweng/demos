"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const context_1 = require("./context");
const fs_extra_1 = require("fs-extra");
const path_1 = require("path");
const AdmZip = require('adm-zip');
class Copy {
    static copy(configName) {
        const context = new context_1.Context();
        fs_extra_1.emptyDirSync(context.android);
        fs_extra_1.emptyDirSync(context.ios);
        fs_extra_1.copySync(context.configFilePath, path_1.join(context.android, context.configFileName));
        fs_extra_1.copySync(context.md5FilePath, path_1.join(context.android, context.md5FileName));
        fs_extra_1.copySync(context.configFilePath, path_1.join(context.ios, context.configFileName));
        fs_extra_1.copySync(context.md5FilePath, path_1.join(context.ios, context.md5FileName));
        const weexboxConfig = require(context.weexboxConfigPath);
        if (weexboxConfig[configName].imagePublicPath === 'bundle://') {
            fs_extra_1.copySync(context.staticPath, context.androidStaticPath);
            fs_extra_1.copySync(context.staticPath, context.iosStaticPath);
        }
        const zip = new AdmZip();
        zip.addLocalFolder(context.wwwFolderPath);
        zip.deleteFile(context.configFileName);
        zip.deleteFile(context.md5FileName);
        zip.writeZip(path_1.join(context.android, 'www.zip'));
        zip.writeZip(path_1.join(context.ios, 'www.zip'));
    }
}
exports.Copy = Copy;
