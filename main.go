package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
	var bedrock, silent, exIcon, montior *bool
	var ip, outputPath *string

	flag.BoolVar(bedrock, "b", false, "Bedrock Server?")
	flag.BoolVar(bedrock, "bedrock", false, "Bedrock Server?")
	flag.StringVar(ip, "ip", "0.0.0.0", "The IP of the server.")
	flag.StringVar(outputPath, "o", "", "The optional path of the output json file.")
	flag.StringVar(outputPath, "output", "", "The optional path of the output json file.")
	flag.BoolVar(silent, "s", false, "Do not output anything to the console.")
	flag.BoolVar(silent, "silent", false, "Do not output anything to the console.")
	flag.BoolVar(exIcon, "ei", true, "Do not get icon data.")
	flag.BoolVar(exIcon, "exIcon", true, "Do not get icon data.")
	flag.BoolVar(montior, "m", false, "Whether or not to use the live monitor mode.")
	flag.BoolVar(montior, "monitor", false, "WWhether or not to use the live monitor mode.")
	_ = montior

	flag.Parse()

	apiURL := ""
	if *bedrock {
		apiURL = "https://api.mcsrvstat.us/bedrock/3/" + *ip
	} else {
		apiURL = "https://api.mcsrvstat.us/3/" + *ip
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(err.Error())
		return
	}

	var status Status
	json.Unmarshal(body, &status)
	if *exIcon {
		status.Icon = ""
	}

	formattedJSON, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	if *outputPath != "" {
		os.WriteFile(*outputPath, formattedJSON, 0666)
	}

	if !*silent {
		fmt.Print(string(formattedJSON))
	}
}
