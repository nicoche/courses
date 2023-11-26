import http.server
import os
import socketserver

DIRECTORY = "_website"

# from https://stackoverflow.com/a/52531444

class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)


if __name__ == '__main__':
    host = os.environ.get("HOST", "0.0.0.0")
    port = int(os.environ.get("PORT", "8000"))

    with socketserver.TCPServer((host, port), Handler) as httpd:
        print("serving at port", port)
        httpd.serve_forever()

