// gpkg is the package manager for gvm
//
package main

import "os"

func main() {
	app := NewApp(os.Args)
	if app == nil {
		os.Exit(1)
	}

	switch app.command {
		case "install":
			app.install()
			break
		case "debug":
			//logger.Info(gpkg.gvm.FindPackage("manhattan"))
			return
			break
		case "uninstall":
			app.uninstall()
			break
		case "empty":
			//os.RemoveAll(filepath.Join(app.PkgsetRoot(), "pkg.gvm"))
			break
		case "build":
			app.build()
			break
		case "list":
			app.list()
			break
		case "sources":
			app.sources()
			break
		case "version":
			app.Info(VERSION)
			app.Debug("DEBUG!")
			break
		case "help":
			app.Message("The following commands are available:")
			app.printUsage()
			break
		default:
			app.Error("Invalid command.")
			app.printUsage()
			break
	}
}

