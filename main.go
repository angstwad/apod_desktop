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
)

func set_background(fname string) {
	fmt.Println("Setting APOD picture to desktop background.")

	osacmd := fmt.Sprintf(`tell application "System Events" to set picture of every desktop to "%s"`, fname)
	cmdbytes := []byte(osacmd)
	scriptf := "/tmp/background.scpt"
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
	img_url := url + "/" + uri
	resp, e := http.Get(img_url)
	if e != nil {
		fmt.Printf("Error downloading APOD photo: %s\n", e.Error())
		os.Exit(1)
	}

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Printf("Error reading response from APOD: %s\n", e.Error())
		os.Exit(1)
	}

	ext := strings.Split(uri, ".")
	fname := "/tmp/apod." + ext[len(ext)-1:][0]
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
	pattern := `<a href="(image\/\d{4}\/\w+\.jpg|png)">`
	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(string(page[:]), 1)
	return match[0][1]
}

func fetch_page(url string) []byte {
	fmt.Println("Fetching APOD page.")
	resp, e := http.Get(url)
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

func main() {
	url := "http://apod.nasa.gov"
	page := fetch_page(url)
	img_uri := get_image_uri(page)
	fname := download_image(url, img_uri)
	set_background(fname)
}
