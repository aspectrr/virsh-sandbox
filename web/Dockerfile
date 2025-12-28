FROM oven/bun:latest

# Set the working directory
WORKDIR /app

# Copy package files and lock file
COPY package.json bun.lock ./

# Install dependencies
RUN bun install

# Copy the rest of the application code
COPY . .

# Expose the port Vite uses (default 5173)
EXPOSE 5173

# Run the start command (adjust if you add a start script)
CMD ["bun", "run", "start"]
