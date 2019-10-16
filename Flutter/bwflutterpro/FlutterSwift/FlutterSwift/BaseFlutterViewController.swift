//
//  BaseFlutterViewController.swift
//  FlutterSwift
//
//  Created by Baird-weng on 2019/7/12.
//  Copyright © 2019 bbw. All rights reserved.
//

import UIKit
import  Flutter
class BaseFlutterViewController: FlutterViewController {

    
//   init() {
//       let appDelegate = UIApplication.shared.delegate as! AppDelegate    
////       super.init(engine: appDelegate.flutterEngine, nibName: nil, bundle: nil)
////        super.init() 
//    }
    
//    required init?(coder aDecoder: NSCoder) {
//        fatalError("init(coder:) has not been implemented")
//    }
//    init() {
//        super.init(engine: (UIApplication.shared.delegate as! AppDelegate).flutterEngine, nibName: nil, bundle: nil)
//    }
//    
//    required init?(coder aDecoder: NSCoder) {
//        fatalError("init(coder:) has not been implemented")
//    }
    override func viewDidLoad() {
        self.setInitialRoute("")
        super.viewDidLoad()
        
        let methodChannel = FlutterMethodChannel.init(name: "222222", binaryMessenger: self)
        methodChannel.setMethodCallHandler { (call, result) in
            
        }

        // Do any additional setup after loading the view.
    }
    deinit {
      print("视图控制器被释放了---------------------")
    }
    

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destination.
        // Pass the selected object to the new view controller.
    }
    */

}
