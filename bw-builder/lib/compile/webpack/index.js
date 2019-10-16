"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const webpack = require("webpack");
const md5_1 = require("../update/md5");
const copy_1 = require("../update/copy");
const ready_1 = require("../update/ready");
const fs_extra_1 = require("fs-extra");
const context_1 = require("../update/context");
const config_1 = require("./config");
class Compile {
    static build(name) {
        const weexConfig = new config_1.Config({ env: name, watch: false }).weexConfig;
        webpack(weexConfig, (err, stats) => {
            process.stdout.write(stats.toString({
                colors: true,
                modules: false,
                warnings: false,
                entrypoints: false,
                assets: false,
                hash: false,
                version: false,
                timings: false,
                builtAt: false,
            }));
            if (err == null) {
                this.update(name);
            }
        });
    }
    static watch(name) {
        const context = new context_1.Context();
        fs_extra_1.emptyDirSync(context.distPath);
        const weexConfig = new config_1.Config({ env: name, watch: true }).weexConfig;
        webpack(weexConfig, (err, stats) => {
            process.stdout.write(stats.toString({
                colors: true,
                modules: false,
                warnings: false,
                entrypoints: false,
                assets: false,
                hash: false,
                version: false,
                timings: false,
                builtAt: false,
            }));
        });
    }
    static update(name) {
        //生成md5文件。
        md5_1.Md5.calculate();
        //拷贝文件
        copy_1.Copy.copy(name);
        //移动文件
        ready_1.Ready.ready();
    }
    
}
exports.Compile = Compile;
