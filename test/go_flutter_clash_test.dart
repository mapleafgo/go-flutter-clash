import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_flutter_clash/go_flutter_clash.dart';
import 'package:go_flutter_clash/model/flutter_clash_config_model.dart';

void main() {
  const MethodChannel channel = MethodChannel('go_flutter_clash');

  TestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    channel.setMockMethodCallHandler((MethodCall methodCall) async {
      return '42';
    });
  });

  tearDown(() {
    channel.setMockMethodCallHandler(null);
  });

  test('start', () async {
    // expect(await GoFlutterClash.start({}, FlutterClashConfig()));
  });
}
