#!/usr/bin/env bash
set -o nounset
set -o errexit
set -o pipefail

function printUsage {
	echo "Usage: $0 interface pkg.implName srcDir destFile" >&2
	exit 255
}

### parse args
if [[ $# -ne 4 ]]; then 
	printUsage
fi

iface=$1
IFS='.' read -r -a array <<< "$2"
if [[ ${#array[@]} -ne 2 ]]; then
	printUsage 
fi
pkgName=${array[0]}
implName=${array[1]}
srcDir=$3
destFile=$4
if [[ -d $destFile ]]; then
	destFile="${destFile}/${implName}.go"
fi

### generate
echo -e "package $pkgName\n\n type $implName struct{}\n" >"$destFile"

~/gospace/bin/impl -dir "$srcDir"  "_ *$implName" $iface >> $destFile

varName=$(tr '[A-Z]' '[a-z]' <<<${iface:0:1})${iface:1}
cat >>"$destFile" <<EOF
func init(){
	$varName = $implName{}
}
EOF
