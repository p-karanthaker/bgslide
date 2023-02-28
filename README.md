# bgslide
An application to set the desktop background to be a slideshow of images.

## Motivations
I use Pop!_OS+GNOME and Windows. I wanted to set my backgrounds as a slideshow of photos that I have taken.

- Pop!_OS doesn't appear to have an option to set the desktop background as a slideshow of images.
- Windows allows slideshows but sets a different image on multiple displays - I wanted the same image on both displays.
- I wanted an excuse to practice Go.

## Supports
- Linux
    - GNOME Desktop Environments
- Windows

## Building/Installing
### Linux
Run your standard `go build` and `go build install`

### Windows
Build with flags so Windows doesn't open a terminal window when running the binary.

`go build -ldflags -H=windowsgui`

`go install -ldflags -H=windowsgui`

## Usage
The `-dir` flag will default to the `Pictures` folder in your home directory if the environment variable `HOME` or `USERPROFILE` exists on your system.
```
Usage of ./bgslide:
  -dir string
        The directory containing the wallpapers. (default "/home/karan/Pictures")
  -interval duration
        The interval for changing wallpaper. E.g. 300s, 5m, 1h. Minimum of 5m (default 30m0s)
```

## Running automatically on startup
### Linux
1. Create a `bgslide.service` file in `/etc/systemd/user` containing:
```ini
[Unit]
Description=A service to set the desktop background to be a slideshow of images.

[Service]
Type=simple
ExecStart=%h/go/bin/bgslide

[Install]
WantedBy=default.target
```
2. `systemctl --user daemon-reload`
3. `systemctl --user enable bgslide.service`
4. `systemctl --user start bgslide`

### Windows
N.B. Windows may flag the application as unsafe and quarantine it.
1. <kbd>Win</kbd>+<kbd>R</kbd>, type `shell:startup`
2. Create a shortcut to the `bgslide.exe` in your GOPATH adding flags if you want to change the defaults.
3. Run the shortcut to start immediately, otherwise it'll run on the next startup.
