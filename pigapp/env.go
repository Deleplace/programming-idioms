package main

// Env encapsulates a Programming-Idioms webapp environment.
type Env struct {
	// Host depends on the target environment.
	// It should not have a trailing slash.
	IsDev           bool
	Host            string
	UseAbsoluteUrls bool
	UseMinifiedCss  bool
	UseMinifiedJs   bool
}

//
// Prod
//

var envProd = Env{
	IsDev: false,
	// Host:            "https://www.programming-idioms.org",
	Host:            "https://pi-go111.appspot.com",
	UseAbsoluteUrls: false,
	UseMinifiedCss:  true,
	UseMinifiedJs:   true,
}

//
// Dev
//

var envDev = Env{
	IsDev:           true,
	Host:            "http://localhost:8080",
	UseAbsoluteUrls: true,
	UseMinifiedCss:  false,
	UseMinifiedJs:   false,
}

// Which one is used ?

var env Env

func initEnv() {
	if isAppengineDevServer() {
		env = envDev
	} else {
		env = envProd
	}
}
