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
cat config.yaml | go run main.go

cat dump/podesta-emails_101.eml | jq -r .Body | tr '\r\n' '\n'

go run main.go -target https://foia.state.gov/Clinton_Email/[5:10]
#curl 'https://foia.state.gov/searchapp/Search/SubmitSimpleQuery?searchText=*&beginDate=false&endDate=false&collectionMatch=Clinton_Email&postedBeginDate=false&postedEndDate=false&caseNumber=false&page=1&start=5&limit=2'
curl 'https://foia.state.gov/searchapp/Search/SubmitSimpleQuery?searchText=*&beginDate=false&endDate=false&collectionMatch=clinton_email&postedBeginDate=false&postedEndDate=false&caseNumber=false&page=1&start=5&limit=2'
https://foia.state.gov/searchapp/DOCUMENTS/Litigation_F-2016-07895_7/Litigation_F-2016-07895/DOC_0C06134835/C06134835.pdf

go run main.go -target https://foia.state.gov/ELSALVADOR/[5:10]
#curl 'https://foia.state.gov/searchapp/Search/SubmitSimpleQuery?searchText=*&beginDate=false&endDate=false&collectionMatch=ELSALVADOR&postedBeginDate=false&postedEndDate=false&caseNumber=false&page=1&start=5&limit=2'
curl 'https://foia.state.gov/searchapp/Search/SubmitSimpleQuery?searchText=*&beginDate=false&endDate=false&collectionMatch=elsalvador&postedBeginDate=false&postedEndDate=false&caseNumber=false&page=1&start=5&limit=2'
https://foia.state.gov/searchapp/DOCUMENTS\\elsalvad\\79c9.PDF

go run main.go -target wikileaks -type podesta-emails -start 5 -end 10
go run main.go -target foia -type clinton_email -start 5 -end 10
