// You have generated a new plugin project without
// specifying the `--platforms` flag. A plugin project supports no platforms is generated.
// To add platforms, run `flutter create -t plugin --platforms <platforms> .` under the same
// directory. You can also find a detailed instruction on how to add platforms in the `pubspec.yaml` at https://flutter.dev/docs/development/packages-and-plugins/developing-packages#plugin-platforms.

import 'dart:async';
import 'dart:convert';

import 'package:flutter/services.dart';

import 'model/flutter_clash_config_model.dart';

class GoFlutterClash {
  static const MethodChannel _channel = const MethodChannel('go_flutter_clash');
  static final Map<String, Function(dynamic)> _callHanders = {};

  /// 初始化clash
  static Future<void> init(String homeDir) async {
    _channel.setMethodCallHandler((MethodCall call) async {
      if (_callHanders.containsKey(call.method)) {
        Function.apply(_callHanders[call.method], [call.arguments]);
      }
    });
    return _channel.invokeMethod('init', homeDir);
  }

  /// 启动clash
  static Future<void> start(
    String profile,
    FlutterClashConfig fcc,
  ) =>
      _channel.invokeMethod('start', [profile, jsonEncode(fcc)]);

  /// 当前开启状态
  static Future<bool> status() => _channel.invokeMethod('status');

  /// 实时网速回调
  static void trafficHandler(Function(dynamic) callback) {
    _callHanders["trafficHandler"] = callback;
  }
}
