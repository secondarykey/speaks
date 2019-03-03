
echo "This shell is for development."
echo "It deletes the file necessary for the operation from the .speaks directory."
echo "still run it?[Y/n]"

read answer

case $answer in
  Y)
    tar cvf backup.tar .speaks
    echo "move .speaks"
    cd .speaks
    echo "remove speaks.ini"
    rm speaks.ini
    echo "remove db file"
    rm *.db
    echo "generate binary.go"
    go-bindata.exe -pkg config -o ../config/binary.go ./...
    echo "success."
    cd ..
    ;;
esac
