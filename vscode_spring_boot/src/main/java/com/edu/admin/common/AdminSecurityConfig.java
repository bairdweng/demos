
package com.edu.admin.common;

import com.edu.admin.service.UserService;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.core.userdetails.UserDetailsService;

@Configuration
@EnableWebSecurity
public class AdminSecurityConfig extends SecurityConfig {

  @Autowired
  private UserService adminService;

  @Bean
  public UserDetailsService userDetailsService() {
    // 获取登录用户信息
    return username -> adminService.loadUserByUsername(username);
  }
}
