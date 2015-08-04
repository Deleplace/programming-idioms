package pigae

import (
	"net/http"
	"sort"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
)

// ApplicationConfig is a global configuration container.
type ApplicationConfig struct {
	// Id of this particular configuration set
	Id int
	// Map of configuration properties
	Toggles Toggles
}

// Load reads properties into this ApplicationConfig.
func (ac *ApplicationConfig) Load(ch <-chan datastore.Property) error {
	ac.Toggles = Toggles{}
	// Todo another modeling with less datastore columns?
	for p := range ch {
		ac.Toggles[p.Name] = (p.Value).(bool)
	}
	return nil
}

// Save writes properties from this ApplicationConfig.
func (ac *ApplicationConfig) Save(ch chan<- datastore.Property) error {
	// Todo another modeling with less datastore columns?
	for n, v := range ac.Toggles {
		ch <- datastore.Property{Name: n, Value: v, NoIndex: false, Multiple: false}
	}
	close(ch)
	return nil
}

// Before first request, toggles are "default" and are not loaded
// from datastore yet
var configTime = "0" //time.Now().Format("2006-01-02_15-04")

func refreshToggles(c appengine.Context) error {
	appConfig, err := dao.getAppConfig(c)
	if err != nil {
		c.Errorf("Error while loading ApplicationConfig from datastore: %v\n", err)
		return err
	}
	toggles = appConfig.Toggles
	configTime = time.Now().Format("2006-01-02_15-04")
	c.Infof("Updated Toggles from memcached or datastore\n")
	// _ = appConfig

	return err
}

//
// A toggle should always be named after the positive feature it represents,
// and default value should be true.
//
var toggles = Toggles{}

func initToggles() {
	// These two toggles block everything
	toggles["online"] = true
	toggles["writable"] = true

	// The toggles below require "online"
	toggles["searchable"] = true

	toggles["loggable"] = false
	toggles["greetings"] = true

	toggles["languageBar"] = true
	toggles["syntaxColoring"] = true

	toggles["licenseDisclaimer"] = true
	toggles["poweredBy"] = true

	// The toggles below require "writable"
	toggles["anonymousWrite"] = false
	toggles["idiomCreation"] = true
	toggles["idiomEditing"] = true
	toggles["implAddition"] = true
	toggles["implEditing"] = true
	toggles["pictureEditing"] = false
	toggles["idiomVotingUp"] = true
	toggles["idiomVotingDown"] = true
	toggles["showIdiomRating"] = true
	toggles["implVotingUp"] = true
	toggles["implVotingDown"] = true
	toggles["showImplRating"] = true
	toggles["languageCreation"] = false

	// Admin
	toggles["administrable"] = true

	// Misc conf
	toggles["isDev"] = env.IsDev
	toggles["themeVirtualVersioning"] = true
	toggles["useAbsoluteUrls"] = env.UseAbsoluteUrls
	toggles["useMinifiedCss"] = env.UseMinifiedCss
	toggles["useMinifiedJs"] = env.UseMinifiedJs
	toggles["useCDN"] = false
}

func toggled(name string) bool {
	return toggles[name]
}

func copyToggles(src Toggles) Toggles {
	dest := make(Toggles, len(src))
	for k, v := range src {
		dest[k] = v
	}
	return dest
}

func allToggleNames() []string {
	names := make([]string, len(toggles))
	i := 0
	for key := range toggles {
		names[i] = key
		i++
	}
	sort.Strings(names)
	return names
}

func togglesMissing(w http.ResponseWriter, r *http.Request, toggleNames ...string) error {
	for _, name := range toggleNames {
		if !toggles[name] {
			return PiError{"Not available for now.", http.StatusForbidden}
		}
	}
	return nil
}
