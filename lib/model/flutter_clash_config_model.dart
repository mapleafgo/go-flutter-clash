import 'package:json_annotation/json_annotation.dart';

part 'flutter_clash_config_model.g.dart';

@JsonSerializable()
class FlutterClashConfig {
  @JsonKey(name: "port")
  int port = 0;
  @JsonKey(name: "socks-port")
  int socksPort = 0;
  @JsonKey(name: "redir-port")
  int redirPort = 0;
  @JsonKey(name: "mixed-port")
  int mixedPort = 7890;
  @JsonKey(name: "allow-lan")
  bool allowLan = false;
  @JsonKey(name: "mode", fromJson: _stringToMode, toJson: _modeToString)
  Mode mode;
  @JsonKey(name: "log-level")
  String logLevel = "error";
  @JsonKey(name: "ipv6")
  bool ipv6 = false;

  FlutterClashConfig({
    this.port,
    this.mode,
    this.allowLan,
    this.mixedPort,
    this.socksPort,
    this.redirPort,
    this.logLevel,
    this.ipv6,
  });

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
