@echo off
@rem  打包脚本

copy  /y server.exe  .\release



@rem

@rem set zdir=".\tool\7z"
@rem %zdir%\7za.exe a  develop.zip server.exe


@rem 7z a archive1.zip subdir\ ：增加subdir文件夹下的所有的文件和子文件夹到archive1.zip中，archived1.zip中的文件名包含subdir\前缀。

@rem 7z a archive2.zip .\subdir\* ：增加subdir文件夹下的所有的文件和子文件夹到archive1.zip中，archived2.zip中的文件名不包含subdir\前缀。


@rem 7z a c:\archive3.zip dir2\dir3\ ：archiive3.zip中的文件名将包含dir2\dir3\前缀，但是不包含c:\dir1前缀。

@rem 7z a Files.7z *.txt -r ： 增加当前文件夹及其子文件夹下的所有的txt文件到Files.7z中。


@rem 7z d archive.zip *.bak -r ：从archive.zip中删除所有的bak文件。

@rem 解压文件到指定的目录
@rem %zdir%\7za.exe e  -y test.zip -o"./testt" -r

pause