
package com.example.demo.controller;

import java.util.HashMap;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
public class TestController {

  @RequestMapping(value = "/hello",method = RequestMethod.GET)
  @ResponseBody
  public String helloHtml(HashMap<String, Object> map) {
    map.put("hello", "欢迎进入HTML页面");
    return "/index";
  }
}