export const sxw = function () {
  /*
  example1();
  example2();
  example3();
  example4();
  example5();
  example6();
  example7();
  example8();
  example9();
  example10();
  example11();
  example12();
  */
  // example13();
  // example14();
  example15();

}
// 展现变量跟函数提升
function example1() {
  var scope = "global scope";
  function checkscope() {
    var scope = "local scope";
    function f() {
      return scope;
    }
    return f();
  }
  checkscope();
  console.log(`example1========${scope}`);
}
function example2() {
  var scope = "global scope";
  function checkscope() {
    var scope = "local scope";
    function f() {
      return scope;
    }
    return f;
  }
  checkscope()();
  console.log(`example2========${scope}`);
}

// 内存存放
function example3() {
  var a = 20;
  var b = a;
  b = 30;
  // a是20，为什么不是30
  console.log(`example3===========${a}`);
}
function example4() {
  var a = { name: '前端开发' }
  var b = a;
  b.name = '进阶';
  // 进阶，为什么不是前端开发
  console.log(`example4===========${a.name}`);
}
function example5() {
  var a = { name: '前端开发' }
  var b = a;
  a = null;
  // null，为什么不是前端开发
  console.log(`example5===========${a}`);
}
// 思考
function example6() {
  var a = { n: 1 };
  var b = a;
  a.x = a = { n: 2 };
  console.log(`example6===========a.x=${JSON.stringify(a)}  b=${JSON.stringify(b)}`);
}

// 闭包
function example7() {
  var data = [];
  for (var i = 0; i < 3; i++) {
    data[i] = function () {
      console.log(`example7==============${i}`);
    };
  }
  data[0]();
  data[1]();
  data[2]();
}
function example8() {
  var data = [];
  for (let i = 0; i < 3; i++) {
    data[i] = function () {
      console.log(`example8==============${i}`);
    };
  }
  data[0]();
  data[1]();
  data[2]();
}

// --------this
function example9() {
  function baz() {
    // 当前调用栈是：baz
    // 因此，当前调用位置是全局作用域

    console.log("baz");
    bar(); // <-- bar的调用位置
  }

  function bar() {
    // 当前调用栈是：baz --> bar
    // 因此，当前调用位置在baz中

    console.log("bar");
    foo(); // <-- foo的调用位置
  }

  function foo() {
    // 当前调用栈是：baz --> bar --> foo
    // 因此，当前调用位置在bar中

    console.log("foo");
  }

  baz();
}
function example10() {
  function test() {
    console.log(`example10------`, this);
  };
  test();
  // console.log(`example10 this.a======${a}`);

  // console.log(`example10 this.b======${this.b}`);

}
// apply,call
function example11() {
  function foo(a, b) {
    console.log("example11" + "a:" + a + "，b:" + b);
  }

  // 把数组”展开“成参数
  foo.apply(null, [2, 3]); // a:2，b:3

  // 使用bind(..)进行柯里化
  var bar = foo.bind(null, 2);
  bar(3); // a:2，b:3 
}
// 拷贝
function example12() {
  // 木易杨
  let a = {
    name: "muyiy",
    book: {
      title: "You Don't Know JS",
      price: "45"
    }
  }
  let b = Object.assign({}, a);
  a.book.title = 'I will be change'
  console.log('example12 b1 =======', b);
}

function example13() {
  let list = ['10', '3', '4', '15', '6'].map((i) => {
    return parseInt(i);
  });
  console.log(list);
}

function example14() {
  function changeObjProperty(o) {
    o.siteUrl = "http://www.baidu.com"
    // new 指向了新的内存空间，新对象的改变跟旧对象无关。
    o = Object()
    o.siteUrl = "http://www.google.com"
  }
  let webSite = new Object();
  changeObjProperty(webSite);
  console.log(webSite.siteUrl);
}

function example15() {
  function Foo() {
    Foo.a = function () {
      console.log(1)
    }
    this.a = function () {
      console.log(2)
    }
  }
  // 以上只是 Foo 的构建方法，没有产生实例，此刻也没有执行

  Foo.prototype.a = function () {
    console.log(3)
  }
  // 现在在 Foo 上挂载了原型方法 a ，方法输出值为 3

  Foo.a = function () {
    console.log(4)
  }
  // 现在在 Foo 上挂载了直接方法 a ，输出值为 4

  Foo.a();
  // 立刻执行了 Foo 上的 a 方法，也就是刚刚定义的，所以
  // # 输出 4
  let obj = new Foo();
  /* 这里调用了 Foo 的构建方法。Foo 的构建方法主要做了两件事：
1. 将全局的 Foo 上的直接方法 a 替换为一个输出 1 的方法。
2. 在新对象上挂载直接方法 a ，输出值为 2。
*/
  obj.a();
  // 因为有直接方法 a ，不需要去访问原型链，所以使用的是构建方法里所定义的 this.a，
  // # 输出 2
  Foo.a();
  // 构建方法里已经替换了全局 Foo 上的 a 方法，所以
  // # 输出 1

}
