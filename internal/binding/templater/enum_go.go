package templater

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/therecipe/qt/internal/binding/converter"
	"github.com/therecipe/qt/internal/binding/parser"
)

func goEnum(e *parser.Enum, _ *parser.Value) string {
	var bb = new(bytes.Buffer)
	defer bb.Reset()

	var t string

	fmt.Fprintf(bb, "//go:generate stringer -type=%v\n//%v\ntype %v int64\nconst (\n", strings.Replace(e.Fullname, ":", "_", -1), e.Fullname, strings.Replace(e.Fullname, ":", "_", -1))

	for _, v := range e.Values {
		switch v.Name {
		case "ByteOrder":
			{

			}

		default:
			{
				if strings.Contains(v.Value, " | ") {
					var tArray = make([]string, 0)
					for _, s := range strings.Split(v.Value, " | ") {
						tArray = append(tArray, converter.GoEnum(v.Name, s, e))
					}
					t = strings.Join(tArray, " | ")
				} else {
					t = converter.GoEnum(v.Name, v.Value, e)
				}
				if strings.HasPrefix(t, "C.") && UseStub() {
					t = "0"
				}
				fmt.Fprintf(bb, "%v__%v %v = %v(%v)\n", strings.Split(e.Fullname, "::")[0], v.Name, strings.Replace(e.Fullname, ":", "_", -1), strings.Replace(e.Fullname, ":", "_", -1), t)
			}
		}
	}

	fmt.Fprint(bb, ")")

	if e.NoConst || strings.Contains(e.Name, "Style") {
		return strings.Replace(bb.String(), "const (", "var (", -1)
	}
	return bb.String()
}
