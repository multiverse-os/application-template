[<img src="https://avatars2.githubusercontent.com/u/24763891?s=400&u=c1150e7da5667f47159d433d8e49dad99a364f5f&v=4"  width="256px" height="256px" align="right" alt="Multiverse OS Logo">](https://github.com/multiverse-os)

## Multiverse: Application Template 
**URL** [multiverse-os.org](https://multiverse-os.org)

An application template for use as a skeleton/boilerplate that can be the basis 
of a template for Multiverse OS's `laboratory` development framework. Intended 
to simplify, enhance, and speed-up the development of Multiverse OS
applications. 

This template can serve as its own boilerplate, providing scaffolding for either 
command-line interface (CLI), the web interface or web server daemon, or a 
graphical user interface (GUI) application.  

The Multiverse OS design guide outlines a structure that is required by all core
utilities and applications included by default (the default is defined by the
type of system being installed, Multiverse OS can be used as a headless server,
or a general use work station, comming fully capable of running OSX, Linux, and
even windows PE binaries). 

This project is intended to both simplify development of new application,
intended for use by Multiverse OS, or provided for the community and developed
collectively by the community. It illustrates some key design principals that
allow it to remain consistent with all other core applications, and it allows
interested developers in quickly learning the key design requirements to
volunteer and join the Multiverse OS development community. 

# Application Framework Package 
After building many Go langauge applications, there are portions of every
application that repeat regardless if it is a web application, command-line tool
or client regardless of complexity. 

This framework is an attempt to provide a framework for the most basic
functionality expected of any complete and feature complete application. Such as
handling PID files, loading/saving configuration files, semantic versioning,
dealing with local data path, and local config path. 

This will likely be used with our web framework pulling out logic that is used
in the web framework into this, and so the web-framework makes use of this while
leaving just the web-framework logic in the web-framework package. 


