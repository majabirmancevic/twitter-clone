FROM nginx:latest

COPY api_gateway/api_gateway.conf /etc/nginx/conf.d/api_gateway.conf
COPY auth_service/certificates /etc/nginx/certificates
RUN rm /etc/nginx/conf.d/default.conf
