# Nats dummy server
# Listens on port 3000

Install:
    - `go install github.com/JeremiahVaughan/nats-dummy@latest`
    - `sudo cp nats-dummy.service /etc/systemd/system/nats-dummy.service`
    - `sudo systemctl enable nats-dummy.service`
    - `sudo systemctl start nats-dummy.service`
    - `sudo systemctl status nats-dummy.service`
