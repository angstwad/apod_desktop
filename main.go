// Copyright 2014 Paul Durivage <pauldurivage@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	url     = "http://apod.nasa.gov"
	scriptf = "/tmp/background.scpt"
	pattern = `<a href="(image\/\d{4}\/\w+\.jpg|png)"`
)

func set_background(fname string) {
	fmt.Println("Setting APOD picture to desktop background.")

	cmdbytes := []byte(fmt.Sprintf(`tell application "System Events" to set picture of every desktop to "%s"`, fname))

	var mode os.FileMode = 0644
	e := ioutil.WriteFile(scriptf, cmdbytes, mode)
	if e != nil {
		fmt.Printf("Error writing AppleScript file: %s\n", e.Error())
		os.Exit(1)
	}

	_, e = exec.Command("/usr/bin/osascript", scriptf).CombinedOutput()
	if e != nil {
		fmt.Printf("Error setting APOD picture to background: \n%s\n", e.Error())
		os.Exit(1)
	}

	e = os.Remove(scriptf)
	if e != nil {
		fmt.Printf("Error deleting file: %s\n", e.Error())
	}
}

func download_image(url string, uri string) string {
	fmt.Printf("Downloading photo...")
	img_url := fmt.Sprintf("%s/%s", url, uri)
	resp, e := http.Get(img_url)
	if e != nil {
		fmt.Printf("Error downloading APOD photo: %s\n", e.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Printf("Error reading response from APOD: %s\n", e.Error())
		os.Exit(1)
	}

	ext := strings.Split(uri, ".")[1]
	fname := fmt.Sprintf("/tmp/apod.%s", ext)
	var mode os.FileMode = 0644
	e = ioutil.WriteFile(fname, body, mode)
	if e != nil {
		fmt.Printf("Error writing file to ", e.Error())
		os.Exit(1)
	}
	fmt.Println("Done.")
	fmt.Printf("Photo saved to %s.\n", fname)

	return fname
}

func get_image_uri(page []byte) string {
	fmt.Println("Getting APOD image URL.")
	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(string(page), 1)
	if match == nil {
		fmt.Println("No image found today!")
		os.Exit(0)
	}
	return match[0][1]
}

func fetch_page(url string) []byte {
	fmt.Println("Fetching APOD page.")
	resp, e := http.Get(url)

	defer resp.Body.Close()

	if e != nil {
		fmt.Printf("Error while contacting APOD: %s\n", e.Error())
		os.Exit(1)
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Printf("Error reading response from APOD: %s\n", e.Error())
		os.Exit(1)
	}
	return body
}

func check_conn() {
	fmt.Print("Checking connection...")
	for i := 0; i < 12; i++ {
		_, e := http.Get(url)
		if e != nil {
			fmt.Println("\nThere was a problem; sleeping for 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Println("Ok.")
		break
	}

}

func main() {
	check_conn()
	page := fetch_page(url)
	img_uri := get_image_uri(page)
	fname := download_image(url, img_uri)
	set_background(fname)
}
