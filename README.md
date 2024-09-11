### 参考サイト

[Go の Lambda 関数ハンドラーの定義](https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-handler.html)

[.zip ファイルアーカイブを使用して Go Lambda 関数をデプロイする](https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-package.html)

[Lambda で Go をデプロイするまで　備忘録 (windows)](https://zenn.dev/ru/scraps/666673ea46bde1)

### コマンド

$env:GOOS = "linux"  
go env GOOS CGO_ENABLED GOARCH  
go get github.com/aws/aws-lambda-go/lambda  
set GOOS=linux  
set GOARCH=amd64  
set CGO_ENABLED=0  
go build -tags lambda.norpc -o bootstrap main.go  
C:\Users\{username}\go\bin\build-lambda-zip.exe -o myFunction.zip bootstrap
