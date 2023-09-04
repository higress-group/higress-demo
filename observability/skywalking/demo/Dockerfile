FROM openjdk:17-jdk-slim

RUN apt update && apt install -y curl

RUN mkdir -p /app
COPY target/demo-0.0.1-SNAPSHOT.jar /app/app.jar
WORKDIR /app
ENTRYPOINT [ "java", "-jar", "app.jar"]

# liuxr25/springbootdemo:svc-sequence