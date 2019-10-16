import 'package:flutter/material.dart';
import 'dart:ui' as ui;
import 'package:bwflutterpro/tabbar/tabbar.dart';
import 'package:bwflutterpro/toolpage/toolpage.dart';

import 'package:bwflutterpro/userinfo/userinfo.dart';
import 'package:bwflutterpro/helper/NativeRouter.dart';
import 'package:bwflutterpro/webview/hmwebview.dart';
void main() => runApp(_widgetForRoute(ui.window.defaultRouteName));

// 根据iOS端传来的route跳转不同界面
Widget _widgetForRoute(String route) {
  switch (route) {
    case 'myApp':
      return new MyApp();
    case 'toolpage':
      return toolpage();
    default:
      print('---------进入了其它------');
      return MyApp();
  }
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: '我的app',
        theme: ThemeData(
          primarySwatch: Colors.blue,
        ),
        home: new tabbar(),
        routes: <String, WidgetBuilder>{
          '/userinfo': (BuildContext context) => new userinfo(),
        });
  }

}

class normal extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    print('--------------------------------');

    return GestureDetector(
      child: Container(
        width: 100,
        height: 100,
        color: Colors.red,
      ),
      onTap: () {
//        NativeRouter.popVc();
      },
    );

//    return   Container(
//      width: 100,
//      height: 100,
//      color: Colors.green,
//      child: new RaisedButton(onPressed: null),
//    );
  }
}
