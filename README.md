#gpkg

gpkg is the package manager for http://github.com/moovweb/gvm.

Once you've installed Go using the instructions found at http://github.com/moovweb/gvm you'll have the gpkg command.

##Sources

gpkg uses a list of sources to find packages. You can add and remove source via the gpkg command:

* Add a source
`gpkg sources add github.com/moovweb`
* Remove a source
`gpkg sources remove github.com/badrepo`

One exception is the special package name "." which prompts gpkg to install from a local folder in the current working directory.

##Packages

Creating a gpkg package is pretty simple.

`````
mkdir example1
cd example1
echo "package main
func main() {
  println(\"Hello World\")
}" >> main.go
gpkg install .
example1
``````  