package report

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	fatihcolor "github.com/fatih/color"
	"github.com/mishamyrt/hapm/internal/manager"
	hapmpkg "github.com/mishamyrt/hapm/internal/package"
	"github.com/mishamyrt/hapm/internal/repository"
)

const tokenGenerateLink = "https://github.com/settings/tokens"

type Reporter struct {
	out io.Writer
}

func New(out io.Writer) Reporter {
	if out == nil {
		out = os.Stdout
	}
	return Reporter{out: out}
}

func (r Reporter) Out() io.Writer {
	return r.out
}

func (r Reporter) NoToken(env string) {
	message := "$" + env + " is not defined.\n"
	message += "Open " + tokenGenerateLink + ",\n"
	message += "generate a personal token and set it in the $" + env + " variable.\n"
	message += "Otherwise you will run into rate limit fairly quickly."
	r.Warning(message)
}

func (r Reporter) WrongFormat(location string) {
	r.Error("Wrong location format: '" + location + "'")
	example := "Package Location can be specified in several formats.\n"
	example += "* Root or tag URL of a repository on GitHub.\n"
	example += "  - https://github.com/mishamyrt/myrt_desk_hass\n"
	example += "  - https://github.com/mishamyrt/myrt_desk_hass/releases/tag/v0.2.4\n"
	example += "* Package name with version separated by the @ symbol\n"
	example += "  - mishamyrt/myrt_desk_hass\n"
	example += "  - mishamyrt/myrt_desk_hass@v0.2.4\n"
	example += "If no version is specified, then latest will be used."
	r.Warning(example)
}

func (r Reporter) Latest(packages []string) {
	message := "No versions are listed for some packages."
	message += "\nThe latest available version will be retrieved and used."
	for _, pkg := range packages {
		message += "\n  " + pkg
	}
	r.Warning(message)
}

func (r Reporter) Exception(action string, err error) {
	r.Error("Error while " + action + ": " + err.Error())
}

func (r Reporter) Warning(text string) {
	_, _ = fmt.Fprintln(r.out, paint(text, fatihcolor.FgYellow))
}

func (r Reporter) Error(text string) {
	_, _ = fmt.Fprintln(r.out, paint(text, fatihcolor.FgRed))
}

func (r Reporter) Diff(diff []manager.PackageDiff, fullName bool, updatesOnly bool) {
	groups := groupDiffByKind(diff)
	keys := make([]string, 0, len(groups))
	for kind := range groups {
		keys = append(keys, kind)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for _, kind := range keys {
		builder.WriteString(formatKind(kind) + "\n")
		for _, pkg := range groups[kind] {
			if updatesOnly {
				builder.WriteString("  " + formatUpdate(pkg) + "\n")
			} else {
				builder.WriteString(formatEntry(pkg, fullName) + "\n")
			}
		}
	}
	_, _ = fmt.Fprint(r.out, builder.String()+"\r")
}

func (r Reporter) Packages(packages []hapmpkg.PackageDescription) {
	groups := groupPackagesByKind(packages)
	keys := make([]string, 0, len(groups))
	for kind := range groups {
		keys = append(keys, kind)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for _, kind := range keys {
		builder.WriteString(formatKind(kind) + "\n")
		for _, pkg := range groups[kind] {
			builder.WriteString(formatPackage(pkg) + "\n")
		}
	}
	_, _ = fmt.Fprint(r.out, builder.String()+"\r")
}

func (r Reporter) Versions(packageName string, versions []string) {
	var builder strings.Builder
	for _, version := range versions {
		builder.WriteString(formatVersion(packageName, version) + "\n")
	}
	_, _ = fmt.Fprint(r.out, builder.String())
}

func (r Reporter) Summary(diff []manager.PackageDiff) {
	if len(diff) == 0 {
		_, _ = fmt.Fprintln(r.out, "There's nothing to do here")
		return
	}
	adds := 0
	deletes := 0
	switches := 0
	for _, pkg := range diff {
		switch pkg.Operation {
		case "add":
			adds++
		case "delete":
			deletes++
		case "switch":
			switches++
		}
	}
	parts := make([]string, 0)
	if adds > 0 {
		parts = append(parts, fmt.Sprintf("installed %s", paint(adds, fatihcolor.FgHiCyan)))
	}
	if deletes > 0 {
		parts = append(parts, fmt.Sprintf("removed %s", paint(deletes, fatihcolor.FgHiCyan)))
	}
	if switches > 0 {
		parts = append(parts, fmt.Sprintf("switched %s", paint(switches, fatihcolor.FgHiCyan)))
	}
	_, _ = fmt.Fprintf(r.out, "\nDone: %s\n", strings.Join(parts, ", "))
}

type Progress struct {
	out      io.Writer
	title    string
	running  bool
	finished bool
	disabled bool
	mu       sync.Mutex
}

var progressSteps = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

func NewProgress(out io.Writer) *Progress {
	if out == nil {
		out = os.Stdout
	}
	return &Progress{out: out, disabled: os.Getenv("HAPM_DISABLE_PROGRESS") == "1"}
}

func (p *Progress) Start(title string) {
	if p.disabled {
		return
	}
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return
	}
	p.title = title
	p.running = true
	p.finished = false
	p.mu.Unlock()
	go p.show()
}

func (p *Progress) Stop() {
	if p.disabled {
		return
	}
	p.mu.Lock()
	p.running = false
	p.mu.Unlock()
	for {
		p.mu.Lock()
		finished := p.finished
		p.mu.Unlock()
		if finished {
			return
		}
		time.Sleep(70 * time.Millisecond)
	}
}

func (p *Progress) show() {
	states := [][2]int{{0, 1}, {2, 1}, {4, 1}}
	prefix := "* " + p.title + " "
	for {
		p.mu.Lock()
		running := p.running
		p.mu.Unlock()
		if !running {
			blank := strings.Repeat(" ", len(prefix)+len(progressSteps))
			_, _ = fmt.Fprint(p.out, blank+"\r")
			p.mu.Lock()
			p.finished = true
			p.mu.Unlock()
			return
		}
		line := prefix
		for i := range states {
			line += progressSteps[states[i][0]]
			if states[i][1] == 1 {
				if states[i][0] == len(progressSteps)-1 {
					states[i][0]--
					states[i][1] = 0
				} else {
					states[i][0]++
				}
			} else {
				states[i][0]--
				if states[i][0] == 0 {
					states[i][1] = 1
				}
			}
		}
		_, _ = fmt.Fprint(p.out, line+"\r")
		time.Sleep(70 * time.Millisecond)
	}
}

func formatKind(kind string) string {
	if kind == "" {
		return paint(":", fatihcolor.Faint)
	}
	label := strings.ToUpper(kind[:1]) + kind[1:] + ":"
	return paint(label, fatihcolor.Faint)
}

func formatUpdate(diff manager.PackageDiff) string {
	nextVersion := paint(diff.Version, fatihcolor.FgYellow)
	currentVersion := paint("("+diff.CurrentVersion+")", fatihcolor.Faint)
	delimiter := paint("@", fatihcolor.Faint)
	return diff.FullName + delimiter + nextVersion + " " + currentVersion
}

func formatEntry(diff manager.PackageDiff, fullName bool) string {
	var textColor fatihcolor.Attribute
	versionStr := diff.Version
	prefix := ""
	switch diff.Operation {
	case "add":
		prefix = "+"
		textColor = fatihcolor.FgGreen
		versionStr = paint(versionStr, fatihcolor.Faint)
	case "switch":
		prefix = "*"
		textColor = fatihcolor.FgYellow
		versionStr = paint(diff.CurrentVersion, fatihcolor.Faint) + " → " + diff.Version
	default:
		prefix = "-"
		textColor = fatihcolor.FgRed
		versionStr = paint(versionStr, fatihcolor.Faint)
	}
	name := repository.RepoName(diff.FullName)
	if fullName {
		name = diff.FullName
	}
	title := paint(prefix+" "+name, textColor)
	return title + "@" + versionStr
}

func formatPackage(pkg hapmpkg.PackageDescription) string {
	version := paint("@"+pkg.Version, fatihcolor.Faint)
	return "  " + pkg.FullName + version
}

func formatVersion(pkg string, version string) string {
	line := paint("- "+pkg+"@", fatihcolor.Faint)
	line += paint(version, fatihcolor.FgYellow)
	return line
}

func paint(value any, attrs ...fatihcolor.Attribute) string {
	if len(attrs) == 0 {
		return fmt.Sprint(value)
	}
	return fatihcolor.New(attrs...).Sprint(value)
}

func groupDiffByKind(packages []manager.PackageDiff) map[string][]manager.PackageDiff {
	groups := map[string][]manager.PackageDiff{}
	for _, pkg := range packages {
		groups[pkg.Kind] = append(groups[pkg.Kind], pkg)
	}
	return groups
}

func groupPackagesByKind(packages []hapmpkg.PackageDescription) map[string][]hapmpkg.PackageDescription {
	groups := map[string][]hapmpkg.PackageDescription{}
	for _, pkg := range packages {
		groups[pkg.Kind] = append(groups[pkg.Kind], pkg)
	}
	return groups
}
