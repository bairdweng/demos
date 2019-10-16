package com.example.demo.controller;

// import java.io.IOException;
// import java.io.InputStream;

// import com.example.demo.mapper.UserMapper;
// import org.apache.ibatis.io.Resources;
// import org.apache.ibatis.session.SqlSessionFactory;
// import org.apache.ibatis.session.SqlSessionFactoryBuilder;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
// import com.example.demo.dao.User;
import cn.hutool.core.lang.Dict;
import lombok.extern.slf4j.Slf4j;

@RestController
@Slf4j
public class HomeController {
  // private SqlSessionFactory sessionFactory;
  // @Autowired
  // private UserMapper userMapper;

  // public void setup() throws IOException {
  //   String resource = "config/mybatis-config.xml";
  //   InputStream in = Resources.getResourceAsStream(resource);
  //   sessionFactory = new SqlSessionFactoryBuilder().build(in);
  // }
  /*
  @RequestMapping("/test")
  public String Index() {
    return "这是一个模板工程";
  }

  @RequestMapping("/getdata")
  public String getData() {
    return "2";
  }

  @GetMapping("/user/{id}")
  public Dict getUser(@PathVariable Integer id) {
    User us = userMapper.getUser(id);
    return Dict.create().set("code", 200).set("msg", "成功").set("data", us);
    /*
     * 使用xml try { this.setup(); SqlSession session = sessionFactory.openSession();
     * try { UserMapper um = session.getMapper(UserMapper.class); User us =
     * um.getUser(id); return Dict.create().set("code", 200).set("msg",
     * "成功").set("data", us);
     * 
     * } catch (Exception e) { return Dict.create().set("code", -1).set("msg",
     * "失败").set("data", e.toString()); } } catch (Exception e) { return
     * Dict.create().set("code", -1).set("msg", "失败").set("data", e.toString()); }
     */
  // }
}