import 'package:flutter/material.dart';

class userinfo extends StatefulWidget {
  @override
  userinfoState createState() => userinfoState();
}

class userinfoState extends State<userinfo> {
  String _showTitle = "~~~";

  @override
  Widget build(BuildContext context) {
    Object arg = ModalRoute.of(context).settings.arguments;
    var params = new Map.from(arg);
    return Scaffold(
      appBar: new AppBar(
          title: new Text(params['userName']),
          leading: new IconButton(
            tooltip: 'Previous choice',
            icon: const Icon(Icons.arrow_back),
            onPressed: () {
              popPage();
            },
          )),
      body: Container(
        color: Colors.grey[300],
        margin: const EdgeInsets.all(0),
        child: ListView(children: getItems()),
      ),
      backgroundColor: Colors.white,
    );
  }

  @override
  void initState() {
    // TODO: implement initState
    super.initState();
  }

  // 头部。
  getHeaderView() {
    return GestureDetector(
      child: Column(
        children: <Widget>[
          Stack(
            alignment: AlignmentDirectional.topCenter,
            children: <Widget>[
              //矩形
              Container(
                margin: const EdgeInsets.fromLTRB(10, 50, 10, 10),
                height: 200,
                child: Column(
                  children: <Widget>[
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: <Widget>[
                        Container(
                          margin: const EdgeInsets.fromLTRB(0, 50, 0, 0),
                          child: Text(
                            "用户名",
                            style: TextStyle(color: Colors.black, fontSize: 20),
                          ),
                        )
                      ],
                    ),
                    Container(
                      margin: const EdgeInsets.fromLTRB(0, 10, 0, 0),
                      width: 300,
                      child: Text(
                        "我就是个性签名~~~~~~~~",
                        overflow: TextOverflow.ellipsis,
                        textAlign: TextAlign.center,
                        style: TextStyle(color: Colors.black38, fontSize: 14),
                        maxLines: 1,
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: <Widget>[
                        Container(
                          margin: EdgeInsets.all(20),
                          child: Column(
                            children: <Widget>[
                              Text("关注",
                                  style: TextStyle(color: Colors.black38)),
                              Container(
                                  padding: EdgeInsets.fromLTRB(0, 5, 0, 0),
                                  child: Text('123',
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 20,
                                      )))
                            ],
                          ),
                        ),
                        Container(
                          margin: EdgeInsets.all(20),
                          child: Column(
                            children: <Widget>[
                              Text("粉丝",
                                  style: TextStyle(color: Colors.black38)),
                              Container(
                                  padding: EdgeInsets.fromLTRB(0, 5, 0, 0),
                                  child: Text('123',
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 20,
                                      )))
                            ],
                          ),
                        ),
                        Container(
                          margin: EdgeInsets.all(20),
                          child: Column(
                            children: <Widget>[
                              Text("赞",
                                  style: TextStyle(color: Colors.black38)),
                              Container(
                                  padding: EdgeInsets.fromLTRB(0, 5, 0, 0),
                                  child: Text('123',
                                      style: TextStyle(
                                        color: Colors.black,
                                        fontSize: 20,
                                      )))
                            ],
                          ),
                        )
                      ],
                    )
                  ],
                ),
                decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(5)),
              ),
              //头像
              Container(
                margin: EdgeInsets.fromLTRB(0, 10, 0, 0),
                height: 80,
                width: 80,
                decoration: BoxDecoration(
                    color: Colors.green,
                    borderRadius: BorderRadius.circular(40)),
              ),
            ],
          ),
          Row(
            children: <Widget>[],
          )
        ],
      ),
    );
  }

  getItems() {
    var items = <Widget>[];
    items.add(getHeaderView());
    for (var i = 0; i < 20; i++) {
      items.add(GestureDetector(
        child: Container(
          margin: const EdgeInsets.fromLTRB(0, 1, 0, 0),
          color: Colors.blue[600],
          alignment: Alignment.center,
          child: Text(_showTitle,
              style: Theme.of(context)
                  .textTheme
                  .display1
                  .copyWith(color: Colors.white)),
        ),
        onTap: () {
          print(i);
          updateText();
        },
      ));
    }
    return items;
  }

  void updateText() {
    setState(() {
      _showTitle = "我被更改了2";
    });
  }

  void popPage() {
    Navigator.of(context).pop();
  }
}
