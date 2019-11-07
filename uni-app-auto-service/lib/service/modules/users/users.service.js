"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const common_1 = require("@nestjs/common");
const core_1 = require("@nestjs/core");
const uuid = require("uuid");
const sqlite_1 = require("../../db/sqlite");
let UsersService = class UsersService {
    getAllUsers() {
        return new Promise((resolve, reject) => {
            sqlite_1.db.all("SELECT * FROM user", function (err, rows) {
                return !err ?
                    resolve(rows) :
                    reject(new core_1.HttpException(err, 500));
            });
        });
    }
    getUser(value) {
        return new Promise((resolve, reject) => {
            sqlite_1.db.get("SELECT * FROM user WHERE id = ?", [value], (err, row) => {
                return !err ?
                    resolve(row) :
                    reject(new core_1.HttpException(err, 500));
            });
        });
    }
    createUser(username, email, password) {
        return new Promise((resolve, reject) => {
            sqlite_1.db.run("INSERT INTO user (id, username, email, password, role)" +
                "VALUES (?, ?, ?, ?, 'user')", [uuid.v1().replace(/-/g, ""), username, email, password], (err) => {
                return !err ?
                    resolve({ 'message': 'User has been registered' }) :
                    reject(new core_1.HttpException(err, 500));
            });
        });
    }
    updateUser(id, username, email, password, role) {
        return new Promise((resolve, reject) => {
            sqlite_1.db.run("UPDATE user SET username=?, email=?, password=?, role=?" +
                "WHERE(id = ?);", [username, email, password, role, id], (err) => {
                return !err ?
                    resolve({ 'message': 'User ' + id + ' has been updated successfully' }) :
                    reject(new core_1.HttpException(err, 500));
            });
        });
    }
    deleteUser(id) {
        return new Promise((resolve, reject) => {
            sqlite_1.db.run("DELETE From user WHERE id = ?", [id], (err) => {
                return !err ?
                    resolve({ 'message': 'User ' + id + ' has been deleted' }) :
                    reject(new core_1.HttpException(err, 500));
            });
        });
    }
};
UsersService = __decorate([
    common_1.Component()
], UsersService);
exports.UsersService = UsersService;
