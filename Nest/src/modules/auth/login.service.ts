import { Component } from '@nestjs/common';
import { HttpException } from '@nestjs/core';
import { AuthHelper } from './auth.helper';
import { db } from '../../db/sqlite';
let process = require('child_process');

@Component()
export class LoginService {

  constructor(
    private authHelper: AuthHelper
  ) { }

  /*******************************************************
   * Basic Login with credentials
   *******************************************************/
  public login(email, password) {
    console.log("==========================", email);

    process.exec('ls', function (error, stdout, d) {
      console.log(error, stdout, d);
    });
    process.on('message', function (m) {
      console.log('Son Listen:', m);
    });

    if (!email) {
      return (new HttpException("Email is required", 422))
    }
    if (!password) {
      return (new HttpException("Password is required", 422))
    }

    return new Promise((resolve, reject) => {
      db.all(
        "SELECT * FROM user WHERE email = ?", [email], (err, user) => {

          if (err) {
            return reject(new HttpException(err, 503))
          }

          else {
            if (user.length === 0) {
              return reject(new HttpException("User does not exist", 401))
            }
            else {
              if (password != user[0].password) {
                return reject(new HttpException("Incorrect password", 401))
              }
              else {
                return resolve(this.authHelper.genToken(user[0]));
              }
            }
          }
        });
    });
  }
}