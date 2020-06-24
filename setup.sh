RC_FILE="$HOME/.bashrc"

if [ "$#" -eq 1 ]; then
  RC_FILE="$1"
fi

sed -i '/GOPATH=/d' $RC_FILE
sed -i '/GOBIN=/d' $RC_FILE
sed -i '/PATH=$GOBIN:$PATH/d' $RC_FILE

echo "Removed GOPATH, GOBIN definitions"

# setup GOPATH in rc file
GOPATH=$(pwd)
GOBIN="${GOPATH}/bin"
NEW_PATH='$GOBIN:$PATH'
echo "export GOPATH=$GOPATH" >> $RC_FILE 
echo "export GOBIN=$GOBIN" >> $RC_FILE 
echo "export PATH=$NEW_PATH" >> $RC_FILE 

echo "Updated $RC_FILE"

# load RC file
source $RC_FILE
echo "sourced $RC_FILE"

