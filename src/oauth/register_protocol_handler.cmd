@chcp 866
%SystemRoot%\System32\reg.exe add "HKEY_CLASSES_ROOT\backuper" /ve /d "URL: Backuper Protocol" /f
%SystemRoot%\System32\reg.exe add "HKEY_CLASSES_ROOT\backuper" /v "URL Protocol" /f
%SystemRoot%\System32\reg.exe add "HKEY_CLASSES_ROOT\backuper\shell" /ve /f
%SystemRoot%\System32\reg.exe add "HKEY_CLASSES_ROOT\backuper\shell\open" /ve /f
%SystemRoot%\System32\reg.exe add "HKEY_CLASSES_ROOT\backuper\shell\open\command" /ve /d "\"%~dp0oauth.exe\" \"%~dp0" %%1" /f