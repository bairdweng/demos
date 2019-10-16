import 'package:flutter/services.dart';

class NativeRouter extends Object {
  //注册一个服务
  static const methodChannel = const MethodChannel('com.pages.your/native_get');

  //回调native
  static popVc() async {
    await methodChannel.invokeListMethod('popvc', '参数');
  }
}
