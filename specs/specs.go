package specs

import "io/ioutil"
import "strings"
import "strconv"

type SpecError struct {
	msg  string
	line int
}

func (e *SpecError) String() string                { return "Spec Error: " + e.msg + " line " + strconv.Itoa(e.line) }
func NewSpecError(msg string, line int) *SpecError { return &SpecError{msg: msg, line: line} }

type Specs struct {
	Source string
	Origin string
	List   map[string]string
}

func NewBlankSpecs(source string) *Specs {
	return &Specs{Source: source}
}

func NewSpecs(pkgfile string) (*Specs, *SpecError) {
	specs := &Specs{}

	specs.List = map[string]string{}

	data, err := ioutil.ReadFile(pkgfile)
	if err != nil {
		return specs, NewSpecError(err.String(), 0)
	}

	for n, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		switch fields[0] {
		case "pkg":
			if len(fields) > 1 {
				pkg := fields[1]
				criteria := ""
				if len(fields) > 2 {
					criteria = strings.Join(fields[2:], " ")
				}
				specs.List[pkg] = criteria
			} else {
				return specs, NewSpecError("Invalid pkg line in "+pkgfile, n+1)
			}
			break
		case ":source":
			if len(fields) > 1 {
				specs.Source = fields[1]
			} else {
				return specs, NewSpecError("Invalid source line in "+pkgfile, n+1)
			}
			break
		case ":origin":
			if len(fields) > 1 {
				specs.Origin = fields[1]
			} else {
				return specs, NewSpecError("Invalid source line in "+pkgfile, n+1)
			}
			break
		default:
			break
		}
	}

	return specs, nil
}

func (specs *Specs) String() (out string) {
	if specs.Source != "" {
		out += ":source " + specs.Source + "\n"
	}
	if specs.Origin != "" {
		out += ":origin " + specs.Origin + "\n"
	}
	for name, spec := range specs.List {
		out += strings.TrimSpace("pkg "+name+" "+spec) + "\n"
	}
	return
}
