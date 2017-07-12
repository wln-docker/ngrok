@echo off
:: ===================
:: 内网穿透客户端脚本
:: ===================
cls
set port=80
set host=www
set /p port=请输入要映射的本地端口号（%port%）：
:port
if %port% gtr 99999999999999999999999999999999 (
    cls
    set /p port=请输入有效的端口号（0至65535）：
    goto port
) else if %port% leq 0 (
    set port=Nan
    goto port
) else if %port% geq 65535 (
    set port=Nan
    goto port
)
:input
set /p host=请输入要映射的子域名：
if "%host%"=="www" (
    goto input
)
ngrok -config ngrok.cfg -subdomain=%host% %port%