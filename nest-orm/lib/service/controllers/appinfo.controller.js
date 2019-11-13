"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
const common_1 = require("@nestjs/common");
const platform_express_1 = require("@nestjs/platform-express");
let AppInfoController = class AppInfoController {
    uploadFile(file) {
        console.log(file);
    }
    async testpackage(res) {
        res.status(common_1.HttpStatus.OK).json({ "code": 1 });
    }
};
__decorate([
    common_1.Post('upload'),
    common_1.UseInterceptors(platform_express_1.FileInterceptor('file')),
    __param(0, common_1.UploadedFile())
], AppInfoController.prototype, "uploadFile", null);
__decorate([
    common_1.Get('testpackage'),
    __param(0, common_1.Response())
], AppInfoController.prototype, "testpackage", null);
AppInfoController = __decorate([
    common_1.Controller()
], AppInfoController);
exports.AppInfoController = AppInfoController;
