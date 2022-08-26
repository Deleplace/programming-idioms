package main

import (
	"net/http"
	"sort"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/log"
)

// ApplicationConfig is a global configuration container.
type ApplicationConfig struct {
	// Id of this particular configuration set
	Id int
	// Map of configuration properties
	Toggles Toggles
}

/* Deprecated
// Load reads properties into this ApplicationConfig.
func (ac *ApplicationConfig) Load(ch <-chan datastore.Property) error {
	ac.Toggles = Toggles{}
	// Another modeling with less datastore columns: use AppConfigProperty instead
	for p := range ch {
		ac.Toggles[p.Name] = (p.Value).(bool)
	}
	return nil
}

// Save writes properties from this ApplicationConfig.
func (ac *ApplicationConfig) Save(ch chan<- datastore.Property) error {
	// Another modeling with less datastore columns: use AppConfigProperty instead
	for n, v := range ac.Toggles {
		ch <- datastore.Property{Name: n, Value: v, NoIndex: false, Multiple: false}
	}
	close(ch)
	return nil
}
*/

// Before first request, toggles are "default" and are not loaded
// from datastore yet.
// TODO: use this to estimate config freshness. Maybe use type time.Time instead.
var configTime = "0" //time.Now().Format("2006-01-02_15-04")

func refreshToggles(ctx context.Context) error {
	appConfig, err := dao.getAppConfig(ctx)
	if err == appConfigPropertyNotFound {
		// Nothing in Memcache, nothing in Datastore!
		// Then, init default (hard-coded) toggle values and persist them.
		initToggles()
		log.Infof(ctx, "Saving default Toggles to Datastore...")
		err := dao.saveAppConfig(ctx, ApplicationConfig{Id: 0, Toggles: toggles})
		if err == nil {
			log.Infof(ctx, "Default Toggles saved to Datastore.")
			configTime = time.Now().Format("2006-01-02_15-04")
		}
		return err
	}
	if err != nil {
		log.Errorf(ctx, "Error while loading ApplicationConfig from datastore: %v", err)
		return err
	}
	toggles = appConfig.Toggles
	configTime = time.Now().Format("2006-01-02_15-04")
	log.Infof(ctx, "Updated Toggles from memcached or datastore\n")
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
	toggles["idiomVotingUp"] = false
	toggles["idiomVotingDown"] = false
	toggles["showIdiomRating"] = false
	toggles["implVotingUp"] = false
	toggles["implVotingDown"] = false
	toggles["showImplRating"] = false
	toggles["languageCreation"] = false

	// Homepage
	toggles["homeBlockCoverage"] = true
	toggles["homeBlockAllIdioms"] = true
	toggles["homeBlockLastUpdated"] = true
	toggles["homeBlockPopular"] = false

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
			return PiErrorf(http.StatusForbidden, "Not available for now.")
		}
	}
	return nil
}
