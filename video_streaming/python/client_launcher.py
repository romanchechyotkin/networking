import sys
from tkinter import Tk

from client import Client

if __name__ == "__main__":
	try:
		server_addr = sys.argv[1]
		server_port = sys.argv[2]
		rtp_port = sys.argv[3]
		filename = sys.argv[4]
	except:
		print("[Usage: ClientLauncher.py Server_name Server_port RTP_port Video_file]\n")	
	
	root = Tk()
	
	# Create a new client
	app = Client(root, server_addr, server_port, rtp_port, filename)
	app.master.title("RTPClient")
	print("RTPClient started")	
	root.mainloop()
