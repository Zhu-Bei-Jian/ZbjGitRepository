@ECHO OFF
protoc --go_out=. -I=.;../../../   def/*.proto
@IF %ERRORLEVEL% NEQ 0 PAUSE
protoc --go_out=. -I=.;../../../  cmsg/*.proto
@IF %ERRORLEVEL% NEQ 0 PAUSE
protoc --go_out=. -I=.;../../../  db/*.proto
@IF %ERRORLEVEL% NEQ 0 PAUSE
protoc --go_out=. -I=.;../../../  smsg/*.proto
protoc --go_out=. -I=.;../../../  gameproto/*.proto
::protoc --go_out=. -I=.;../../../../  wmsg/*.proto
@IF %ERRORLEVEL% NEQ 0 PAUSE
::protoc --go_out=. -I=.;../../../../  netframe/*.proto
::protoc --go_out=. -I=.;../../../../  logicmsg/*.proto
protoc --go_out=. -I=.;../../../  gameconf/*.proto
@IF %ERRORLEVEL% NEQ 0 PAUSE

ECHO.
ECHO Compile .proto To .go Done!
@IF %ERRORLEVEL% NEQ 0 PAUSE