import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../widgets/loader.dart';

class LoaderService {
  static void show(BuildContext context, String message) {
    if (message.isEmpty){
      showDialog(context: context, barrierDismissible: false, builder: (_) => AppLoader(message: message));
    }
    else {
      showDialog(context: context, barrierDismissible: false, builder: (_) => const AppLoader(message: "Loading...",));
    }
  }
  static void hide(BuildContext context) {
    Navigator.of(context, rootNavigator: true).pop();
  }
}