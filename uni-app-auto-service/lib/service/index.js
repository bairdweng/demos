"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = require("@nestjs/core");
const express = require("express");
const app_modules_1 = require("./app.modules");
const server = express();
const app = core_1.NestFactory.create(app_modules_1.ApplicationModule, server);
app.listen(3000, () => {
    console.log('Typescript Nest app & Express server running on port 3000');
});
