package com.edu.admin.service;

import com.edu.admin.common.CommonResult;

import org.springframework.transaction.annotation.Transactional;

import cn.hutool.core.lang.Dict;

public interface UserService {
  // 用户注册
  @Transactional
  CommonResult register(String username, String password);
}