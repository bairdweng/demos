package com.edu.admin.service;

import com.edu.admin.common.CommonResult;

import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.transaction.annotation.Transactional;

public interface UserService {
  // 用户注册
  @Transactional
  CommonResult register(String username, String password);

  // 用户登录
  @Transactional
  CommonResult login(String username, String password);

  UserDetails loadUserByUsername(String username);

}