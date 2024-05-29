FROM nginx:latest

RUN echo "hello from the auth service" > /usr/share/nginx/html/index.html

EXPOSE 80
