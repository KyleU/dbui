SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"
cd "$DIR"

for arch in amd64 386
do
  for os in darwin linux windows
  do
    echo "Building [$os $arch]..."
    env GOOS=$os GOARCH=$arch make build-release
    mkdir -p ./build/$os/$arch
    if [ "$os" = "windows" ]; then
      mv ./build/release/dbui.exe ./build/$os/$arch/dbui.exe
    else
      echo "$os"
      mv ./build/release/dbui ./build/$os/$arch/dbui
    fi
  done
done
