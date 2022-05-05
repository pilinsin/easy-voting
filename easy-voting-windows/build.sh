fyne package -os windows -icon ../Icon.png
tar xJvf easy-voting-windows.tar.xz
mv usr/local easy-voting-windows
echo ./bin/easy-voting-windows >> easy-voting-windows/launch.sh
chmod +711 easy-voting-windows/launch.sh

tar cJvf easy-voting-windows.tar.xz easy-voting-windows
rm -r usr easy-voting-windows Makefile
