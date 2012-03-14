package gpkg

type Builder struct {
	tmpdir string
}

func NewBuilder() *Builder {
	return &Builder{}
}

//func (builder *Builder) 



//import "path/filepath"
//import "fmt"
//import "strings"
//import "os"
/*
const LOG_LEVEL = DEBUG
type PackageNode struct {
	spec string
	ready bool
	versions[] Version
	Package
	parent *PackageNode
}

type Builder struct {
	Logger
	tmpdir string
	pkg_node *PackageNode
	*Sources
	deps map[string]*PackageNode
}

func NewPackageNode(p *Package, spec string, versions[] Version) *PackageNode {
	pn := &PackageNode{}
	pn.Package = *p
	pn.spec = spec
	pn.versions = versions
	return pn
}

func NewBuilder(pkg_node *PackageNode, tmpdir string) *Builder {
	builder := &Builder{}
	builder.tmpdir = tmpdir
	builder.pkg_node = pkg_node
	builder.Logger = NewLogger("", LOG_LEVEL)
	builder.deps = make(map[string]*PackageNode, 255)
	return builder
}

func (builder *Builder) traceDep(p *PackageNode) {
	var stack[]*PackageNode
	count := 0
	pp := p
	for pp.parent != nil {
		count++
		pp = pp.parent
	}

	stack = make([]*PackageNode, count)
	pos := 0
	pp = p
	for pp.parent != nil {
		stack[count - pos - 1] = pp
		pos++
		pp = pp.parent
	}
	println(pp.Name())
	for n, p := range stack {
		println(strings.Repeat("  ", n+1), p.Name())
	}

	fmt.Println(strings.Repeat("  ", count+1), " *", "version:", p.version.String())
	fmt.Println(strings.Repeat("  ", count+1), " *", "spec:", p.spec)
	fmt.Println(strings.Repeat("  ", count+1), " *", "choices:", p.versions)
}*/
/*
func (builder *Builder) resolveConflict(offender *PackageNode, deffender *PackageNode) bool {
	if deffender.spec == "*" {
		deffender.spec = offender.spec
		deffender.version = offender.version
		deffender.versions = offender.versions
		// Downgrade redo
		builder.Debug("Version challege failed. Retroactive downgrade version of", deffender.Name(), "in", deffender.parent.Name())
		return true
	} else if offender.spec == "*" {
		builder.Debug("Version challege failed. Graceful downgrade version of", offender.Name(), "in", offender.parent.Name())
		return false
	}

	builder.traceDep(deffender)
	builder.traceDep(offender)
	builder.Fatal("Version conflict\n")
	return false
}*/
/*
func (builder *Builder) ResolveVersions(depth int, container Container, p *PackageNode) {
	for dep, spec := range p.Specs.List() {
		//debug_trace := fmt.Sprint(strings.Repeat("  ", depth), p.Name(), " ", p.Version() + "->" + dep, " ", spec)
		debug_trace := fmt.Sprint(strings.Repeat("  ", depth), dep, " ", spec)
		dep_pkg := builder.Sources.FindBySpec(dep, spec)
		if dep_pkg == nil {
			builder.Fatal("Couldn't find " + dep + " " + spec + " in any sources\n")
		}

		dep_pkg.parent = p

		if builder.deps[dep] != nil {
			if dep_pkg.Version() != builder.deps[dep].Version() {
				if !builder.resolveConflict(dep_pkg, builder.deps[dep]) {
					// Dependecies updated rerun resolver with new version
					builder.ResolveVersions(depth + 1, container, builder.pkg_node)
				}
			}
		} else {
			builder.Debug(debug_trace, "found:", len(dep_pkg.versions), "using:", dep_pkg.Version())
			builder.deps[dep] = dep_pkg
		}

		if dep_pkg != nil {
			dep_pkg.Clone(container)
			builder.ResolveVersions(depth + 1, container, dep_pkg)
		}
	}
}*/
/*
func (builder *Builder) BuildPackage(src Container, dest Container, p *PackageNode) {
	out := ""

	builder.Info("Building", p.name + "-" + p.Version())

	p.Clone(src)

	builder.Debug("== Cleaning", p.name)
	good, out := p.Clean()
	builder.Debug(out)
	if good == false {
		builder.Fatal("Error cleaning")
	}
	builder.Debug("== Building", p.name)
	good, out = p.Build()
	builder.Debug(out)
	if good == false {
		builder.Fatal("Error building")
	}
	builder.Debug("== Testing", p.name)
	good, out = p.Test()
	builder.Debug(out)
	if good == false {
		builder.Fatal("Failed test")
	}
	builder.Debug("== Installing", p.name)
	good, out = p.Install()
	builder.Debug(out)
	if good == false {
		builder.Fatal("Error installing")
	}
}

func (builder *Builder) BuildDependencies(depth int, src Container, dest Container, p *PackageNode) {
	for name, _ := range p.Specs.List() {
		dep := builder.deps[name]
		if dep.ready == true {
			continue
		}
		builder.BuildDependencies(depth + 1, src, dest, builder.deps[name])
	}

	if p == builder.pkg_node {
		return
	}

	builder.BuildPackage(src, dest, p)
	p.ready = true
}
*/
/*func (builder *Builder) Build() {
	p := builder.pkg_node
	builder.Info("Installing", p.Name(), p.Version())

	src_container := NewContainerWithRoot(filepath.Join(builder.tmpdir, "src"))
	build_container := NewContainerWithRoot(filepath.Join(builder.tmpdir, "build"))
	deps_container := NewContainerWithRoot(filepath.Join(builder.tmpdir, "deps"))
	defer func() {
		build_container.Close()
		deps_container.Close()
		src_container.Close()
	}()
	p.Clone(src_container)

	builder.Info("== Resolve Versions", p.name)
	//builder.ResolveVersions(1, src_container, p)


	builder.Info("== Build Dependencies", p.name)
	gopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", deps_container.Root())
	builder.BuildDependencies(1, src_container, deps_container, p)


	builder.Info("== Build", p.name, p.Version())
	os.Setenv("GOPATH", build_container.Root() + ":" + deps_container.Root())
	builder.BuildPackage(src_container, build_container, p)
	os.Setenv("GOPATH", gopath)
}
*/
