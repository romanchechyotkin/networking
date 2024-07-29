# run

```bash
python -m venv venv

# The standard RTSP port is 554, but you will need to choose a port number greater than 1024.
python server.py <port>

python client_launcher.py <server_host> <server_port> <rtp_port> <filename>
``` 