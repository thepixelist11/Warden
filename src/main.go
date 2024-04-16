package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/term"
)

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

func main() {
	width, height := getTerminalSize()

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

	var body string
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
		body, err := io.ReadAll(res.Body)
		_ = body
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Open(path.Join(loadFile))
		if err != nil {
			log.Fatal(err)
		}
		_ = file
	}
	var status Status
	// json.Unmarshal(body, &status)

	_ = body
	if exIcon && !monitor {
		status.Icon = ""
	}

	formattedJSON := formatStatus(status)

	if monitor {
		renderMonitor(status, width, height)
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

func renderMonitor(status Status, width int, height int) {
	// termWidth, termHeight := ui.TerminalDimensions()
	var images []image.Image
	image, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(parseIconData(status.Icon))))
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}
	images = append(images, image)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	img := widgets.NewImage(nil)
	// img.SetRect(0, 0, 20, 15)
	img.Border = false

	hostname := widgets.NewParagraph()
	hostname.Text = status.Hostname
	// hostname.SetRect(0, 0, 20, 5)
	hostname.Border = false
	hostname.TextStyle.Fg = ui.Color(6)

	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)
	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0, img),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0, hostname),
		),
	)

	render := func() {
		img.Image = images[0]
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
