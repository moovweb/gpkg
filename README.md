#gpkg

gpkg is the package manager for http://github.com/moovweb/gvm.

Once you've installed Go using the instructions found at http://github.com/moovweb/gvm you'll have the gpkg command.

##Creating your first package

Creating a gpkg package is pretty simple.

`````
mkdir example1
cd example1
echo "package main
func main() {
  println(\"Hello World\")
}" >> main.go
gpkg build example1
example1
``````  

##Creating and using a custom library

gpkg uses a special Package.gvm file to make imports available during compile time. See the following example:

`````
mkdir lib1
cd lib1
echo "package lib1
func Hello(name string) {
  println(\"Hello\", name)
}" >> lib1.go
gpkg build lib1
cd ..

mkdir example2
cd example2
echo "pkg lib1" >> Package.gvm
echo "package main
import \"lib1\"
func main() {
  lib1.Hello(\"Josh\")
}" >> main.go
gpkg build example2
example2
`````

##Sources

gpkg uses a list of sources to find packages for `gpkg install`. You can add and remove source via the gpkg command:

* Add a source
`gpkg sources add github.com/moovweb`
* Remove a source
`gpkg sources remove github.com/badrepo`