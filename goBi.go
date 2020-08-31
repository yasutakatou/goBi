/*
 * Simple URL's security score check tool.
 *
 * @author    yasutakatou
 * @copyright 2020 yasutakatou
 * @license   GPL-3.0 License
 * @FYI: https://qiita.com/MasatoraAtarashi/items/eec4642fe1e6ce79304d
 *       - 語尾に自動でやんすをつけてくれるアプリ
 */
/*
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	hook "github.com/robotn/gohook"
	"github.com/taglme/string2keyboard"
)

var (
	debug                   bool
	zenkaku                 bool
	user32                  = syscall.MustLoadDLL("user32.dll")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
)

type gobiData struct {
	Mae      []int
	MaeCount int
	Ato      string
}

func main() {
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Zenkaku := flag.Bool("zenkaku", true, "[-zenkaku=zenkaku mode (true is enable)]")
	_Delete := flag.Int("del", 8, "[-del=string delete key]")
	_Split := flag.String("split", "@", "[-split=string for split]")
	flag.Parse()

	debug = bool(*_Debug)
	zenkaku = bool(*_Zenkaku)

	var gobi []gobiData

	for i := 0; i < flag.NArg(); i++ {
		if strings.Index(flag.Arg(i), string(*_Split)) != -1 {
			strs := strings.Split(flag.Arg(i), string(*_Split))
			count, err := strconv.Atoi(strs[0])
			if err == nil {
				if count > 0 && count < 256 {
					gobi = append(gobi, gobiData{Mae: []int{count}, MaeCount: 0, Ato: strs[1]})
				}
			} else {
				gobi = append(gobi, gobiData{Mae: intsConvert(strs[0]), MaeCount: 0, Ato: strs[1]})
			}
		}
	}

	if len(gobi) == 0 {
		fmt.Println("no convert defined!")
		os.Exit(1)
	}

	if debug == true {
		fmt.Println(gobi)
	}

	do(gobi, int(*_Delete))

	os.Exit(0)
}

func intsConvert(strs string) []int {
	var ints []int

	for _, c := range strs {
		ints = append(ints, int(c))
	}
	return ints
}

func do(gobi []gobiData, deleteKey int) {
	EvChan := hook.Start()

	currentHwnd := GetWindow("GetForegroundWindow")

	shiftFlag := 0
	inputCount := 0

	for ev := range EvChan {
		if ev.Kind == 4 || ev.Kind == 5 { //KeyHold = 4,KeyUp = 5
			if inputCount > 0 {
				inputCount = inputCount - 1
			} else {
				shiftFlag, strs := keyHoldUp(int(ev.Rawcode), int(ev.Kind), shiftFlag)
				if shiftFlag == 256 {
					hook.End()
					return
				}
				if len(strs) > 0 {
					nowHwnd := GetWindow("GetForegroundWindow")
					if currentHwnd == nowHwnd {
						result := checkRuleAndGo(gobi, strs)
						if result != 0 {
							delKey(len(gobi[(result - 1)].Mae))
							string2keyboard.KeyboardWrite(gobi[(result - 1)].Ato)
							inputCount = len(gobi[(result - 1)].Mae) + len(gobi[(result - 1)].Ato)
	
							if debug == true {
								fmt.Println("type: " + gobi[(result - 1)].Ato)
							}
							for i := 0; i < len(gobi); i++ {
								gobi[i].MaeCount = 0
							}
						}
					} else {
						currentHwnd = GetWindow("GetForegroundWindow")
						for i := 0; i < len(gobi); i++ {
							gobi[i].MaeCount = 0
						}
					}
				}	
			}
		}

	}
}

func delKey(keys int) {
	if zenkaku == true {
		if (keys / 2) < 1 {
			keys = 1
		} else {
			keys = keys / 2
		}
	}
	for i := 0; i < keys; i++ {
		string2keyboard.KeyboardWrite("\\b")
	}

}

func checkRuleAndGo(gobi []gobiData, strs string) int {
	for i := 0; i < len(gobi); i++ {
		if len(gobi[i].Mae) > gobi[i].MaeCount {
			if gobi[i].Mae[gobi[i].MaeCount] == int(([]rune(strs))[0]) {
				gobi[i].MaeCount = gobi[i].MaeCount + 1
			}

			if len(gobi[i].Mae) == gobi[i].MaeCount {
				return (i + 1)
			}
		}
	}
	if debug == true {
		fmt.Println([]rune(strs))
		fmt.Println(gobi)
	}
	return 0
}

func keyHoldUp(Rawcode, Kind, shiftFlag int) (int, string) {
	switch Rawcode {
	case 160:
		if Kind == 4 {
			return Rawcode, ""
		} else {
			return 0, ""
		}
	case 27: //Default Escape
		return 256, ""
	default:
		if Kind == 5 {
			if shiftFlag == 160 {
				return shiftFlag, string(Rawcode)
			} else {
				return shiftFlag, strings.ToLower(string(Rawcode))
			}
		}
	}
	return shiftFlag, ""
}

func GetWindow(funcName string) uintptr {
	hwnd, _, _ := syscall.Syscall(procGetForegroundWindow.Addr(), 6, 0, 0, 0)
	if debug == true {
		fmt.Printf("currentWindow: handle=0x%x.\n", hwnd)
	}
	return hwnd
}
