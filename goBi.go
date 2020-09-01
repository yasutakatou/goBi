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
	"regexp"
	"strconv"
	"strings"
	"syscall"

	hook "github.com/robotn/gohook"
	"github.com/taglme/string2keyboard"
	"gopkg.in/ini.v1"
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
	_Config := flag.String("config", ".goBi", "[-config=config file)]")
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Zenkaku := flag.Bool("zenkaku", true, "[-zenkaku=zenkaku mode (true is enable)]")
	_Delete := flag.Int("del", 8, "[-del=string delete key]")
	_Split := flag.String("split", "@", "[-split=string for split]")
	flag.Parse()

	debug = bool(*_Debug)
	zenkaku = bool(*_Zenkaku)
	delKey := int(*_Delete)
	Split := string(*_Split)
	Config := string(*_Config)

	var gobi []gobiData
	var bakGobi []string

	if Exists(Config) == true {
		delKey, Split, gobi, bakGobi = loadConfig(Config, Split)
	} else {
		for i := 0; i < flag.NArg(); i++ {
			if strings.Index(flag.Arg(i), Split) != -1 {
				strs := strings.Split(flag.Arg(i), Split)
				count, err := strconv.Atoi(strs[0])
				if err == nil {
					if count > 0 && count < 256 {
						gobi = append(gobi, gobiData{Mae: []int{count}, MaeCount: 0, Ato: strs[1]})
						bakGobi = append(bakGobi, flag.Arg(i))
					}
				} else {
					gobi = append(gobi, gobiData{Mae: intsConvert(strs[0]), MaeCount: 0, Ato: strs[1]})
					bakGobi = append(bakGobi, flag.Arg(i))
				}
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

	do(gobi, delKey)

	saveConfig(Config, delKey, Split, bakGobi, Split)
	os.Exit(0)
}

func saveConfig(filename string, Del int, Spl string, gob []string, Split string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println(zenkaku)
	writeFile(file, "[ZENKAKU]")
	if zenkaku == true {
		writeFile(file, "0")
	} else {
		writeFile(file, "1")
	}

	writeFile(file, "[DELKEY]")
	writeFile(file, strconv.Itoa(Del))

	writeFile(file, "[SPLIT]")
	writeFile(file, Spl)

	writeFile(file, "[GOBI]")
	for i := 0; i < len(gob); i++ {
		writeFile(file, gob[i])
	}
}

func writeFile(file *os.File, strs string) bool {
	_, err := file.WriteString(strs + "\n")
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func loadConfig(filename, Split string) (int, string, []gobiData, []string) {
	var Del int
	var Spl string
	var gob []gobiData
	var bakGob []string

	loadOptions := ini.LoadOptions{}
	loadOptions.UnparseableSections = []string{"ZENKAKU", "DELKEY", "SPLIT", "GOBI"}

	cfg, err := ini.LoadSources(loadOptions, filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	setSingleConfigBool(&zenkaku, "ZENKAKU", cfg.Section("ZENKAKU").Body())
	setSingleConfigInt(&Del, "DELKEY", cfg.Section("DELKEY").Body())
	setSingleConfigStr(&Spl, "SPLIT", cfg.Section("SPLIT").Body())
	setSingleConfigGobi(&gob, &bakGob, "GOBI", cfg.Section("GOBI").Body(), Split)

	return Del, Spl, gob, bakGob
}

func setSingleConfigInt(config *int, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			tmp, err := strconv.Atoi(v)
			if err == nil {
				*config = tmp
			}
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigBool(config *bool, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			if strings.Index(v, "0") != -1 {
				*config = true
				} else {
				*config = false
			}
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigStr(config *string, configType, datas string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			*config = v
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func setSingleConfigGobi(config *[]gobiData, bakConfig *[]string, configType, datas, Split string) {
	if debug == true {
		fmt.Println(" -- " + configType + " --")
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(datas, -1) {
		if len(v) > 0 {
			if strings.Index(v, Split) != -1 {
				strs := strings.Split(v, Split)
				count, err := strconv.Atoi(strs[0])
				if err == nil {
					if count > 0 && count < 256 {
						*config = append(*config, gobiData{Mae: []int{count}, MaeCount: 0, Ato: strs[1]})
						*bakConfig = append(*bakConfig, v)
					}
				} else {
					*config = append(*config, gobiData{Mae: intsConvert(strs[0]), MaeCount: 0, Ato: strs[1]})
					*bakConfig = append(*bakConfig, v)
				}
			}			
		}
		if debug == true {
			fmt.Println(v)
		}
	}
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
							inputCount = len(gobi[(result-1)].Mae) + len(gobi[(result-1)].Ato)

							if debug == true {
								fmt.Println("type: " + gobi[(result-1)].Ato)
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
