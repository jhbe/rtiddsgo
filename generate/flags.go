package generate

const flags = `// #cgo CFLAGS: -DRTI_UNIX -DRTI_LINUX -DRTI_64BIT -m64 -I{{.RtiInstallDir}}/include -I{{.RtiInstallDir}}/include/ndds -I/usr/include/x86_64-linux-gnu
// #cgo LDFLAGS: -L{{.RtiLibDir}} -lnddsczd -lnddscorezd -ldl -lnsl -lm -lpthread -lrt -Wl,--no-as-needed`
