
package com.edu.admin.service.impl;

import java.util.List;

import com.edu.admin.common.CommonResult;
import com.edu.admin.dao.FxUser;
import com.edu.admin.dao.FxUserExample;
import com.edu.admin.mapper.FxUserMapper;
import com.edu.admin.service.UserService;
import org.springframework.stereotype.Service;
// import org.springframework.security.crypto.password.PasswordEncoder;

import cn.hutool.core.collection.CollectionUtil;

import org.springframework.beans.factory.annotation.Autowired;

@Service
public class UserServiceImpl implements UserService {
  @Autowired
  private FxUserMapper mapper;
  // @Autowired
  // private PasswordEncoder passwordEncoder;

  @Override
  public CommonResult register(String username, String password) {

    FxUserExample example = new FxUserExample();
    example.createCriteria().andUserNameEqualTo(username);
    List<FxUser> users = mapper.selectByExample(example);
    if (!CollectionUtil.isEmpty(users)) {
      return CommonResult.failed("该用户已经存在");
    }
    FxUser user = new FxUser();
    user.setUserName(username);
    // user.setPassword(passwordEncoder.encode(password));
    user.setPassword(password);
    mapper.insert(user);
    return CommonResult.success("注册成功", "message");
  }

}