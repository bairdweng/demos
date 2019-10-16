import 'package:flutter/material.dart';

class toolpage extends StatefulWidget {
  @override
  toolpageState createState() => toolpageState();
}

class toolpageState extends State<toolpage> {
  String _showTitle = "我是谁？";

  @override
  Widget build(BuildContext context) {
    return ListView(children: getItems());
  }

  @override
  void initState() {
    // TODO: implement initState
    super.initState();
  }

  getItems() {
    var items = <Widget>[];
    for (var i = 0; i < 2; i++) {
      items.add(GestureDetector(
        child: Container(
          margin: const EdgeInsets.fromLTRB(10, 1, 10, 0),
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
}
