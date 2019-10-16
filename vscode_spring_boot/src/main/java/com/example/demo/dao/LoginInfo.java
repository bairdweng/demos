package com.example.demo.dao;

import io.swagger.annotations.ApiModelProperty;
import lombok.Data;

@Data
public class LoginInfo {
  @ApiModelProperty(value = "用户名")
  private String userName;
  @ApiModelProperty(value = "密码")
  private String passWord;
}