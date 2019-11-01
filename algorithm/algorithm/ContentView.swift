//
//  ContentView.swift
//  algorithm
//
//  Created by Baird-weng on 2019/10/23.
//  Copyright © 2019 bairdweng. All rights reserved.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("特殊视图2")) {
                    NavigationLink(destination: ControllerPage<UIKitController>()) {
                        PageRow(title: "UIViewController", subTitle: "打开 UIViewController")
                    }
                    NavigationLink(destination: ControllerPage<UIKitController>()) {
                        PageRow(title: "UIViewController", subTitle: "打开 UIViewController")
                    }
                }
            }
                .listStyle(GroupedListStyle())
                .navigationBarTitle(Text("Example"), displayMode: .large)
                .navigationBarItems(trailing: Button(action: {
                    print("Tap")
                    insertionSort().start()

                }, label: {
                    Text("Right").foregroundColor(.orange)
                    }))

        }
    }
}

#if DEBUG
    struct ContentView_Previews: PreviewProvider {
        static var previews: some View {
            ContentView().colorScheme(.dark)
        }
    }
#endif

