package main

import "os"

func (gpkg *Gpkg) sources() {
	if len(os.Args) < 2 {
		for _, src := range gpkg.gvm.sources {
			gpkg.logger.Info(src.root)
		}
		return
	}

	command := readCommand()
	if command == "add" {
		gpkg.gvm.AddSource(os.Args[1])
		gpkg.logger.Message("Added", os.Args[1], "to sources")
	} else if command == "remove" {
		if !gpkg.gvm.RemoveSource(os.Args[1]) {
			gpkg.logger.Error("Couldn't remove", os.Args[1])
		}
		gpkg.logger.Message("Removed", os.Args[1], "from sources")
	} else {
	}
}
