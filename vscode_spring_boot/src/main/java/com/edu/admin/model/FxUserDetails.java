
package com.edu.admin.model;

import com.edu.admin.dao.FxUser;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.Arrays;
import java.util.Collection;

public class FxUserDetails implements UserDetails {
  /**
   *
   */
  private static final long serialVersionUID = 1L;
  private FxUser user;

  public FxUserDetails(FxUser details) {
    this.user = details;
  }

  @Override
  public Collection<? extends GrantedAuthority> getAuthorities() {
    // 返回当前用户的权限
    return Arrays.asList(new SimpleGrantedAuthority("TEST"));
  }

  @Override
  public String getPassword() {
    return user.getPassword();
  }

  @Override
  public String getUsername() {
    return user.getUserName();
  }

  @Override
  public boolean isAccountNonExpired() {
    return true;
  }

  @Override
  public boolean isAccountNonLocked() {
    return true;
  }

  @Override
  public boolean isCredentialsNonExpired() {
    return true;
  }

  @Override
  public boolean isEnabled() {
    return true;
  }
}