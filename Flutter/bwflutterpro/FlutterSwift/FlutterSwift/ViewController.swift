//
//  ViewController.swift
//  FlutterSwift
//
//  Created by Baird-weng on 2019/6/12.
//  Copyright © 2019 bbw. All rights reserved.
//

import UIKit
import Flutter
import FlutterPluginRegistrant
class ViewController: UIViewController, FlutterStreamHandler {

    override func viewDidLoad() {
        super.viewDidLoad()
        self.view.backgroundColor = .red
        // Do any additional setup after loading the view.
    }

    @IBAction func pushFlutterPage(_ sender: Any) {
        let eventName = "com.pages.your/native_get"
        let flutter = BaseFlutterViewController()
//        let methodChannel = FlutterMethodChannel.init(name: eventName, binaryMessenger: flutter)
//
//        
//        let eventChannel = FlutterEventChannel.init(name: "flutterswift", binaryMessenger: flutter)
//        eventChannel.setStreamHandler(self);
        self.navigationController?.pushViewController(flutter, animated: true)
    }
    // 当Flutter 初始化完毕发送消息。
    func onListen(withArguments arguments: Any?, eventSink events: @escaping FlutterEventSink) -> FlutterError? {
        events(["title": "222"]);
        return nil;
    }
    func onCancel(withArguments arguments: Any?) -> FlutterError? {

      print("========", arguments ?? []);

      return nil
    }
}

