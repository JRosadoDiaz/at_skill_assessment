# Action Target Skill Assessment

This is a Golang project that creates pings a given list of hosts and displays the results onto a webpage.

## Installation
Note: This project uses the Golang programming language which can be installed here: https://go.dev/.

Clone the project

```bash
  git clone https://github.com/JRosadoDiaz/at_skill_assessment.git
```

Go to the project directory

```bash
  cd at_skill_assessment
```


Once running, a webserver will run at the specified port or on `http://localhost:8000` by default. There are two alternative methods to run the project:

### Run locally through command line
This project utilizes priviledged pings and therefore requires advanced priviledges to opperate correctly. For this reason, this project may not work on Windows opperating systems unless priviledges are elevated and firewalls disabled.
this project includes 4 command line flags, all with their own default values: 

`hosts`: Comma-seperated list of hosts to ping

`port`: Port number the web server will run on

`interval`: The interval between pings

`count`: Number of times the host will be pinged, 0 is ping indefinitely
```bash
  sudo go run ./cmd/main.go -hosts="www.google.com,www.reddit.com" -port="8000" -interval=5 -count=0
```

### Running as a daemon
For Linux machines, a systemd file can be found in `/init/pinger.service`. To install, the service file must be copied into your `/etc/systemd/system/` directory and activate
```bash
  // Copies service to systemd folder
  sudo cp ./init/pinger.service /etc/systemd/system/pinger.service
  
  // Reloads list of daemons
  sudo systemctl daemon-reload
  
  // Activates the service
  sudo systemctl start pinger.service
  
  // Checks status of service
  sudo systemctl status pinger.service
  
  // Disable and remove service
  sudo systemctl disable pinger.service
  sudo systemctl stop pinger.service
```
