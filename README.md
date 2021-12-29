# atosGo
用Go语言仿atos，实现dSYM符号文件的解析

### atos_dSYM.go
解析dSYM文件符号

### atos_systemSymbol.go
解析系统库符号，如CorFoundation

### 结果对比

#### atosGo结果：
-[NSObject(YYModel) yy_modelToJSONString] (in SweetWords) (NSObject+YYModel.m:1545)
 
#### symbolicatecrash结果：
-[NSObject(YYModel) yy_modelToJSONString] (in SweetWords) (NSObject+YYModel.m:1545)
 
#### atos结果：
-[NSObject(YYModel) yy_modelToJSONString] (in SweetWords) (NSObject+YYModel.m:1545)

#### 注：按需提取系统符号，可以使用这个脚本 parse_ipsw.sh
详情请看文章 https://blog.csdn.net/qq_22326601/article/details/122208364
