package cli

import (
	"encoding/json"
	"os"
)

type App struct {
	Start       string
	StartParams []Param
	Kill        string
	KillParams  []Param
}

func (a *App) StartValues() map[string]interface{} {
	ret := map[string]interface{}{}
	for _, p := range a.StartParams {
		if p.Default != "" {
			ret[p.Name] = p.Default
		}
	}
	return ret
}

func (a *App) KillValues() map[string]interface{} {
	ret := map[string]interface{}{}
	for _, p := range a.KillParams {
		if p.Default != "" {
			ret[p.Name] = p.Default
		}
	}
	return ret
}

type Param struct {
	Name    string
	Type    string
	Default string
}

var rendezVousCliVerb = "rendez-vous"

var args []string

func Args(a ...string) {
	args = a
}

func init() {
	Args(os.Args...)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func RendezVous(app App) {
	if len(args) > 1 && args[1] == rendezVousCliVerb {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		fail(enc.Encode(app))
		os.Exit(0)
	}
}

// RendezVous(App{
// 	Name: "blog",
// 	Start: "blog -listen {{.listen}}",
// 	StartParams: []Params{
//     Param{Name:"listen",Type:"address",Default:":8090"}
//   },
// })
