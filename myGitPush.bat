@echo off
REM Git Push Automation Script for Windows

:: Prompt for project directory
set /p projectDir=Bitte Pfad zum Projektverzeichnis eingeben (z.B. C:\labbi-app): 

:: Change to project directory
cd /d "%projectDir%" || (
    echo Fehler: Konnte Verzeichnis %projectDir% nicht finden.
    pause
    exit /b 1
)

echo.
echo Aktuelles Git-Repository: %cd%
git status

echo.
git add .

echo.
:: Prompt for commit message
set /p commitMsg=Bitte Commit-Nachricht eingeben: 

if "%commitMsg%"=="" (
    echo Fehler: Commit-Nachricht darf nicht leer sein.
    pause
    exit /b 1
)

git commit -m "%commitMsg%"

echo.
git pull --rebase

echo.
git push origin main

echo.
echo Fertig! Druecke eine beliebige Taste, um zu beenden.
pause
