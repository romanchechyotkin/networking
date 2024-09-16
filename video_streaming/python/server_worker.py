import socket
import threading
from random import randint
from collections import defaultdict

from rtp_packet import RtpPacket
from video_stream import VideoStream


class ServerWorker:
    SETUP = 'SETUP'
    PLAY = 'PLAY'
    PAUSE = 'PAUSE'
    TEARDOWN = 'TEARDOWN'

    INIT = 0
    READY = 1
    PLAYING = 2
    state = INIT

    OK_200 = 0
    FILE_NOT_FOUND_404 = 1
    CON_ERR_500 = 2

    def __init__(self, clientInfo):
        self.clientInfo = clientInfo
        self.rooms = defaultdict(list)
        self.events = dict()
        self.workers = dict()

    def run(self):
        threading.Thread(target=self.recvRtspRequest).start()

    def recvRtspRequest(self):
        """Receive RTSP request from the client."""
        conn_socket = self.clientInfo['rtspSocket'][0]
        print("connection socket", conn_socket)
        while True:
            data = conn_socket.recv(256)
            if data:
                print("Data received:\n" + data.decode("utf-8"))
                self.processRtspRequest(data.decode("utf-8"), conn_socket)

    def processRtspRequest(self, data, conn_socket):
        """Process RTSP request sent from the client."""
        # Get the request type
        print("got request", data)

        request = data.split('\n')
        line1 = request[0].split(' ')
        request_type = line1[0]

        # Get the media file name
        filename = line1[1]

        # Get the RTSP sequence number
        seq = request[1].split(' ')

        # Process SETUP request
        if request_type == self.SETUP:
            if self.state == self.INIT:
                # Update state
                print("processing SETUP\n")

                try:
                    self.clientInfo['videoStream'] = VideoStream(filename)
                    self.state = self.READY
                except IOError:
                    self.replyRtsp(self.FILE_NOT_FOUND_404, seq[1])

                # Generate a randomized RTSP session ID
                self.clientInfo['session'] = randint(100000, 999999)

                # Send RTSP reply
                self.replyRtsp(self.OK_200, seq[1])

                # Get the RTP/UDP port from the last line
                self.clientInfo['rtpPort'] = request[2].split(' ')[3].strip(";")

                room_id = request[2].split(' ')[-1]
                client_conn = (conn_socket.getpeername(), self.clientInfo['rtpPort'])

                self.rooms[room_id].append(client_conn)


        # Process PLAY request
        elif request_type == self.PLAY:
            if self.state == self.READY:
                print("processing PLAY")
                room_id = request[2].split(' ')[-1]

                self.state = self.PLAYING

                # Create a new socket for RTP/UDP
                self.clientInfo["rtpSocket"] = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

                self.replyRtsp(self.OK_200, seq[1])

                # Create a new thread and start sending RTP packets
                self.events.setdefault(room_id, threading.Event())
                self.workers.setdefault(room_id, threading.Thread(target=self.sendRtp, args=(room_id,)))
                self.workers[room_id].start()

                print(f"ROOM={room_id} workers={self.workers[room_id]} events={self.events[room_id]}")

        # Process PAUSE request
        elif request_type == self.PAUSE:
            if self.state == self.PLAYING:
                print("processing PAUSE\n")
                room_id = request[2].split(' ')[-1]

                self.state = self.READY

                self.events[room_id].set()

                self.replyRtsp(self.OK_200, seq[1])

        # Process TEARDOWN request
        elif request_type == self.TEARDOWN:
            print("processing TEARDOWN\n")
            room_id = request[2].split(' ')[-1]

            self.events[room_id].set()

            self.replyRtsp(self.OK_200, seq[1])

            # Close the RTP socket
            self.clientInfo['rtpSocket'].close()

    def sendRtp(self, room_id):
        """Send RTP packets over UDP."""
        while True:
            self.events[room_id].wait(0.05)

            # Stop sending if request is PAUSE or TEARDOWN
            if self.events[room_id].isSet():
                break

            data = self.clientInfo['videoStream'].nextFrame()

            if data:
                frameNumber = self.clientInfo['videoStream'].frameNbr()
                try:
                    for conn in self.rooms[room_id]:
                        address = conn[0][0]
                        port = int(conn[1])
                        self.clientInfo['rtpSocket'].sendto(self.makeRtp(data, frameNumber), (address, port))
                except Exception as e:
                    print("Connection Error", e)
                # print('-'*60)
                # traceback.print_exc(file=sys.stdout)
                # print('-'*60)

    def makeRtp(self, payload, frameNbr):
        """RTP-packetize the video data."""
        version = 2
        padding = 0
        extension = 0
        cc = 0
        marker = 0
        pt = 26  # MJPEG type
        seqnum = frameNbr
        ssrc = 0

        rtpPacket = RtpPacket()

        rtpPacket.encode(version, padding, extension, cc, seqnum, marker, pt, ssrc, payload)

        return rtpPacket.getPacket()

    def replyRtsp(self, code, seq):
        """Send RTSP reply to the client."""
        if code == self.OK_200:
            # print("200 OK")
            reply = 'RTSP/1.0 200 OK\nCSeq: ' + seq + '\nSession: ' + str(self.clientInfo['session'])
            connSocket = self.clientInfo['rtspSocket'][0]
            connSocket.send(reply.encode())

        # Error messages
        elif code == self.FILE_NOT_FOUND_404:
            print("404 NOT FOUND")
        elif code == self.CON_ERR_500:
            print("500 CONNECTION ERROR")
