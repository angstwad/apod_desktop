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
    "time"
    "log"
    "errors"
)

const (
    baseurl = "http://apod.nasa.gov"
    pageurl = "http://apod.nasa.gov/apod/astropix.html"
    pattern = `(?i)<img src="(image\/\d{4}\/\w+\.jpg|png|jpeg)"`
)

func setBackground(fname string) {
    fmt.Println("Setting APOD picture to desktop background.")

    cmd := `tell application "System Events"
    set desktopCount to count of desktops
    repeat with desktopNumber from 1 to desktopCount
        tell desktop desktopNumber
            set picture to "%s"
        end tell
    end repeat
end tell`

    cmdbytes := []byte(fmt.Sprintf(cmd, fname))

    tmp, e := ioutil.TempFile("", "apod")
    if e != nil {
        log.Fatal("Error opening temp file.", e.Error())
    }

    defer os.Remove(tmp.Name())

    if _, e := tmp.Write(cmdbytes); e != nil {
        log.Fatal(fmt.Sprintf("Error writing AppleScript file: %s\n", e.Error()))
    }

    if _, e = exec.Command("/usr/bin/osascript", tmp.Name()).CombinedOutput(); e != nil {
        log.Fatal(fmt.Sprintf("Error setting APOD picture to background: \n%s\n", e.Error()))
    }
}

func downloadAndSet(url string, uri string) {
    fmt.Println("Downloading photo...")
    imgURL := fmt.Sprintf("%s/%s", url, uri)
    resp, e := http.Get(imgURL)
    if e != nil {
        log.Fatal(fmt.Sprintf("Error downloading APOD photo: %s\n", e.Error()))
    }

    defer resp.Body.Close()

    body, e := ioutil.ReadAll(resp.Body)
    if e != nil {
        log.Fatal(fmt.Sprintf("Error reading response from APOD: %s\n", e.Error()))
    }

    tmp, e := ioutil.TempFile("", "apod")
    if e != nil {
        log.Fatal("Error opening temp file. ", e.Error())
    }

    defer os.Remove(tmp.Name())

    if _, err := tmp.Write(body); err != nil {
        log.Fatal("Error writing file to ", e.Error())
    }
    fmt.Println("Done.")
    fmt.Printf("Image saved to %s.\n", tmp.Name())

    setBackground(tmp.Name())
}

func getImgURI(page []byte) string {
    fmt.Println("Getting APOD image URL.")
    re := regexp.MustCompile(pattern)
    match := re.FindAllStringSubmatch(string(page), 1)
    if match == nil {
        log.Fatal("No image found today!")
    }
    fmt.Printf("Image URI is '%s'.\n", match[0][1])
    return match[0][1]
}

func fetchPage(url string) []byte {
    fmt.Println("Fetching APOD page.")
    resp, e := http.Get(url)

    defer resp.Body.Close()

    if e != nil {
        log.Fatal(fmt.Sprintf("Error while contacting APOD: %s\n", e.Error()))
    }
    body, e := ioutil.ReadAll(resp.Body)
    if e != nil {
        log.Fatal(fmt.Sprintf("Error reading response from APOD: %s\n", e.Error()))
    }
    return body
}

func checkConn() error {
    fmt.Print("Checking connection...")
    tries := 12
    for i := 0; i < tries; i++ {
        if _, e := http.Get(pageurl); e != nil {
            if i == 0 {
                fmt.Println()
            }
            fmt.Printf("Couldn't connect, sleeping for 5 seconds; remaining tries: %v\n", tries - (i + 1))
            time.Sleep(5 * time.Second)
            continue
        }
        fmt.Println("Ok.")
        return nil
    }
    return errors.New("could not connect to apod.nasa.gov")
}

func main() {
    if err := checkConn(); err != nil {
        log.Fatal(err)
    }
    page := fetchPage(pageurl)
    imgURI := getImgURI(page)
    downloadAndSet(baseurl, imgURI)
}
