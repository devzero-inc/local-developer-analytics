package daemon

const (
	DaemonPlistFilePath    = "Library/LaunchAgents"
	DaemonPlistName        = "lda.plist"
	DaemonServicedFilePath = ".config/systemd/user"
	DaemonServicedName     = "lda.service"
	DaemonPermission       = 0644
)

var (
	DaemonPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>devzer.io.lda</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/zvonimirtomesic/Projects/Codilas/devzero/lda/lda</string>
        <string>collect</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/devzer.io.lda.out.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/devzer.io.lda.err.log</string>
</dict>
</plist>`

	DaemonServiced = `[Unit]
Description=lda
After=network.target

[Service]
Type=simple
ExecStart=/Users/zvonimirtomesic/Projects/Codilas/devzero/lda/lda
Restart=always

[Install]
WantedBy=multi-user.target`
)
