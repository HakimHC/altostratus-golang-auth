FROM nginx:latest

RUN mkdir -p /usr/share/nginx/html/api/auth

RUN echo "hello from the auth service" > /usr/share/nginx/html/api/auth/index.html

EXPOSE 80
