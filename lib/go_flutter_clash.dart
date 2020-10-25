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

  /// 初始化clash
  static Future<void> init(String homeDir) =>
      _channel.invokeMethod('init', [homeDir]);

  /// 启动clash
  static Future<void> start(
    String profile,
    FlutterClashConfig fcc,
  ) =>
      _channel.invokeMethod('start', [profile, jsonEncode(fcc)]);
}
