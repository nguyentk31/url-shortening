target "default" {
  context = "./"
  dockerfile = "Dockerfile"
  tags = [ "myapp:latest" ]
}