
echo "This shell is for development."
echo "It deletes the file necessary for the operation from the .speaks directory."
echo "Still run it?[Y/n]"

read answer

case $answer in
  Y)

    echo "Backup .speaks"
    tar cvf - .speaks | gzip - > backup_`date '+%Y%m%d%H%M%S'`.tar.gz

    echo "Move .speaks"
    cd .speaks

    echo "Remove speaks.ini"
    rm speaks.ini

    echo "Remove db file"
    rm *.db

    case $? in
      0)
        echo "generate binary.go"
        go-bindata.exe -pkg config -o ../config/binary.go ./...
        echo "Success!"
      ;;
      1)
        echo "DB remove error"
      ;;
    esac

    cd ..

    ;;
esac
