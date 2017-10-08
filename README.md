This is a prototype of a personal VPN/Sentry server (https://www.vpnsentry.net/). This is a reverse engineering effort as the code for VPN/Sentry is not open source. In its current state, the Android app is only able to successfully connect to the server, and not actually do any data transfer.

## Domain setup
* You'll need to make a DNS A record for france.yourdomain.com that points to wherever you want to run the server

## Build custom Android APK that points at france.yourdomain.com
This will allow you to build a custom version of the VPN/Sentry android APK that can point at your own server.
* You'll need Java installed to build this
* Clone this repository and cd into apk directory
* Download an APK version of VPN/Sentry to file and save as vpn_entry.apk
* Decompile APK and copy ServiceTunnel\$CommandHandler.smali file into this directory and replace two instances of 'vpnsentry.net' with your own domain that you control
* Run ./build.sh
* Modified apk is now at vpn_sentry_new.apk

## Run prototype VPN/Sentry server 
This will allow you to build a custom version of the VPN/Sentry android APK that can point at your own server.
* Generate certs for france.yourdomain.com (e.g. https://letsencrypt.org/). I used acme (https://github.com/google/acme) tool to create them.
* Install Go (https://golang.org/)
* go run main.go -c ~/.config/acme/france.yourdomain.com.crt -k ~/.config/acme/france.yourdomain.com.key -v

## Setup mitmproxy
* Install mitmproxy (https://mitmproxy.org/) 
* mitmdump -vw filename -d --raw-tcp --insecure --tcp <france.yourdomain.com>

## Setup Android emulator (optional: makes it quicker to test)
* Install Android Studio (https://developer.android.com/studio/index.html) and setup an emulated device
* adb install vpn_sentry_new.apk
* In emulator, Go to Settings->Proxy->Manual Proxy Configuration and point it at mitmproxy (e.g. ip:8080)

## Testing
* On Android device or emulator you should now be able to start the modified VPN/Sentry app and have it connect to your endpoint. You should see all data flowing through mitmproxy in unencrypted form.