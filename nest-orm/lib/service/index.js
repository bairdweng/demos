"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = require("@nestjs/core");
const express = require("express");
const app_modules_1 = require("./app.modules");
const server = express();
async function bootstrap() {
    const app = await core_1.NestFactory.create(app_modules_1.ApplicationModule);
    await app.listen(3000);
}
bootstrap();
