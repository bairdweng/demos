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
let UsersController = class UsersController {
    constructor(usersService) {
        this.usersService = usersService;
    }
    async getAllUsers(res) {
        const users = await this.usersService.getAllUsers();
        res.status(common_1.HttpStatus.OK).json(users);
    }
    async getUser(res, id) {
        const user = await this.usersService.getUser(id);
        res.status(common_1.HttpStatus.OK).json(user);
    }
    async createUser(res, username, email, password) {
        const result = await this.usersService.createUser(username, email, password);
        res.status(common_1.HttpStatus.CREATED).json(result);
    }
    async updateUser(res, id, username, email, password, role) {
        const result = await this.usersService.updateUser(id, username, email, password, role);
        res.status(common_1.HttpStatus.ACCEPTED).json(result);
    }
    async deleteUser(res, id) {
        const result = await this.usersService.deleteUser(id);
        res.status(common_1.HttpStatus.ACCEPTED).json(result);
    }
};
__decorate([
    common_1.Get(),
    __param(0, common_1.Response())
], UsersController.prototype, "getAllUsers", null);
__decorate([
    common_1.Get('/:id'),
    __param(0, common_1.Response()),
    __param(1, common_1.Param('id'))
], UsersController.prototype, "getUser", null);
__decorate([
    common_1.Post(),
    __param(0, common_1.Response()),
    __param(1, common_1.Body('username')),
    __param(2, common_1.Body('email')),
    __param(3, common_1.Body('password'))
], UsersController.prototype, "createUser", null);
__decorate([
    common_1.Put('/:id'),
    __param(0, common_1.Response()),
    __param(1, common_1.Param('id')),
    __param(2, common_1.Body('username')),
    __param(3, common_1.Body('email')),
    __param(4, common_1.Body('password')),
    __param(5, common_1.Body('role'))
], UsersController.prototype, "updateUser", null);
__decorate([
    common_1.Delete('/:id'),
    __param(0, common_1.Response()),
    __param(1, common_1.Param('id'))
], UsersController.prototype, "deleteUser", null);
UsersController = __decorate([
    common_1.Controller('user')
], UsersController);
exports.UsersController = UsersController;
