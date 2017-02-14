#!/bin/sh

curl https://wikileaks.org/dnc-emails/get/1
curl https://wikileaks.org/hbgary-emails/get/1
curl https://wikileaks.org/podesta-emails/get/1

#https://wikileaks.org/podesta-emails/get/regexp
#https://wikileaks.org/podesta-emails/get/%d

go run main.go -target https://wikileaks.org/podesta-emails/get/40890
go run main.go -target https://wikileaks.org/podesta-emails/get/[40890]
go run main.go -target https://wikileaks.org/podesta-emails/get/[:]
go run main.go -target https://wikileaks.org/podesta-emails/get/[5:]
go run main.go -target https://wikileaks.org/podesta-emails/get/[:10]
go run main.go -target https://wikileaks.org/podesta-emails/get/[5:10]

cat config.json | go run main.go
