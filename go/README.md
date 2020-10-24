# go_flutter_clash

This Go package implements the host-side of the Flutter [go_flutter_clash](https://github.com/fanlide/go-flutter-clash) plugin.

## Usage

Import as:

```go
import go_flutter_clash "github.com/fanlide/go-flutter-clash/go"
```

Then add the following option to your go-flutter [application options](https://github.com/go-flutter-desktop/go-flutter/wiki/Plugin-info):

```go
flutter.AddPlugin(&go_flutter_clash.GoFlutterClashPlugin{}),
```
