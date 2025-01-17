package parser

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/therecipe/qt/internal/utils"
)

const (
	SIGNAL = "signal"
	SLOT   = "slot"
	PROP   = "prop"

	IMPURE = "impure"
	PURE   = "pure"

	PLAIN            = "plain"
	CONSTRUCTOR      = "constructor"
	COPY_CONSTRUCTOR = "copy-constructor"
	MOVE_CONSTRUCTOR = "move-constructor"
	DESTRUCTOR       = "destructor"

	CONNECT    = "Connect"
	DISCONNECT = "Disconnect"
	CALLBACK   = "callback"

	GETTER = "getter"
	SETTER = "setter"

	VOID = "void"

	TILDE = "~"

	MOC = "moc"
)

func IsPackedList(v string) bool {
	return (strings.HasPrefix(v, "QList<") ||
		//TODO: QLinkedList
		strings.HasPrefix(v, "QVector<") ||
		strings.HasPrefix(v, "QStack<") ||
		strings.HasPrefix(v, "QQueue<")) &&
		//TODO: QSet

		strings.Count(v, "<") == 1 //TODO:
}

func UnpackedList(v string) string {
	return CleanValue(UnpackedListDirty(v))
}

func UnpackedListDirty(v string) string {
	return strings.Split(strings.Split(v, "<")[1], ">")[0]
}

func IsPackedMap(v string) bool {
	return (strings.HasPrefix(v, "QMap<") ||
		strings.HasPrefix(v, "QMultiMap<") ||
		strings.HasPrefix(v, "QHash<") ||
		strings.HasPrefix(v, "QMultiHash<")) &&

		strings.Count(v, "<") == 1 //TODO:
}

func UnpackedMap(v string) (string, string) {
	var splitted = strings.Split(UnpackedList(v), ",")
	return splitted[0], splitted[1]
}

func UnpackedMapDirty(v string) (string, string) {
	var splitted = strings.Split(UnpackedListDirty(v), ",")
	return splitted[0], splitted[1]
}

func CleanValue(v string) string {
	if IsPackedList(cleanValueUnsafe(v)) || IsPackedMap(cleanValueUnsafe(v)) {
		var inside = strings.Split(strings.Split(v, "<")[1], ">")[0]
		return strings.Replace(cleanValueUnsafe(v), strings.Split(strings.Split(cleanValueUnsafe(v), "<")[1], ">")[0], inside, -1)
	}
	return cleanValueUnsafe(v)
}

func cleanValueUnsafe(v string) string {
	for _, b := range []string{"*", "const", "&amp", "&", ";"} {
		v = strings.Replace(v, b, "", -1)
	}
	return strings.TrimSpace(v)
}

func CleanName(name, value string) string {
	switch name {
	case
		"type",
		"func",
		"range",
		"string",
		"int",
		"map",
		"const",
		"interface",
		"select",
		"strings",
		"new",
		"signal",
		"ptr",
		"register":
		{
			return name[:len(name)-2]
		}

	case "":
		{
			var v = strings.Replace(CleanValue(value), ".", "", -1)
			if len(v) >= 3 {
				return fmt.Sprintf("v%v", strings.ToLower(v[:2]))
			} else {
				return fmt.Sprintf("v%v", strings.ToLower(v))
			}
		}
	}

	return name
}

var LibDeps = map[string][]string{
	"Core":          {"Widgets", "Gui", "Svg"}, //Widgets, Gui //Svg (needed because it's more convenient)
	"AndroidExtras": {"Core"},
	"Gui":           {"Widgets", "Core"}, //Widgets
	"Network":       {"Core"},
	"Xml":           {"XmlPatterns", "Core"}, //XmlPatterns
	"DBus":          {"Core"},
	"Nfc":           {"Core"},
	"Script":        {"Core"},
	"Sensors":       {"Core"},
	"Positioning":   {"Core"},
	"Widgets":       {"Gui", "Core"},
	"Sql":           {"Widgets", "Gui", "Core"}, //Widgets, Gui
	"MacExtras":     {"Gui", "Core"},
	"Qml":           {"QuickControls2", "Network", "Core"},
	"WebSockets":    {"Network", "Core"},
	"XmlPatterns":   {"Network", "Core"},
	"Bluetooth":     {"Core"},
	"WebChannel":    {"Network", "Qml", "Core"}, //Network (needed for static linking ios)
	"Svg":           {"Widgets", "Gui", "Core"},
	"Multimedia":    {"MultimediaWidgets", "Widgets", "Network", "Gui", "Core"},                     //MultimediaWidgets, Widgets
	"Quick":         {"QuickControls2", "QuickWidgets", "Widgets", "Network", "Qml", "Gui", "Core"}, //QuickWidgets, Widgets, Network (needed for static linking ios)
	"Help":          {"Sql", "CLucene", "Network", "Widgets", "Gui", "Core"},                        //Sql + CLucene + Network (needed for static linking ios)
	"Location":      {"Positioning", "Quick", "Gui", "Core"},
	"ScriptTools":   {"Script", "Widgets", "Core"}, //Script, Widgets
	"UiTools":       {"Widgets", "Gui", "Core"},
	"X11Extras":     {"Gui", "Core"},
	"WinExtras":     {"Widgets", "Gui", "Core"},
	"WebEngine":     {"Widgets", "WebEngineWidgets", "WebChannel", "Network", "WebEngineCore", "Quick", "Gui", "Qml", "Core"}, //Widgets, WebEngineWidgets, WebChannel, Network
	"TestLib":       {"Widgets", "Gui", "Core"},                                                                               //Widgets, Gui
	"SerialPort":    {"Core"},
	"SerialBus":     {"Core"},
	"PrintSupport":  {"Widgets", "Gui", "Core"},
	//"PlatformHeaders": []string{}, //TODO: uncomment
	"Designer": {"UiPlugin", "Widgets", "Gui", "Xml", "Core"},
	"Scxml":    {"Network", "Qml", "Core"}, //Network (needed for static linking ios)
	"Gamepad":  {"Gui", "Core"},

	"Purchasing":        {"Core"},
	"DataVisualization": {"Gui", "Core"},
	"Charts":            {"Widgets", "Gui", "Core"},
	//"Quick2DRenderer":   {}, //TODO: uncomment

	//"NetworkAuth":    {"Network", "Gui", "Core"},
	"Speech":         {"Core"},
	"QuickControls2": {"Quick", "QuickWidgets", "Widgets", "Network", "Qml", "Gui", "Core"}, //Quick, QuickWidgets, Widgets, Network, Qml, Gui (needed for static linking ios)

	"Sailfish": {"Core"},
	"WebView":  {"Core"},

	MOC:         make([]string, 0),
	"build_ios": {"Core", "Gui", "Network", "Sql", "Xml", "Nfc", "Script", "Sensors", "Positioning", "Widgets", "Qml", "WebSockets", "XmlPatterns", "Bluetooth", "WebChannel", "Svg", "Multimedia", "Quick", "Help", "Location", "ScriptTools", "MultimediaWidgets", "UiTools", "PrintSupport", "WebView"},
}

var Libs = []string{
	"Core",
	"AndroidExtras",
	"Gui",
	"Network",
	"Xml",
	"DBus",
	"Nfc",
	"Script", //depreached (planned) in 5.6
	"Sensors",
	"Positioning",
	"Widgets",
	"Sql",
	"MacExtras",
	"Qml",
	"WebSockets",
	"XmlPatterns",
	"Bluetooth",
	"WebChannel",
	"Svg",
	"Multimedia",
	"Quick",
	"Help",
	"Location",
	"ScriptTools", //depreached (planned) in 5.6
	"UiTools",
	"X11Extras",
	"WinExtras",
	"WebEngine",
	"TestLib",
	"SerialPort",
	"SerialBus",
	"PrintSupport",
	//"PlatformHeaders", //missing imports/guards
	"Designer",
	"Scxml",
	"Gamepad",

	"Purchasing",        //GPLv3 & LGPLv3
	"DataVisualization", //GPLv3
	"Charts",            //GPLv3
	//"Quick2DRenderer",   //GPLv3

	//"NetworkAuth",
	"Speech",
	"QuickControls2",

	"Sailfish",
	"WebView",
}

func ShouldBuild(module string) bool {
	return true
}

func GetLibs() []string {
	for i := len(Libs) - 1; i >= 0; i-- {
		switch {
		case !(runtime.GOOS == "darwin" || runtime.GOOS == "linux") && (Libs[i] == "WebEngine" || Libs[i] == "WebView"),
			runtime.GOOS != "windows" && Libs[i] == "WinExtras",
			runtime.GOOS != "darwin" && Libs[i] == "MacExtras",
			runtime.GOOS != "linux" && Libs[i] == "X11Extras":
			{
				Libs = append(Libs[:i], Libs[i+1:]...)
			}

		case utils.QT_VERSION() != "5.8.0" && Libs[i] == "Speech":
			{
				Libs = append(Libs[:i], Libs[i+1:]...)
			}
		}
	}
	return Libs
}

func Dump() {
	for _, c := range State.ClassMap {
		var bb = new(bytes.Buffer)
		defer bb.Reset()

		fmt.Fprint(bb, "funcs\n\n")
		for _, f := range c.Functions {
			fmt.Fprintln(bb, f)
		}

		fmt.Fprint(bb, "\n\nenums\n\n")
		for _, e := range c.Enums {
			fmt.Fprintln(bb, e)
		}

		utils.MkdirAll(utils.GoQtPkgPath("internal", "binding", "dump", c.Module))
		utils.SaveBytes(utils.GoQtPkgPath("internal", "binding", "dump", c.Module, fmt.Sprintf("%v.txt", c.Name)), bb.Bytes())
	}
}

func SortedClassNamesForModule(module string, template bool) []string {
	var output = make([]string, 0)
	for _, class := range State.ClassMap {
		if class.Module == module {
			output = append(output, class.Name)
		}
	}
	sort.Stable(sort.StringSlice(output))

	if State.Moc && template {
		var items = make(map[string]string)

		for _, cn := range output {
			if class, ok := State.ClassMap[cn]; ok {
				items[cn] = class.Bases
			}
		}

		var tmpOutput = make([]string, 0)

		for len(items) > 0 {
			for item, dep := range items {

				var c, ok = State.ClassMap[dep]
				if ok && c.Module != MOC {
					tmpOutput = append(tmpOutput, item)
					delete(items, item)
					continue
				}

				for _, key := range tmpOutput {
					if key == dep {
						tmpOutput = append(tmpOutput, item)
						delete(items, item)
						break
					}
				}

			}
		}
		output = tmpOutput
	}

	return output
}

func SortedClassesForModule(module string, template bool) []*Class {
	var (
		classNames = SortedClassNamesForModule(module, template)
		output     = make([]*Class, len(classNames))
	)
	for i, name := range classNames {
		output[i] = State.ClassMap[name]
	}
	return output
}
