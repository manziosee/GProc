@echo off
echo ========================================
echo GProc - Go Process Manager Demo
echo ========================================
echo.

echo 1. Starting a long-running process...
gproc.exe start demo-app d:\project\GProc\test-app.exe --auto-restart --max-restarts 3
echo.

echo 2. Listing all processes...
gproc.exe list
echo.

echo 3. Waiting 5 seconds...
timeout /t 5 /nobreak > nul
echo.

echo 4. Viewing last 10 lines of logs...
gproc.exe logs demo-app --lines 10
echo.

echo 5. Stopping the process...
gproc.exe stop demo-app
echo.

echo 6. Final process list...
gproc.exe list
echo.

echo Demo completed!