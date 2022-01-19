# Kgicon CLI

> Command Line Interface Tool for KgIcons Laravel Package

For use in https://github.com/kaisargaming/package-kg-icons

## Usage

Available functionality as of now;

- `kgicon-cli prep {hero|majestic}`, this command will basically prepares the icon files from upstream repository and varied in the provider will copy, rename and adjust target directory inside `resources/providers/`

- `kgicon-cli createlist resources/provider/{hero|majestic}` will create a `json` formatted file in the `resources/js` based on the provider code eg. `hero.json`, it can be used to loop through the icon names for documentation or cheatsheet.

## Build

Please execute your own compile if you are not using `macos`. Clone and execute `go build` inside the cloned repo path.