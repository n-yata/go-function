### 参考サイト

[Go の Lambda 関数ハンドラーの定義](https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-handler.html)

[.zip ファイルアーカイブを使用して Go Lambda 関数をデプロイする](https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-package.html)

[Lambda で Go をデプロイするまで　備忘録 (windows)](https://zenn.dev/ru/scraps/666673ea46bde1)

### コマンド（Powershell）
go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest  

$env:GOOS = "linux"  
go env GOOS CGO_ENABLED GOARCH  
set GOOS=linux  
set GOARCH=amd64  
set CGO_ENABLED=0  
go build -tags "lambda.norpc,dev" -o bootstrap .  
build-lambda-zip -o myFunction.zip bootstrap  

### binファイルの場所
C:\Users\{username}\go\bin\  

### Secrets Managerのライブラリ
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
