import http.server
import os
import socketserver

# from https://stackoverflow.com/a/52531444

if __name__ == '__main__':
    host = os.environ.get("HOST", "0.0.0.0")
    port = os.environ.get("PORT", 8000)
    directory = os.environ.get("DIRECTORY", "_website")

    class Handler(http.server.SimpleHTTPRequestHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, directory=directory, **kwargs)

    with socketserver.TCPServer((host, port), Handler) as httpd:
        print("serving at port", port)
        httpd.serve_forever()

