#! /bin/bash
set -e

function read_dir(){

    if [ "$#" -lt 2 ]; then
        echo "wrong paramters count, exit"
        exit
    fi
    
    out_path=$2

    for file in `ls $1` ;do #注意此处这是两个反引号，表示运行系统命令
        
        if [ -d $1/$file ] ;then #注意此处之间一定要加上空格，否则会报错
            
            read_dir $1/$file $2
        else
            
            file_list=("CoreFoundation" "CoreAutoLayout" "CoreGraphics" "Foundation" "FrontBoardServices" "GraphicsServices" "PhysicsKit" "QuartzCore" "UIKit" "UIKitCore" "libdispatch.dylib" "libdyld.dylib" "libobjc.A.dylib" "libc++abi.dylib" "libsystem_platform.dylib" "libsystem_pthread.dylib" "libsystem_c.dylib" "libsystem_kernel.dylib" "libsystem_malloc.dylib")
        
            if [ $(contains "${file_list[@]}" $file) == "y" ]; then
            
#                mv $1/$file $1/${new_name}
                o_path=$1/$file
                out_path=$2/$2_$file
                cp -rf ${o_path} ${out_path}

                echo "--拷贝文件--"${out_path}
   
            fi
     
        fi
    done
}

function contains() {
    local n=$#
    local value=${!n}
    for ((i=1;i < $#;i++)) {
        if [ "${!i}" == "${value}" ]; then
            echo "y"
            return 0
        fi
    }
    echo "n"
    return 1
}

ipsw="ipsw"
zipDir="zip"
current_path=`pwd`

ipswPath=""$current_path"/"$ipsw"/"
zipDirpath=""$current_path"/"$zipDir"/"

if [ -d $ipswPath ];then
    rm -r $ipswPath
    echo "rm -- "$ipswPath""
fi

if [ -d $zipDirpath ];then
    rm -r $zipDirpath
    echo "rm -- "$zipDirpath""

fi

mkdir $ipsw
mkdir $zipDir


files=$(find . -name "*.ipsw") 
for fnanme in $files 
do 
	s1=${fnanme#*_}
	sysVersion=${s1%%_*}
	echo "fetch sysVersion ---- " $sysVersion

	ipsw dyld extract $fnanme $ipsw 

	cd $ipswPath

	cache_arm_file=$(find . -name "dyld_*") 

	first_path=${cache_arm_file%/*}
	arm=${cache_arm_file##*_}

	version_arm="${sysVersion}_${arm}"
	mkdir $version_arm

	ipsw dyld split $cache_arm_file

	read_dir $first_path $version_arm

	zip_name="${version_arm}.zip"
	# path="${ipswPath}/${version_arm}"
	echo "version_arm --- "$version_arm
	# source_path="${}"
	zip -r $zip_name $version_arm 

	mv $zip_name $zipDirpath

	cd $ipswPath
	rm -rf *

	cd .. 
done

echo "提取完成，请查看zip文件夹"


