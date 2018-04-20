
::  -------------------------------------------------------------------------------------------------
::
::  Script build.cmd is used for building executable files from sources cmd/backuper and cmd/oauth
::  Executables are stored in ./bin directory alongside with config.json
::
::  If you build Backuper for the first time you need to pass '/i' as parameter (like 'build.cmd /i') 
::  to register oauth.exe as a handler for backuper:// URI protocol. It is necessary to receive
::  API tokens
::
::  -------------------------------------------------------------------------------------------------

@echo off

mkdir bin

go get github.com/fsnotify/fsnotify

cd cmd\backuper
go build
move /y backuper.exe ..\..\bin

cd ..\oauth
go build
move /y oauth.exe ..\..\bin
xcopy /y register_protocol_handler.cmd ..\..\bin

cd ..\..\core\settings
xcopy /y config.json ..\..\bin

cd ..\..\bin

if "%1" == "/i" (
    powershell.exe Start-Process register_protocol_handler.cmd -Verb runAs
)

del /q register_protocol_handler.cmd