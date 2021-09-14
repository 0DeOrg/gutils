module gutils

go 1.16

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-resty/resty/v2 v2.6.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/viper v1.8.1
	go.uber.org/zap v1.17.0
)

replace gutils => ../gutils
