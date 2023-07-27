@echo off
if "%1"=="swag" goto:swag
set CONFIG_FILE=./configs/%2-config-default.yaml
if "%1"=="run" goto:run
if "%1"=="test" goto:test
goto:eof

:swag
if "%2"=="" goto:error_second_parameter_required
swag init -q -d ./cmd/%2,./internal/%2/handler,./pkg/utils -p snakecase -o ./internal/%2/docs
goto:eof

:run
if "%2"=="" goto:error_second_parameter_required
go build -o auth.exe ./cmd/%2/main.go && %2.exe
goto:eof

:test
if "%2"=="" goto:error_second_parameter_required
go clean -testcache
if "%2"=="all" goto:test_all
go test -v -run=Auth .
goto:eof

:test_all
set CONFIG_FILE=./configs/auth-config-default.yaml
go test -v -run=Auth .

:error_second_parameter_required
echo "error: second parameter required"
goto:eof