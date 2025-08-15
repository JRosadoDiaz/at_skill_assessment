This is the developer systemd unit file

** INSTALLATION **
When ready to deploy as a service, copy file into the following linux directory:

/etc/systemd/system/

This can be done through the following commands:

-- Copy to dir
sudo cp pinger.service /etc/systemd/system/pinger.service

-- Reload systemd to see new service
sudo systemctl daemon-reload

-- Enable/Disable the service to start automatically on boot
sudo systemctl enable/disable pinger.service

-- Start the service manually
sudo systemctl start pinger.service

-- Check status
sudo systemctl status pinger.service