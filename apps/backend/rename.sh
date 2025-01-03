#!/usr/bin/env sh

export CUR="github.com/suttapak/starter"
export NEW="github.com/suttapak/starter"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;