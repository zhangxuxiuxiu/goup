#!/usr/bin/env bash

type=$1
constructor=$2
file=$3

#varname="${type,}"
varname=$(tr '[A-Z]' '[a-z]' <<<${type:0:1})${type:1}

cat >> $file <<EOF

var $varname $type

func Ref${type}() ${type}{
	return $varname
}

func init(){
	$varname = $constructor
}

EOF


