# dns-lookup-tool
A simple DNS lookup tool built using WebSockets, Vue.js, and Golang

**Build and Run with Docker**
1. Clone and enter repo
```
git clone git@github.com:cvasq/dns-lookup-tool.git
```
2. Build the docker image locally
```
docker build -t cvasquez/dns-lookup-tool .
```                                                                                                                                                                                           
3. Run the Docker container and map the listening HTTP port to localhost:
_Default Listening Port: 8080_
```
docker run -it --rm -p 8080:8080 cvasquez/dns-lookup-tool:latest
```                                                                                                                                                                                           
4. Navigate to http://localhost:8080/ in your browser
