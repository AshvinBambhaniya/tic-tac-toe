# Stage 1: Build the application
FROM node:22-alpine AS build

WORKDIR /app

# Install dependencies first to leverage caching
COPY package*.json ./
RUN npm install

# Copy source and build
COPY . .
RUN npm run build

# Stage 2: Runtime image
FROM node:22-alpine

WORKDIR /app

# Install runtime system dependencies
RUN apk add --no-cache nginx supervisor curl

# Create directory for nginx
RUN mkdir -p /run/nginx

# Copy configuration files from the build stage/context
COPY ./docker/nginx.conf /etc/nginx/nginx.conf
COPY ./docker/supervisord.conf /etc/supervisord.conf

# Copy ONLY the build output and necessary files for production
COPY --from=build /app/.output ./.output

# Expose the Nginx port
EXPOSE 80

# Environment variables for production (defaults)
ENV NODE_ENV=production

# Command to start everything
CMD ["supervisord", "-c", "/etc/supervisord.conf"]
