@ECHO off
:: VERSION
FOR /F "tokens=* USEBACKQ" %%F IN (`git describe --tags --always`) DO (SET VERSION=%%F)
:: ECHO VERSION: %VERSION%
:: BUILD_DATE
FOR /f "tokens=2 delims==" %%G in ('wmic os get localdatetime /value') do set datetime=%%G
SET BUILD_DATE=%datetime:~0,4%-%datetime:~4,2%-%datetime:~6,2% %time:~0,8%
:: ECHO BUILD_DATE: %BUILD_DATE%
:: COMMIT_HASH
FOR /F "tokens=* USEBACKQ" %%F IN (`git rev-parse HEAD`) DO (SET COMMIT_HASH=%%F)
:: ECHO COMMIT_HASH: %COMMIT_HASH%
:: BRANCH
FOR /F "tokens=* USEBACKQ" %%F IN (`git branch --show-current`) DO (SET BRANCH=%%F)
:: ECHO BRANCH: %BRANCH%
:: Build
set LDFLAGS="-w -s -X github.com/feimingxliu/quicksearch/pkg/about.Branch=%BRANCH% -X github.com/feimingxliu/quicksearch/pkg/about.Version=%VERSION% -X 'github.com/feimingxliu/quicksearch/pkg/about.BuildDate=%BUILD_DATE%' -X github.com/feimingxliu/quicksearch/pkg/about.CommitHash=%COMMIT_HASH%"
:: ECHO LDFLAGS: %LDFLAGS%
go build -ldflags=%LDFLAGS% -o bin/quicksearch.exe github.com/feimingxliu/quicksearch/cmd/quicksearch