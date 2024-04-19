package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	xtc "github.com/tomnomnom/xtermcolor"
	"golang.org/x/term"
)

type rgb struct {
	r int
	g int
	b int
}

type Status struct {
	Online   bool   `json:"online"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Hostname string `json:"hostname"`
	Debug    struct {
		Ping          bool `json:"ping"`
		Query         bool `json:"query"`
		Srv           bool `json:"srv"`
		Querymismatch bool `json:"queryismismatch"`
		IPinsrv       bool `json:"ipinsrv"`
		Animatedmotd  bool `json:"animatedmotd"`
		Cachehit      bool `json:"cachehit"`
		Cachetime     int  `json:"cachetime"`
		Cacheexpire   int  `json:"cacheexpire"`
		Apiversion    int  `json:"apiversion"`
	} `json:"debug"`
	Version  string `json:"version"`
	Protocol struct {
		Version int    `json:"version"`
		Name    string `json:"name"`
	} `json:"protocol"`
	Icon     string `json:"icon"`
	Software string `json:"software"`
	Mapname  struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		Html  string `json:"html"`
	} `json:"map"`
	Gamemode     string `json:"gamemode"`
	Serverid     string `json:"serverid"`
	Eula_blocked bool   `json:"eula_blocked"`
	Motd         struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		Html  []string `json:"html"`
	} `json:"motd"`
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
		List   []struct {
			Name string `json:"name"`
			UUID string `json:"uuid"`
		} `json:"list"`
	} `json:"players"`
	Plugins []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"plugins"`
	Mods []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"mods"`
}

const TERMWIDTH = 160
const TERMHEIGHT = 80

func main() {
	var bedrock, silent, exIcon, monitor bool
	var ip, outputPath, loadFile string

	flag.BoolVar(&bedrock, "b", false, "Bedrock Server?")
	flag.BoolVar(&bedrock, "bedrock", false, "Bedrock Server?")
	flag.StringVar(&ip, "ip", "0.0.0.0", "The IP of the server.")
	flag.StringVar(&outputPath, "o", "", "The optional path of the output json file.")
	flag.StringVar(&outputPath, "output", "", "The optional path of the output json file.")
	flag.BoolVar(&silent, "s", false, "Do not output anything to the console.")
	flag.BoolVar(&silent, "silent", false, "Do not output anything to the console.")
	flag.BoolVar(&exIcon, "ei", true, "Do not get icon data.")
	flag.BoolVar(&exIcon, "exIcon", true, "Do not get icon data.")
	flag.BoolVar(&monitor, "m", false, "Whether or not to use the live monitor mode.")
	flag.BoolVar(&monitor, "monitor", false, "Whether or not to use the live monitor mode.")
	flag.StringVar(&loadFile, "lf", "", "Debug mode. Intended for development. Will use pre-saved JSON instead of making an api request.")
	flag.StringVar(&loadFile, "loadFile", "", "Debug mode. Intended for development. Will use pre-saved JSON instead of making an api request.")

	flag.Parse()

	fd := int(os.Stdin.Fd())
	termOriginal, err := term.GetState(fd)
	if err != nil {
		log.Fatal(err)
	}

	var body []byte
	if loadFile == "" {
		var apiURL string
		if bedrock {
			apiURL = "https://api.mcsrvstat.us/bedrock/3/" + ip
		} else {
			apiURL = "https://api.mcsrvstat.us/3/" + ip
		}

		// Get status
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.ReadFile(path.Join(loadFile))
		if err != nil {
			log.Fatal(err)
		}
		body = file
	}
	var status Status
	json.Unmarshal(body, &status)

	if exIcon && !monitor {
		status.Icon = ""
	}

	formattedJSON := formatStatus(status)

	if monitor {
		renderMonitor(status)
	} else {
		if outputPath != "" {
			saveDataToFile(outputPath, formattedJSON)
		}

		if !silent {
			printData(string(formattedJSON))
		}
	}

	term.Restore(fd, termOriginal)
}

func printData(jsonData string) {
	fmt.Println(jsonData)
}

func saveDataToFile(path string, jsonData []byte) {
	os.WriteFile(path, jsonData, 0666)
}

func formatStatus(status Status) []byte {
	ret, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

func parseIconData(iconData string) string {
	splitData := strings.Split(iconData, ",")
	base64Data := splitData[1]
	return base64Data
}

func renderMonitor(status Status) {
	var images []image.Image
	if status.Icon != "" {
		image, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(parseIconData(status.Icon))))
		if err != nil {
			log.Fatalf("failed to decode image: %v", err)
		}
		images = append(images, image)
		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()
	}
	img := widgets.NewImage(nil)
	if len(images) > 0 {
		primColorRGB := getPrimaryImageColor(images[0])
		primColor := xtc.FromColor(color.RGBA{uint8(primColorRGB.r), uint8(primColorRGB.g), uint8(primColorRGB.b), 0})
		img.BorderStyle.Fg = ui.Color(primColor)
	}

	topPadding := int(math.Floor(TERMHEIGHT*(1.0/3)*(1.0/5)/2.0)) - 1

	hostname := widgets.NewParagraph()
	hostname.Text = status.Hostname
	hostname.TextStyle.Fg = ui.ColorCyan
	hostname.BorderStyle.Fg = ui.ColorCyan
	hostname.PaddingTop = topPadding
	hostname.PaddingLeft = 1

	serverOnline := widgets.NewParagraph()
	if status.Online {
		serverOnline.Text = "Online"
		serverOnline.TextStyle.Fg = ui.Color(118)
		serverOnline.BorderStyle.Fg = ui.Color(118)
	} else {
		serverOnline.Text = "Offline"
		serverOnline.TextStyle.Fg = ui.ColorRed
		serverOnline.BorderStyle.Fg = ui.ColorRed
	}
	serverOnline.PaddingTop = topPadding
	serverOnline.PaddingLeft = 1

	onlinePlayers := widgets.NewParagraph()
	onlinePlayers.Text = fmt.Sprintf("Online: %d/%d", status.Players.Online, status.Players.Max)
	onlinePlayers.TextStyle.Fg = ui.ColorCyan
	onlinePlayers.BorderStyle.Fg = ui.ColorCyan
	onlinePlayers.PaddingTop = topPadding
	onlinePlayers.PaddingLeft = 1

	if status.Port == "" {
		status.Port = "25565"
	}
	ipPort := widgets.NewParagraph()
	ipPort.Text = fmt.Sprintf("%s:%s", status.IP, status.Port)
	ipPort.TextStyle.Fg = ui.Color(184)
	ipPort.BorderStyle.Fg = ui.Color(184)
	ipPort.PaddingTop = topPadding
	ipPort.PaddingLeft = 1

	serverVersion := widgets.NewParagraph()
	serverVersion.Text = status.Version
	serverVersion.TextStyle.Fg = ui.Color(105)
	serverVersion.BorderStyle.Fg = ui.Color(105)
	serverVersion.PaddingTop = topPadding
	serverVersion.PaddingLeft = 1

	motd := widgets.NewParagraph()
	motd.Text = strings.Join(status.Motd.Clean, "\n")
	motd.TextStyle.Fg = ui.Color(2)
	motd.BorderStyle.Fg = ui.Color(2)
	motd.PaddingLeft = 1
	motd.PaddingTop = 1

	playerList := widgets.NewList()
	playerList.Rows = getPlayerNamesList(status)
	playerList.WrapText = false
	playerList.Title = onlinePlayers.Text
	playerList.TitleStyle.Fg = ui.ColorCyan
	playerList.BorderStyle.Fg = ui.ColorCyan

	grid := ui.NewGrid()
	grid.SetRect(0, 0, TERMWIDTH, TERMHEIGHT)
	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/3, img),
			ui.NewCol(2.0/3,
				ui.NewRow(1.0/5,
					ui.NewCol(1.0/2, hostname),
					ui.NewCol(1.0/2, serverOnline),
				),
				ui.NewRow(1.0/5, ipPort),
				ui.NewRow(1.0/5, serverVersion),
				ui.NewRow(2.0/5, motd),
			),
		),
		ui.NewRow(2.0/3,
			ui.NewRow(1.0, playerList),
		),
	)

	render := func() {
		if len(images) > 0 {
			img.Image = images[0]
		}
		ui.Render(grid)
	}
	render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
		render()
	}
}

func getTerminalSize() (int, int) {
	if !term.IsTerminal(0) {
		return -1, -1
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return -1, -1
	}
	return width, height
}

func getPlayerNamesList(status Status) []string {
	players := status.Players.List
	var ret []string
	for _, player := range players {
		ret = append(ret, fmt.Sprintf(" - %s", player.Name))
	}
	return ret
}

// TODO: Use K-means clustering instead of average color
func getPrimaryImageColor(img image.Image) rgb {
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	count := 0
	var avg rgb
	avg.r = 0
	avg.g = 0
	avg.b = 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			col := img.At(x, y)
			r, g, b, _ := col.RGBA()
			rFloat := float64(r) / 65535.0 * 255.0
			gFloat := float64(g) / 65535.0 * 255.0
			bFloat := float64(b) / 65535.0 * 255.0
			if !colorCloseTo(rgb{int(rFloat), int(gFloat), int(bFloat)}, rgb{0, 0, 0}, 30) {
				avg.r += int(rFloat)
				avg.g += int(gFloat)
				avg.b += int(bFloat)
				count++
			}
		}
	}
	if count == 0 {
		return rgb{0, 0, 0}
	}
	avg.r /= count
	avg.g /= count
	avg.b /= count
	return avg
}

func colorCloseTo(col1 rgb, col2 rgb, threshold int) bool {
	dist := (col2.r-col1.r)*(col2.r-col1.r) + (col2.g-col2.g)*(col2.g-col1.g) + (col2.b-col1.b)*(col2.b-col1.b)
	return dist < threshold*threshold
}
