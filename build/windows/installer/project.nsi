!include "MUI2.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"

!define APPNAME "runbit"
!define APPVERSION "1.4.0"

Name "${APPNAME}"
OutFile "..\..\bin\runbit-installer.exe"
InstallDir "$PROGRAMFILES\runbit"
RequestExecutionLevel admin

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Spanish"

Section "Main"
  SetShellVarContext all
  SetOutPath $INSTDIR
  
  File "..\..\bin\runbit.exe"
  
  CreateShortcut "$DESKTOP\runbit.lnk" "$INSTDIR\runbit.exe"
  CreateDirectory "$SMPROGRAMS\runbit"
  CreateShortcut "$SMPROGRAMS\runbit\runbit.lnk" "$INSTDIR\runbit.exe"
  
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayName" "${APPNAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayVersion" "${APPVERSION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "Publisher" "Moisessantos45"
SectionEnd

Section "Uninstall"
  SetShellVarContext all
  Delete "$INSTDIR\*.*"
  RMDir "$INSTDIR"
  Delete "$DESKTOP\runbit.lnk"
  Delete "$SMPROGRAMS\runbit\runbit.lnk"
  RMDir "$SMPROGRAMS\runbit"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}"
SectionEnd
