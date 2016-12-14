package logs

import (
	"fmt"

	seelog "github.com/cihub/seelog"
)

var Logger seelog.LoggerInterface

func loadAppConfig() {
	appConfig := `<seelog >
    <outputs formatid="main">
        <filter levels="info">    
            <console />    
        </filter>
        <filter levels="info">
            <rollingfile formatid="info" type="size" filename="../log/refresh_roll.log" maxsize="100000" maxrolls="5" />
        </filter>
        <filter levels="critical,error">
            <file formatid="critical" path="../log/refresh_critical.log"/>
        </filter>
    </outputs>
    <formats>
        <format id="main" format="%Date/%Time [%LEV] %Msg%n"/>
        <format id="info" format="%Line %Date/%Time [%LEV] %Msg%n"/>
        <format id="critical" format="Date/%Time [%LEV] %Func %Msg %n"/>
    </formats>
	</seelog>`

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	UseLogger(logger)
}

func init() {
	DisableLog()
	loadAppConfig()
}

// DisableLog disables all library log output
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}
