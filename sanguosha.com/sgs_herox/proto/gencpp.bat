@ECHO OFF
copy /y "%cd%\def\*.proto" "%cd%\..\AIServer\proto\"
copy /y "%cd%\cmsg\*.proto" "%cd%\..\AIServer\proto\"
copy /y "%cd%\gameproto\*.proto" "%cd%\..\AIServer\proto\"
copy /y "%cd%\gameconf\*.proto" "%cd%\..\AIServer\proto\"

ECHO.
ECHO Finish copy for AI .go Done!
@IF %ERRORLEVEL% NEQ 0 PAUSE