package main

import (
	"fmt"
)

func main() {
	fmt.Println("application-cli")
	fmt.Println("===============")
	fmt.Println("application command-line interface template; intended to be")
	fmt.Println("used with the library. all applications: cli tools, GUI tools")
	fmt.Println("and even web applications should be built as a Go library.\n")
	fmt.Println("after the go library is written, it is then called into the")
	fmt.Println("UI and only code related to presenting the library through")
	fmt.Println("the specified UI is in this file. so this file contains all")
	fmt.Println("the features accesible through the command-line interface.\n")
	fmt.Println("ansi coloring, terminal printing, and so on...")

	// step 1) load config values
	// env, _ := env.Parse(os.Env())
	// flags, _ := flags.Parse(os.Args())
	// file, _ := file.Load(os.File())
	//
	// signal catching on the config update signal should be handled in the
	// config library.
	//
	// (optional) File can have fnotify hooked on it to catch any changes during
	// software operation.
	//
	//
	// the order is always:
	//
	//   [env => flags => config => default values]
	//
	// the application should never error out, crash, or have any issues when
	// loading the config via this chain.
	//
	// it should always fall back to reasonable defaults or prompt the user
	// for input if defaults can not be assumed.
	//
	// use the combination of data that is collected from the method above to
	// build the application config. if there is no saved config, we should be
	// saving the constructed config values, not just the defaults.
	//
	// in the case of falling back to values that would error if default is
	// used, like in the common case the port already taken. it should *always*
	// just check the next value, so if 3000 is the default, and its taken,
	// do not error, check 3001. and repeat until a free port is located.
	//
	// if you want the error to occur for the purpose of preventing duplicate
	// processes. then you need to error out on specifically that issue, not
	// duplicate ports. this is conceptually very important, these design
	// principals are repeated at every scope of the Multiverse OS project.
	// errors must be correct, applications must never fail on issues that could
	// easily be solved better design; leaving the uesr to manually do what
	// should be computer work for no real reason outside of developers seeing
	// other applications using poor designs and repeating them.
	//
	//
	// config, _ := config.Parse(env, flags, file)
	//
	// step 2) build application object
	//
	// app := App.New()
	//
	// app has a method ran depending on flags and commands (if commands are)
	// being used. and perhaps the value/conent (last item of the cli command
	// if not a flag or command).
	//
	// if it is a command, then just do a switch case and execute the method
	// corresponding to the command.
	//
	// or preferably establish command hooks like a web application, and
	// run and then pass in the input vaues

	// switch {
	// case command.Is(xAction):
	//   result, err := app.Action()
	//   if err != nil {
	//   	app.ErrorOutput(err)
	//      os.Exit(1)
	//   } elsif result == x {
	// 		app.Output("result was x!")
	//   } else {
	//      app.Output
	// }
	//
	// default: // no command specified in a cli tool, falls back to help
	//   app.Output(app.Help())
	//
	// Or if a flag is specifying daemon mode. Have a method that holds open
	// until a signal is given.
	// This file will be pretty small, the core logic should be in the library.
	// Unless its a cli tool primarily, then it may have a large cli tool. try
	// to keep it to a single file.
	//
	// but if it does not, then have all the files in this folder, and avoid
	// special cli tool libraries when possible, they should still be in the
	// primary library not in this folder.
}
