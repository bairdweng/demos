package com.edu.admin.controller;

import java.util.List;

import com.edu.admin.common.CommonResult;
import com.edu.admin.dao.FxUser;
// import com.edu.admin.dao.LoginInfo;
import com.edu.admin.mapper.FxUserMapper;
import com.github.pagehelper.PageHelper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import cn.hutool.core.lang.Dict;
import io.swagger.annotations.ApiOperation;
import io.swagger.v3.oas.annotations.parameters.RequestBody;
import com.edu.admin.service.UserService;

@RestController
public class UserController {
  @Autowired
  private FxUserMapper mapper;

  @Autowired
  private UserService userService;

  /*
   * @ApiOperation("登录")
   * 
   * @RequestMapping(value = "/login", method = RequestMethod.POST)
   * 
   * @ResponseBody public Dict login(@RequestBody LoginInfo info) { // FxUser us =
   * mapper.selectByPrimaryKey(Long.valueOf(id)); // String d = "123123"; // d =
   * info. // LoginInfo info = new LoginInfo(); // info.passWord = "2"; return
   * Dict.create().set("code", 200).set("msg", "成功").set("data", info); }
   */

  @ApiOperation("注册")
  @RequestMapping(value = "/register", method = RequestMethod.POST)
  public CommonResult register(@RequestParam String username, @RequestParam String password) {
    return userService.register(username, password);
  }

  @ApiOperation("获取用户信息")
  @RequestMapping(value = "/user", method = RequestMethod.GET)
  public Dict getUser(@RequestParam(value = "id", defaultValue = "10000") String id) {
    FxUser us = mapper.selectByPrimaryKey(Long.valueOf(id));
    return Dict.create().set("code", 200).set("msg", "成功").set("data", us);
  }

  @ApiOperation("获取所有用户信息")
  @RequestMapping(value = "/user/all", method = RequestMethod.GET)
  public Dict getAllUser(@RequestParam(value = "pageSize", defaultValue = "5") Integer pageSize,
      @RequestParam(value = "pageNum", defaultValue = "1") Integer pageNum) {
    PageHelper.startPage(pageNum, pageSize);
    // List<FxUser> uss = mapper.selectAll();

    return Dict.create().set("code", 200).set("msg", "成功").set("data", null);
  }
}