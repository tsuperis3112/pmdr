# pmdr - A simple Pomodoro timer for your terminal

`pmdr` is a minimalist, yet powerful Pomodoro Technique timer designed to run as a background daemon, helping you stay focused and productive right from your command line.

## Overview

The Pomodoro Technique is a time management method that uses a timer to break down work into intervals, traditionally 25 minutes in length, separated by short breaks. `pmdr` brings this technique to your terminal with a simple, unobtrusive daemon-based approach.

It runs in the background, so you can start a timer and forget about it. Use simple commands to check the status, pause, or stop your session. With powerful hook support, you can integrate `pmdr` with your desktop environment to get native notifications or trigger any custom script.

## Features

- **Daemon-based:** Runs as a background process, leaving your terminal free.
- **Simple Commands:** An intuitive command set (`start`, `status`, `pause`, `resume`, `stop`, `config`).
- **Customizable Timers:** Easily configure work, short break, and long break durations via config file or command-line flags.
- **Spoken Notifications:** Speaks notifications at the beginning of each session (e.g., "Work session started") using native OS text-to-speech engines.
- **Powerful Hooks:** Execute any shell command on timer events (e.g., session completion), allowing for native desktop notifications and other integrations.
- **Configuration-driven:** Simple YAML configuration file for easy customization.

## Installation

If you have a Go environment set up, you can install `pmdr` with a single command:

```sh
go install github.com/tsuperis3112/pmdr@latest
```

Make sure your `$(go env GOPATH)/bin` directory is in your system's `PATH`.

## Usage

`pmdr` is controlled through simple commands.

### Timer Controls

- **`pmdr start [flags]`**: Starts a new Pomodoro session. You can override config settings with flags.
  - `-w, --work <duration>`: Set work session duration (e.g., `25m`).
  - `-s, --short-break <duration>`: Set short break duration (e.g., `5m`).
  - `-l, --long-break <duration>`: Set long break duration (e.g., `15m`).
  - `-c, --cycles <number>`: Set number of work cycles before a long break.
- **`pmdr status`**: Shows the current status of the timer (e.g., session type, remaining time).
- **`pmdr pause`**: Pauses the current session.
- **`pmdr resume`**: Resumes a paused session.
- **`pmdr stop`**: Stops the timer and the daemon completely.

### Configuration Management

- **`pmdr config init`**: Creates a default configuration file.
- **`pmdr config status`**: Shows the path of the configuration file being used.
- **`pmdr config edit`**: Opens the current configuration file in your default editor.

## Configuration

On the first run, `pmdr` doesn't require a configuration file and will use sensible defaults. You can create a default file with `pmdr config init` to customize its behavior.

**Config file locations:**

`pmdr` looks for a configuration file in the following locations, in this order. Both `.yaml` and `.yml` extensions are supported.

1. `./.pmdr.yaml` or `./.pmdr.yml` (in the current directory)
2. `~/.pmdr/config.yaml` or `~/.pmdr/config.yml` (in your home directory)
3. `$XDG_CONFIG_HOME/pmdr/config.yaml` or `$XDG_CONFIG_HOME/pmdr/config.yml` (e.g. `~/.config/pmdr/config.yaml`)

### Example `config.yaml`

```yaml
# pmdr configuration file
# For more information, see: https://github.com/tsuperis3112/pmdr

# Timer durations (any valid Go time duration string, e.g., "25m", "1h30m")
work_duration: 25m
short_break_duration: 5m
long_break_duration: 15m

# Number of work cycles before a long break
pomo_cycles: 4

# Hooks: execute shell commands on events
hooks:
  # Triggered when a work session finishes
  work:
    # Example for macOS native notification
    # - "osascript -e 'display notification \"Work session complete! Time for a break.\" with title \"Pmdr\"'"
    # Example for Linux native notification (with libnotify)
    # - "notify-send \"Pmdr\" \"Work session complete! Time for a break.\""
  # Triggered when a short break session finishes
  short_break:
    # - "osascript -e 'display notification \"Break is over! Time for work.\" with title \"Pmdr\"'"
  # Triggered when a long break session finishes
  long_break:
    # - "osascript -e 'display notification \"Long break is over! Time for work.\" with title \"Pmdr\"'"""
```

### Desktop Notifications via Hooks

The `hooks` feature is the recommended way to get desktop notifications, which can be combined with the built-in spoken notifications.

- **macOS:**

  ```yaml
  hooks:
    work:
      - "osascript -e 'display notification "Time for a break!" with title "Pmdr - Work Complete"'"
    short_break:
      - "osascript -e 'display notification "Break is over!" with title "Pmdr - Back to Work"'"
  ```

- **Linux (with `libnotify`):**

  ```yaml
  hooks:
    work:
      - "notify-send "Pmdr - Work Complete" "Time for a break!""
    short_break:
      - "notify-send "Pmdr - Back to Work" "Break is over!""
  ```

## License

Copyright (c) 2025 Takeru Furuse.

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/tsuperis3112/pmdr/issues)
