
package com.edu.admin.service.impl;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.naming.AuthenticationException;

import com.edu.admin.common.CommonResult;
import com.edu.admin.common.JwtTokenUtil;
import com.edu.admin.dao.FxUser;
import com.edu.admin.dao.FxUserExample;
import com.edu.admin.mapper.FxUserMapper;
import com.edu.admin.model.FxUserDetails;
import com.edu.admin.service.UserService;
import org.springframework.stereotype.Service;
import cn.hutool.core.collection.CollectionUtil;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.password.PasswordEncoder;

@Service
public class UserServiceImpl implements UserService {
  @Autowired
  private JwtTokenUtil jwtTokenUtil;
  @Autowired
  private FxUserMapper mapper;
  @Autowired
  private PasswordEncoder passwordEncoder;

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
    user.setPassword(passwordEncoder.encode(password));
    mapper.insert(user);
    return CommonResult.success("注册成功", "message");
  }

  @Override
  public CommonResult login(String username, String password) {
    FxUserExample example = new FxUserExample();
    example.createCriteria().andUserNameEqualTo(username);
    List<FxUser> users = mapper.selectByExample(example);
    if (CollectionUtil.isEmpty(users)) {
      return CommonResult.failed("用户不存在");
    }
    FxUser user = users.get(0);
    if (passwordEncoder.matches(password, user.getPassword())) {
      FxUserDetails details = new FxUserDetails(user);
      UsernamePasswordAuthenticationToken authentication = new UsernamePasswordAuthenticationToken(details, null,
          details.getAuthorities());
      SecurityContextHolder.getContext().setAuthentication(authentication);
      Map<String, Object> tokenMap = new HashMap<>();
      // mmp.put("key", "value");
      tokenMap.put("token", this.getToken(username));
      // 这里可能需要生成token。
      return CommonResult.success(tokenMap, "登录成功");
    } else {
      return CommonResult.failed("密码错误");
    }
  }

  private String getToken(String username) {
    String token = null;
    UserDetails userDetails = loadUserByUsername(username);
    token = jwtTokenUtil.generateToken(userDetails);
    return token;
  }

  @Override
  public UserDetails loadUserByUsername(String username) {
    FxUserExample example = new FxUserExample();
    example.createCriteria().andUserNameEqualTo(username);
    List<FxUser> users = mapper.selectByExample(example);
    FxUser user = users.get(0);
    FxUserDetails details = new FxUserDetails(user);
    return details;
  }

}