This small application was a learning exercise to develop my Angular and Go skills.
The application allows you to take notes using the Markdown syntax.
The application is backed by a SQLite database.

# Installing 

The current project is made to be able to run under Linux and Windows.

## Requirement

* For the backend
    * [Go compiler](https://golang.org/dl/)
    * [mingw-w64](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe/download) - For Windows (Be sure to add the `bin` folder of mingw-w64 to your `PATH` to make gcc available in the commandline.)
    * GCC - For Linux
* For the frontend
    * [NodeJS](https://nodejs.org/en/)

## How to install

To build the project, simply run the `build.sh` script available at the root of the project.
The result of the build will be available in the `build` folder.

# Using

Simply run the `gonote(.exe)` executable available in the `build` folder.
Once the application tells you it's listening open a web browser and navigate to: [localhost:8080](http://localhost:8080/).

# What if I don't want to build the sources?

First you'll need to add two environment variable. On Windows you can run the powershell script `set-env.ps1` to add the required variables. On linux you'll need to add the following lines to your `~/.bashrc`.

```
export GOPATH="<PATH TO THE PROJECT ROOT ON YOUR SYSTEM.>"
export CGO_ENABLED=1
```

Then you still need to install the dependancies for both the frontend and the backend.
You can do that by running `npm install` from the `/ngNote/` folder and then `go get` from the `/src/gonote/` folder.

Then the backend can be run using the `go run main.go` command from `/src/gonote`.
But if you don't place the frontend files into the `public` folder inside the project you'll get an error when you'll try to access the application, therefore you have two solutions for the frontend.

1. You can build the frontend using the command `npx ng build` from the `/ngNote/` folder and move the files located in `/ngNote/dist/` in the `/src/gonote/public/` folder. OR
2. You can run the frontend on a different port using `npx ng serve` then access the application using the address indicated in the terminal.

# Can I host the application on the internet?

I guess you could host you copy of the application, but it's not really web ready since the frontend always place the request using localhost.

If you want to host a web version, you'll need to modify the frontend base api url. The variable is located in: `/ngNote/src/app/class/common.ts`.

example
```ts
- export const ApiUrl: string = "http://localhost:8080/api";
+ export const ApiUrl: string = "http://my-domain:8080/api";
```

If you want to change the hosting port you can modify `/src/gonote/main.go`.

example
```go
- log.Fatal(http.ListenAndServe(":8080", nil))
+ log.Fatal(http.ListenAndServe(":80", nil))
```

Once you've done your changes, just build the application with the `build.sh` script.

# Thanks

I hope you'll like it even if it's not super useful nor secured.
