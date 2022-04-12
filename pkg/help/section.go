package help

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"k8s.io/klog/v2"
)

type Section struct {
	distance *Distance
	name     string
	usetip   string
	contains map[string]interface{}
	childs   map[string]interface{}
	parent   *Section
}

func NewSection(name string, usetip string, distance *Distance) *Section {
	if len(name) < 1 {
		return nil
	}

	var dist Distance
	if distance == nil {
		dist = NewDistance()
	} else {
		dist = *distance
	}

	return &Section{
		distance: &dist,
		name:     name,
		usetip:   usetip,
		contains: nil,
		childs:   nil,
		parent:   nil,
	}
}

func (section *Section) AddTip(name, kind, tips, value string) {
	if len(name) < 1 {
		klog.Errorln("error tip name: \"\"")
		return
	}

	if len(kind) < 1 {
		kind = DefaultValueType
	}

	if section.contains == nil {
		section.contains = make(map[string]interface{})
	}

	section.contains[name] = Tip{
		Name:         name,
		ValueType:    kind,
		TipStr:       tips,
		DefaultValue: value,
	}

	section.distance.UpdateTip(name, kind) // update length
}

func (section *Section) AddSection(newSection *Section) {
	// unuseful section
	if newSection == nil || len(newSection.name) < 1 {
		return
	}

	if section.childs == nil {
		section.childs = make(map[string]interface{})
	}

	newSection.parent = section
	section.childs[newSection.name] = newSection
}

func (section *Section) PrintSectionWithOffset(w io.Writer, offset int, format func(string) string) {
	if format == nil {
		format = FormatNone
	}

	if section == nil {
		return
	}

	offsetString := strings.Repeat(" ", offset)
	fmt.Fprintf(w, "%s%s:\n", offsetString, section.name)
	if len(section.usetip) > 0 {
		fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", offset+section.distance.tipOffset), format(section.usetip))
	}

	if section.contains != nil {
		tipKeys := getKeys(section.contains)
		if len(tipKeys) > 0 {
			fmt.Fprintln(w)
		}
		for _, key := range tipKeys {
			section.contains[key].(Tip).Fprint(w, offset+section.distance.tipOffset, section.distance.maxKeyLength, section.distance.maxValueTypeLength)
		}
	}

	if section.childs != nil {
		childKeys := getKeys(section.childs)
		for _, key := range childKeys {
			if section.childs[key] == nil {
				continue
			}

			fmt.Println()
			section.childs[key].(*Section).PrintSectionWithOffset(w, offset+section.distance.childOffset, format)
		}
	}
}

func (section *Section) PrintSection(w io.Writer, format func(string) string) {
	section.PrintSectionWithOffset(w, 0, format)
}

func FormatHeader(head string) func(string) string {
	return func(raw string) string {
		return head + raw
	}
}

func FormatClamp(head, tail string) func(string) string {
	return func(raw string) string {
		return head + raw + tail
	}
}

func FormatNone(raw string) string {
	return raw
}

func getKeys(m map[string]interface{}) []string {
	if m == nil || len(m) < 1 {
		return []string{}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if DefaultSortPrint {
		sort.Strings(keys)
	}
	return keys
}
