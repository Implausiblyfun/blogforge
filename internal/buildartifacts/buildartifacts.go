package buildartifacts

import (
	"encoding/json"
	"fmt"
	"os"
)

// Build stores the build aritifact's information for later use.
// Its a little nasty to expose it like this but its fairly lightweight and likely used in informational spaces across the code base.
var Build Info

// Info stores all the info around the current build and the like.
type Info struct {
	Version string `json:"version,omitempty"`
}

// LoadsBuildInfo for later access.
func LoadBuildInfo() {
	fmt.Println("Loading Build Information")
	Build = Info{
		Version: "unknown",
	}

	f, err := os.Open("./version.json")
	if err != nil {
		fmt.Println("Unknown version please get a valid build")
		return
	}

	dc := json.NewDecoder(f)
	err = dc.Decode(&Build)
	if err != nil {
		// fmt as we dont have a handle on the logger due to this being in an init right now.
		fmt.Println("Error decoding the build info")
	}

}
