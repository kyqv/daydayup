docker build -t demo:v1.0 .

docker images

docker run -p 8080:8080 demo:v1.0

docker run -d -p 8080:8080 demo:v1.0