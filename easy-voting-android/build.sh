fyne package -os android -appID com.example.easy-voting -icon ../Icon.png
tar xJvf easy-voting-android.tar.xz
mv usr/local easy-voting-android
echo ./bin/easy-voting-android >> easy-voting-android/launch.sh
chmod +711 easy-voting-android/launch.sh

tar cJvf easy-voting-android.tar.xz easy-voting-android
rm -r usr easy-voting-android Makefile