
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

    echo "Remove Icon file"
    rm -r public/images/icon

    echo "Remove Upload file"
    rm -r data

    cd ..
    echo "generate statik.go"
    statik.exe -src=.speaks -p config -f
    echo "Success!"

    cd ..

    ;;
esac
