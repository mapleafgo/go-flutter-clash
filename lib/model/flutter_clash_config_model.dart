import 'package:json_annotation/json_annotation.dart';

part 'flutter_clash_config_model.g.dart';

@JsonSerializable()
class FlutterClashConfig {
  @JsonKey(name: "port")
  int port;
  @JsonKey(name: "socks-port")
  int socksPort;
  @JsonKey(name: "redir-port")
  int redirPort;
  @JsonKey(name: "tproxy-port")
  int tproxyPort;
  @JsonKey(name: "mixed-port")
  int mixedPort;
  @JsonKey(name: "allow-lan")
  bool allowLan;
  @JsonKey(name: "mode", fromJson: _stringToMode, toJson: _modeToString)
  Mode mode;
  @JsonKey(name: "log-level")
  String logLevel;
  @JsonKey(name: "ipv6")
  bool ipv6;

  FlutterClashConfig({
    this.port,
    this.mode,
    this.allowLan,
    this.tproxyPort,
    this.mixedPort,
    this.socksPort,
    this.redirPort,
    this.logLevel,
    this.ipv6,
  });

  factory FlutterClashConfig.defaultConfig() => FlutterClashConfig(
        port: 0,
        socksPort: 0,
        redirPort: 0,
        tproxyPort: 0,
        mixedPort: 7890,
        allowLan: false,
        mode: Mode.Rule,
        logLevel: "info",
        ipv6: false,
      );

  factory FlutterClashConfig.fromJson(Map<String, dynamic> json) =>
      _$FlutterClashConfigFromJson(json);
  Map<String, dynamic> toJson() => _$FlutterClashConfigToJson(this);
}

Mode _stringToMode(String mode) =>
    ModeMap.entries.firstWhere((t) => t.value == mode).key;

String _modeToString(Mode mode) => ModeMap[mode];

enum Mode { Rule, Global, Direct }
const ModeMap = {
  Mode.Rule: "Rule",
  Mode.Global: "Global",
  Mode.Direct: "Direct",
};
