package main
/*
import "os"
import "reflect"
import "strings"
import "gpkg"
//import "strconv"
import "path/filepath"

type AppInterface interface {
	Find()
}

type App struct {
	AppInterface
	args[] string
	config string
	sources *gpkg.Sources
}

func (app *App) Install() {*/
/*	name := app.readCommand()
	spec := strings.Join(app.args[1:], " ")
	var version *gpkg.Version
	if spec != "" {
		version = app.sources.FindBySpec(name, spec)
	} else {
		version = app.sources.Find(name)
	}
	if version == nil {
		panic("NOOOO!")
	}
	builder := gpkg.NewBuilder(name, version, filepath.Join("/home/jbussdieker/Desktop", strconv.Itoa(os.Getpid())))
	builder.Sources = app.sources
	println(builder)*/
//	builder.Build()

/*	var p *gpkg.PackageNode
	p = nil
	name := app.readCommand()
	spec := strings.Join(app.args[1:], " ")
	if spec != "" {
		p = app.sources.FindBySpec(name, spec)
	} else {
		p = app.sources.Find(name)
	}
	if p != nil {
		builder := gpkg.NewBuilder(p, filepath.Join("/home/jbussdieker/Desktop", strconv.Itoa(os.Getpid())))
		builder.Sources = app.sources
		builder.Build()
	} else {
		println("Package not found")
	}*/
/*}

func (app *App) readCommand() string {
	if len(app.args) > 1 {
		app.args = app.args[1:]
		return app.args[0]
	}
	return ""
}

func NewApp(args[] string, config string) *App {
	app := &App{
		args: os.Args,
		config: config,
	}	
	app.sources = gpkg.NewSources(filepath.Join(config, "sources"))
	return app
}
*/
func main() {
/*	app := NewApp(os.Args, "/home/jbussdieker/.gvm/config")
	command := app.readCommand()

	app_interface := reflect.ValueOf(app)
	v := app_interface.MethodByName(strings.Title(command))
	if v.IsValid() {
		args := make([]reflect.Value, 0) 
		v.Call(args)
	} else {
		println("Invalid command")
		os.Exit(1)
	}*/
}
