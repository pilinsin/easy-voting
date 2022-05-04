fyne package -os linux -icon ../Icon.png
tar xJvf easy-voting-linux.tar.xz
mv usr/local easy-voting-linux
echo ./bin/easy-voting-linux >> easy-voting-linux/launch.sh
chmod +711 easy-voting-linux/launch.sh

tar cJvf easy-voting-linux.tar.xz easy-voting-linux
rm -r usr easy-voting-linux Makefile
