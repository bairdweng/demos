//
//  insertionSort.swift
//  algorithm
//
//  Created by Baird-weng on 2019/10/23.
//  Copyright © 2019 bairdweng. All rights reserved.
//

import UIKit

class insertionSort: NSObject {


    func start() {
       let list = [ 10, -1, 3, 9, 2, 27, 8, 5, 1, 3, 0, 26 ]
       let items1 =  insertionSort(array: list, <)
       let items2 =  insertionSort(array: list, >)
       print(items1)
       print(items2)
    }

    func insertionSort<T>(array: [T], _ isOrderedBefore: (T, T) -> Bool) -> [T] {
        var a = array
        for x in 1..<a.count {
            var y = x
            let temp = a[y]
            while y > 0 && isOrderedBefore(temp, a[y - 1]) {
                a[y] = a[y - 1]
                y -= 1
            }
            a[y] = temp
        }
        return a
    };
}
